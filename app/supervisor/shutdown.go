package supervisor

import (
	"time"

	"github.com/tarscloud/k8stars/app/genconf"
	"github.com/tarscloud/k8stars/app/prestop"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/adminf"
)

func (c *launchCmd) shutdown() error {
	// get config from file
	prestop.Prestop(c.waitStopTime)

	// get config from file
	gConf, err := genconf.GetGlobalConf()
	if err != nil {
		return err
	}
	sConf := &gConf.Conf

	// notify shutdown
	adminClient := &adminf.AdminF{}
	comm := tars.NewCommunicator()
	comm.StringToProxy(sConf.Application+"."+sConf.Server+".adminObj@"+sConf.LocalEndpoint, adminClient)
	go adminClient.Shutdown()

	deadline := time.Now().Add(time.Second * 60)
	for range time.NewTicker(time.Second).C {
		if err := CheckServerStatus(sConf); err != nil {
			break
		}
		if time.Now().After(deadline) {
			log.Debugf("shutdown anyway")
			break
		}
	}
	return nil
}
