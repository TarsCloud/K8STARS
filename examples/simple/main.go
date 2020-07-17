package main

import (
	"fmt"
	"net/http"

	"github.com/tarscloud/k8stars/examples/simple/autogen/App"
	"github.com/TarsCloud/TarsGo/tars"
)

var (
	comm = tars.NewCommunicator()
)

func main() {
	app := &App.SimpleServer{}
	imp := &simpleServerImp{}
	cfg := tars.GetServerConfig()
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".MainObj")

	mux := &tars.TarsHttpMux{}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		client := &App.SimpleServer{}
		comm.StringToProxy(cfg.App+"."+cfg.Server+".MainObj", client)
		ret, err := client.Sum(1, 2)
		if err != nil {
			fmt.Fprintf(w, "invoke error %v", err)
			return
		}
		fmt.Fprintf(w, "sum is %d", ret)
	})
	tars.AddHttpServant(mux, cfg.App+"."+cfg.Server+".HttpObj")
	tars.Run()
}
