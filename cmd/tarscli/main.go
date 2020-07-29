package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars/util/rogger"
	"github.com/tarscloud/k8stars/app/genconf"
	"github.com/tarscloud/k8stars/app/hzcheck"
	"github.com/tarscloud/k8stars/app/notify"
	"github.com/tarscloud/k8stars/app/prestop"
	"github.com/tarscloud/k8stars/app/supervisor"
	"github.com/tarscloud/k8stars/app/syncdir"
	"github.com/tarscloud/k8stars/tinycli"
)

var VERSION string = "unstable"

func main() {
	if len(os.Args) == 2 && os.Args[1] == "version" {
		fmt.Println(VERSION)
		return
	}

	rogger.SetLevel(rogger.DEBUG)
	defer rogger.FlushLogger()
	apps := map[string]*tinycli.App{
		"supervisor": {
			Name:  "supervisor",
			Usage: "Supervisor is a launcher for tars server",
			Cmd:   supervisor.NewCmd(),
		},
		"genconf": {
			Name:  "genconf",
			Usage: "Genconf is a tool to generate tars config file",
			Cmd:   genconf.NewCmd(),
		},
		"hzcheck": {
			Name:  "hzcheck",
			Usage: "HZcheck is a readiness probe for tars server",
			Cmd:   hzcheck.NewCmd(),
		},
		"prestop": {
			Name:  "prestop",
			Usage: "Prestop is a pre-stop script for tars server",
			Cmd:   prestop.NewCmd(),
		},
		"syncdir": {
			Name:  "syncdir",
			Usage: "syncdir is a command for syncing files",
			Cmd:   syncdir.NewCmd(),
		},
		"notify": {
			Name:  "notify",
			Usage: "notify is a command for notifying tars server",
			Cmd:   notify.NewCmd(os.Args[1:]),
		},
	}
	if len(os.Args) <= 1 || os.Args[1] == "help" || os.Args[1] == "-h" {
		for k, app := range apps {
			fmt.Printf("------ %s -----\n", app.Name)
			app.Run([]string{k, "help"}, os.Environ())
		}
		return
	}
	app, ok := apps[os.Args[1]]
	if !ok {
		fmt.Println("command not found", os.Args[1])
		return
	}

	if err := app.Run(os.Args[1:], os.Environ()); err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		rogger.FlushLogger()
		os.Exit(1)
	}
}
