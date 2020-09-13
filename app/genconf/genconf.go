package genconf

import (
	"fmt"
	"path/filepath"

	"github.com/tarscloud/k8stars/consts"
	"github.com/tarscloud/k8stars/logger"
	"github.com/tarscloud/k8stars/tinycli"
)

var (
	log = logger.GetLogger()
)

type genConfCmd struct {
	svrConf     ServerConf
	startScript string
	buildServer string
	mergeConf   string
}

// NewCmd returns an instances of genConfCmd
func NewCmd() tinycli.Cmd {
	tarsPath := consts.TarsPath
	// default
	svrConf := defaultServerConf()
	confPath := filepath.Join(tarsPath, "bin", consts.ServerMetaFile)
	if err := parseServerConf(confPath, &svrConf); err != nil {
		// log.Errorf("ParseServerConf error %v", err)
	}
	return &genConfCmd{svrConf: svrConf}
}

// InitFlag initializes options from environment variables
func (c *genConfCmd) InitFlag(setter tinycli.EnvFlagSetter) {
	setter.SetString("TARS_APPLICATION", &c.svrConf.Application, c.svrConf.Application, "Application name")
	setter.SetString("TARS_SERVER", &c.svrConf.Server, c.svrConf.Server, "Server name of the running tars server")
	setter.SetString("TARS_BUILD_SERVER", &c.buildServer, c.buildServer, "Server name when compiled and built")
	setter.SetString("TARS_SET_ID", &c.svrConf.SetID, c.svrConf.SetID, "SetID of server")
	setter.SetString("TARS_LOCATOR", &c.svrConf.Locator, c.svrConf.Locator, "Object and endpoint of locator")
	setter.SetString("TARS_MERGE_CONF", &c.mergeConf, c.mergeConf, "Config file path to merge to config")
	setter.SetInt("TARS_EP_UPTIME", &c.svrConf.RefreshEndpointInterval, c.svrConf.RefreshEndpointInterval, "Interval(seconds) of refreshing endpoints")
}

// Start run command in background
func (c *genConfCmd) Start() error {
	// parse server config file
	sConf := &c.svrConf
	if sConf.Application == "" || sConf.Server == "" {
		return fmt.Errorf("application or server name can not be empty")
	}

	// generate server conf
	log.Debug("start")
	if c.buildServer != "" {
		log.Debugf("found build server %v", c.buildServer)
	} else {
		c.buildServer = sConf.Server
	}
	return generateConf(sConf, c.buildServer, c.mergeConf)
}
