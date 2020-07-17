package supervisor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/tarscloud/k8stars/algorithm/recentuse"
	"github.com/tarscloud/k8stars/app/genconf"
	"github.com/tarscloud/k8stars/consts"
	"github.com/tarscloud/k8stars/logger"
	"github.com/tarscloud/k8stars/tinycli"
)

var (
	log = logger.GetLogger()
)

type launchCmd struct {
	ctx           context.Context
	cancelFunc    context.CancelFunc
	recentRestart *recentuse.RecentUse
	startPath     string
	stopPath      string

	checkIntv time.Duration

	beforeCheckScript    string
	checkScriptTimeout   time.Duration
	reportStatusInterval time.Duration
	waitStopTime         time.Duration
	disableFlow          string
	onExitChan           chan bool
}

// NewCmd returns an instances of launchCmd
func NewCmd() tinycli.Cmd {
	c := &launchCmd{
		onExitChan:    make(chan bool, 10),
		recentRestart: recentuse.NewRecentUse(time.Minute * 5),
	}
	c.ctx, c.cancelFunc = context.WithCancel(context.Background())
	return c
}

// InitFlag initializes options from environment variables
func (c *launchCmd) InitFlag(setter tinycli.EnvFlagSetter) {
	sp := filepath.Join(consts.TarsPath, "bin", "start.sh")
	setter.SetString("TARS_START_PATH", &c.startPath, sp, "Path of start script")
	setter.SetString("TARS_STOP_PATH", &c.stopPath, "", "Path of stop script")
	setter.SetDuration("TARS_REPORT_INTERVAL", &c.reportStatusInterval, "30s", "Time interval of reporting state")
	setter.SetString("TARS_DISABLE_FLOW", &c.disableFlow, "", "None empty string to turn off the flow")
	setter.SetDuration("TARS_CHECK_INTERVAL", &c.checkIntv, "10s", "Time interval of checking status")
	setter.SetString("TARS_BEFORE_CHECK_SCRIPT", &c.beforeCheckScript, "", "Run script before check")
	setter.SetDuration("TARS_CHECK_SCRIPT_TIMEOUT", &c.checkScriptTimeout, "2s", "Max running time of script")
	setter.SetDuration("TARS_PRESTOP_WAITTIME", &c.waitStopTime, "80s", "Wait time before stop")
}

func (c *launchCmd) checkStatus() bool {
	// get config from file
	gConf, err := genconf.GetGlobalConf()
	if err != nil {
		return false
	}
	sConf := &gConf.Conf

	err = CheckServerStatus(sConf)
	if err != nil {
		log.Errorf("Check error %v", err)
		return false
	}
	return true
}

func (c *launchCmd) restartAndNotify() {
	c.restartSever()
	// get config from file
	gConf, err := genconf.GetGlobalConf()
	if err != nil {
		return
	}
	sConf := &gConf.Conf
	noitfyMsg(sConf, "[alarm] down, server is inactive")
}

func (c *launchCmd) restartSever() {
	stopPath := filepath.Join(consts.TarsPath, "data", "stop")
	if _, err := os.Stat(stopPath); !os.IsNotExist(err) {
		log.Debug("skip restart")
		return
	}
	// kill first
	binPath := filepath.Join(consts.TarsPath, "bin")
	stopCmd := c.stopPath
	if stopCmd == "" {
		stopCmd = fmt.Sprintf("ps -ef | grep '%s' | grep -v grep |  awk '{print $2}' | xargs kill -9", binPath)
	}
	log.Debug(stopCmd)
	cmd := exec.Command("sh", "-c", stopCmd)
	cmd.Run()
	time.Sleep(time.Second * 1)

	go c.startSever()

}

func (c *launchCmd) startSever() {
	cmd := exec.Command("sh", "-c", c.startPath)
	p := filepath.Join(consts.TarsPath, "log", "start.log")
	outfile, err := os.Create(p)
	if err == nil {
		cmd.Stderr = outfile
		cmd.Stdout = outfile
	}

	log.Debugf("start server %s", c.startPath)
	go func() {
		err := cmd.Run()
		if err != nil {
			log.Errorf("server stop error %v", err)
		}
		if outfile != nil {
			outfile.Close()
		}
		c.onExitChan <- true
	}()
}

func (c *launchCmd) preCheck() {
	if c.beforeCheckScript == "" {
		return
	}
	cmd := exec.Command("sh", "-c", c.beforeCheckScript)
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.checkScriptTimeout)
	go func() {
		cmd.Run()
		cancelFunc()
	}()
	<-ctx.Done()
}

func (c *launchCmd) check() {
	go reapProcess()
	shutdownChan := waitShutdown()
	checkTk := time.NewTicker(c.checkIntv)

	isStop := false
	for {
		select {
		case <-checkTk.C:
			go c.preCheck()
			beActive := c.checkStatus()
			log.Debugf("check status ret: %v", beActive)
			c.keepAlive(beActive)
			if isStop {
				if !c.recentRestart.KeepAlive("restart") {
					c.restartAndNotify()
					isStop = false
				}
			}
		case <-c.onExitChan:
			isStop = true
			if !c.recentRestart.KeepAlive("restart") {
				c.keepAlive(false)
				c.restartAndNotify()
				isStop = false
				continue
			}
		case <-shutdownChan:
			c.shutdown()
			c.cancelFunc()
			return
		}
	}
}

func (c *launchCmd) registerServer() error {
	// get config from file
	gConf, err := genconf.GetGlobalConf()
	if err != nil {
		return err
	}
	sConf := &gConf.Conf
	err = registerNode(sConf, c.disableFlow != "")
	if err != nil {
		log.Errorf("RegisterNode error %v", err)
	} else {
		log.Debug("registerNode succ")
	}
	return err
}

// Start run the command
func (c *launchCmd) Start() error {
	// redirect stderr/stderr to supervisor.log
	p := filepath.Join(consts.TarsPath, "log", "supervisor.log")
	outfile, err := os.Create(p)
	if err == nil {
		os.Stderr = outfile
		os.Stdout = outfile
	}
	defer func() {
		if outfile != nil {
			outfile.Close()
		}
	}()

	genconf := &tinycli.App{
		Cmd: genconf.NewCmd(),
	}
	if err := genconf.Run(os.Args[1:], os.Environ()); err != nil {
		return err
	}
	for {
		if err := c.registerServer(); err == nil {
			break
		}
		time.Sleep(time.Minute)
	}
	c.preCheck()
	c.startSever()
	go c.check()
	<-c.ctx.Done()
	return nil
}
