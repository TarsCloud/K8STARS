package prestop

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/tarscloud/k8stars/app/genconf"
	"github.com/tarscloud/k8stars/consts"
	"github.com/tarscloud/k8stars/logger"
	"github.com/tarscloud/k8stars/tarsproxy"
	"github.com/tarscloud/k8stars/tarsregistry/autogen/Tars"
	"github.com/tarscloud/k8stars/tinycli"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/adminf"
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
	req := &Tars.OnPrestopReq{NodeName: consts.LocalIP}
	err = client.OnPrestop(context.Background(), req)
	if err != nil {
		log.Debugf("Prestop error %v", err)
	}
	// notify prestop
	adminClient := &adminf.AdminF{}
	comm := tars.NewCommunicator()
	comm.StringToProxy("AdminObj@"+sConf.LocalEndpoint, adminClient)
	adminClient.TarsSetTimeout(1000)
	adminClient.Notify("prestop")

	if waitStopTime > 0 {
		time.Sleep(waitStopTime)
	}
	return nil
}
