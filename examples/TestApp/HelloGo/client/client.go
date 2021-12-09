package main

import (
	"fmt"

	"github.com/TarsCloud/TarsGo/tars"
	"github.com/tarscloud/k8stars/examples/TestApp/HelloGo/TestApp"
)

func main() {
	comm := tars.NewCommunicator()
	obj := fmt.Sprintf("TestApp.HelloGo.SayHelloObj")
	app := new(TestApp.SayHello)
	//cfg := tars.GetServerConfig()

	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 192.168.0.55 -p 17890")

	comm.StringToProxy(obj, app)
	var out, i int32
	i = 123
	ret, err := app.Add(i, i*2, &out)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret, out)
}
