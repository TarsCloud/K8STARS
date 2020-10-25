package notify

import (
	"fmt"
	"strings"

	"github.com/tarscloud/k8stars/app/genconf"
	"github.com/tarscloud/k8stars/tinycli"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/adminf"
	"github.com/TarsCloud/TarsGo/tars/util/rogger"
)

type notifyCmd struct {
	cmd string
}

// NewCmd returns an instances of notifyCmd
func NewCmd(args []string) tinycli.Cmd {
	cmd := strings.Join(args[1:], " ")
	return &notifyCmd{cmd: cmd}
}

// InitFlag initializes options from environment variables
func (c *notifyCmd) InitFlag(setter tinycli.EnvFlagSetter) {
}

// Start starts the command
func (c *notifyCmd) Start() error {
	// get config from file
	rogger.SetLevel(rogger.OFF)

	gConf, err := genconf.GetGlobalConf()
	if err != nil {
		return err
	}
	sConf := &gConf.Conf

	// notify
	adminClient := &adminf.AdminF{}
	comm := tars.NewCommunicator()
	comm.StringToProxy("AdminObj@"+sConf.LocalEndpoint, adminClient)
	adminClient.TarsSetTimeout(5000)
	ret, err := adminClient.Notify(c.cmd)
	if err != nil {
		return err
	}
	fmt.Printf("notify succ:%s\n", ret)
	return nil
}
