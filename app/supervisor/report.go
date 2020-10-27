package supervisor

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/tarscloud/k8stars/app/genconf"
	"github.com/tarscloud/k8stars/consts"
	"github.com/tarscloud/k8stars/tarsproxy"
	"github.com/tarscloud/k8stars/tarsregistry/autogen/Tars"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/notifyf"
)

func (c *launchCmd) keepAlive(checkSucc bool) error {
	// get config from file
	gConf, err := genconf.GetGlobalConf()
	if err != nil {
		log.Errorf("GetGlobalConf error %v", err)
		return err
	}
	sConf := &gConf.Conf

	// first time: register node
	stateFile := filepath.Join(consts.TarsPath, "data", consts.CheckStatusFile)
	st, _ := os.Stat(stateFile)

	lastState := ""
	if bs, err := ioutil.ReadFile(stateFile); err == nil {
		lastState = string(bs)
	}
	state := "inactive"
	if checkSucc {
		state = "active"
	}

	if st != nil && st.ModTime().Add(c.reportStatusInterval).After(time.Now()) {
		if lastState == state {
			return nil
		}
	}

	if err := ioutil.WriteFile(stateFile, []byte(state), 0644); err != nil {
		log.Debugf("WriteFile error %v", err)
	}

	// invoke tars registry and set server status
	client := tarsproxy.GetRegistryClient(sConf.Locator)
	if client == nil {
		log.Debugf("GetRegistryClient return nil, Locator %s", sConf.Locator)
		return nil
	}
	req := Tars.KeepAliveReq{
		NodeName: consts.LocalIP,
		State:    state,
		Application: sConf.Application,
		Server: sConf.Server,
		SetID: sConf.SetID,
	}
	if err := client.KeepAlive(context.Background(), &req); err != nil {
		log.Debugf("KeepAlive error %v", err)
		return nil
	}
	return nil
}

func registerNode(sConf *genconf.ServerConf, disableFlow bool) error {
	// invoke tars registry and register the endponts
	client := tarsproxy.GetRegistryClient(sConf.Locator)
	if client == nil {
		return fmt.Errorf("get client failed")
	}
	adapters := make([]Tars.AdapterConf, len(sConf.Adapters))
	for i := range sConf.Adapters {
		sv := sConf.Adapters[i]
		adapters[i] = Tars.AdapterConf{
			Servant:      fmt.Sprintf("%s.%s.%s", sConf.Application, sConf.Server, sv.Object),
			Endpoint:     sv.Endpoint,
			Protocol:     sv.Protocol,
			MaxConns:     int32(sv.MaxConns),
			QueueCap:     int32(sv.QueueCap),
			QueueTimeout: int32(sv.QueueTimeout),
			ThreadNum:    int32(sv.Threads),
		}
	}
	req := Tars.OnStartupReq{
		Application: sConf.Application,
		Server:      sConf.Server,
		SetID:       sConf.SetID,
		NodeName:    consts.LocalIP,
		Adapters:    adapters,
		DisableFlow: disableFlow,
		State:       "activating",
		Version:     os.Getenv("SERVER_VERSION"),
	}
	return client.OnStartup(context.Background(), &req)
}

func noitfyMsg(sConf *genconf.ServerConf, msg string) error {
	// invoke tars registry and register the endponts
	client := tarsproxy.GetNotifyClient(sConf.Locator)
	req := &notifyf.ReportInfo{
		EType:    notifyf.ReportType_REPORT,
		SApp:     sConf.Application,
		SServer:  sConf.Server,
		SSet:     sConf.SetID,
		SMessage: msg,
		ELevel:   notifyf.NOTIFYLEVEL_NOTIFYERROR,
	}
	return client.ReportNotify(context.Background(), req)
}
