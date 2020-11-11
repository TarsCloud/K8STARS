package prestop

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/adminf"
	"github.com/tarscloud/k8stars/app/genconf"
	"github.com/tarscloud/k8stars/consts"
	"github.com/tarscloud/k8stars/logger"
	"github.com/tarscloud/k8stars/tarsproxy"
	Tars "github.com/tarscloud/k8stars/tarsregistry/autogen/tars"
	"github.com/tarscloud/k8stars/tinycli"
)

var (
	log = logger.GetLogger()
)

type prestopCmd struct {
	waitStopTime time.Duration
}

// NewCmd returns an instances of prestopCmd
func NewCmd() tinycli.Cmd {
	return &prestopCmd{}
}

// InitFlag initializes options from environment variables
func (c *prestopCmd) InitFlag(setter tinycli.EnvFlagSetter) {
	setter.SetDuration("TARS_PRESTOP_WAITTIME", &c.waitStopTime, "80s", "Wait time before stop")
}

// Start starts the command
func (c *prestopCmd) Start() error {
	return Prestop(c.waitStopTime)
}

// Prestop notifies registry and waits
func Prestop(waitStopTime time.Duration) error {

	preStopFile := filepath.Join(consts.TarsPath, "data", "prestop")
	if _, err := os.Stat(preStopFile); !os.IsNotExist(err) {
		log.Info("Prestop has run before")
		return nil
	}
	ioutil.WriteFile(preStopFile, []byte("done"), 0644)

	// set state to file
	stateFile := filepath.Join(consts.TarsPath, "data", consts.CheckStatusFile)
	ioutil.WriteFile(stateFile, []byte(consts.StateDeactivating), 0644)
	// wait for keeplive
	time.Sleep(time.Second)

	// get config from file
	gConf, err := genconf.GetGlobalConf()
	if err != nil {
		return err
	}
	sConf := &gConf.Conf

	// invoke tars registry
	client := tarsproxy.GetRegistryClient(sConf.Locator)
	if client == nil {
		log.Debugf("GetRegistryClient return nil, Locator %s", sConf.Locator)
		return nil
	}

	// set deactivating
	req1 := &Tars.KeepAliveReq{
		NodeName:    gConf.LocalIP,
		Server:      sConf.Server,
		Application: sConf.Application,
		State:       consts.StateDeactivating,
	}
	err = client.KeepAlive(context.Background(), req1)
	if err != nil {
		log.Debugf("KeepAlive error %v", err)
	}

	// notify
	adminClient := &adminf.AdminF{}
	comm := tars.NewCommunicator()
	comm.StringToProxy("AdminObj@"+sConf.LocalEndpoint, adminClient)
	adminClient.TarsSetTimeout(1000)
	adminClient.Notify("prestop")

	// wait
	if waitStopTime > 0 {
		time.Sleep(waitStopTime)
	}

	// prestop
	req := &Tars.OnPrestopReq{
		NodeName:    gConf.LocalIP,
		Application: sConf.Application,
		Server:      sConf.Server,
	}
	err = client.OnPrestop(context.Background(), req)
	if err != nil {
		log.Debugf("Prestop error %v", err)
	}

	return nil
}
