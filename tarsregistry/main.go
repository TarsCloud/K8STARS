package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"
	mtars "github.com/tarscloud/k8stars/tarsregistry/autogen/Tars"
	"github.com/tarscloud/k8stars/tarsregistry/store"
)

var (
	logger = tars.GetLogger("registry")
)

func main() {
	cfg := tars.GetServerConfig()
	app := &mtars.Tarsregistry{}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:pass@tcp(127.0.0.1:3306)/db_tars"
	}
	logger.Infof("dsn %s", dsn)
	driver, err := store.NewMysqlDB(dsn)
	if err != nil {
		panic(err)
	}
	imp := &registryImp{
		driver: driver,
	}
	startReq := &mtars.OnStartupReq{
		NodeName:    cfg.LocalIP,
		Application: cfg.App,
		Server:      cfg.Server,
		SetID:       cfg.Setdivision,
		State:       "active",
		Version:     os.Getenv("SERVER_VERSION"),
	}
	obj := cfg.App + "." + cfg.Server + ".Registry"
	for _, v := range cfg.Adapters {
		if v.Obj == obj {
			ep := fmt.Sprintf("%s -h %s -p %d", v.Endpoint.Proto, v.Endpoint.Host, v.Endpoint.Port)
			startReq.Adapters = []mtars.AdapterConf{
				{
					Servant:      obj,
					Endpoint:     ep,
					Protocol:     "tars",
					MaxConns:     1000,
					ThreadNum:    10,
					QueueCap:     10000,
					QueueTimeout: 60000,
				},
			}
		}
	}
	if len(startReq.Adapters) == 0 {
		panic("Registry object not config")
	}

	go imp.keepAlive(startReq)

	app.AddServantWithContext(imp, obj)
	tars.Run()
}
