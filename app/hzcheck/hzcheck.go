package hzcheck

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/tarscloud/k8stars/consts"
	"github.com/tarscloud/k8stars/logger"
	"github.com/tarscloud/k8stars/tinycli"
)

var (
	log = logger.GetLogger()
)

type hzCheckCmd struct {
	updateInterval time.Duration
}

// NewCmd returns an instances of launchCmd
func NewCmd() tinycli.Cmd {
	return &hzCheckCmd{}
}

// InitFlag initializes options from environment variables
func (c *hzCheckCmd) InitFlag(setter tinycli.EnvFlagSetter) {
	setter.SetDuration("TARS_REPORT_INTERVAL", &c.updateInterval, "30s", "Time interval of checking status")
}

func (c *hzCheckCmd) Start() error {
	// check last report status and time
	stateFile := filepath.Join(consts.TarsPath, "data", consts.CheckStatusFile)
	st, err := os.Stat(stateFile)
	if err != nil {
		return fmt.Errorf("not ready")
	}
	if st.ModTime().Add(c.updateInterval * 2).Before(time.Now()) {
		return fmt.Errorf("supervisor not alive")
	}
	bs, err := ioutil.ReadFile(stateFile)

	// maybe concurrent read and write, try again
	for tryCount := 10; tryCount > 0; tryCount-- {
		if err != nil || string(bs) == "" {
			time.Sleep(time.Millisecond * 50)
			bs, err = ioutil.ReadFile(stateFile)
		} else {
			break
		}
	}
	if string(bs) == consts.StateActive || string(bs) == consts.StateDeactivating {
		return nil
	}
	return fmt.Errorf("not active:%s", string(bs))
}
