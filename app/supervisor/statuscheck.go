package supervisor

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/tarscloud/k8stars/app/genconf"
)

var connTimeout = time.Millisecond * 300

// CheckServerStatus checks the status of server
func CheckServerStatus(sConf *genconf.ServerConf) error {
	sList := make([]string, 0)
	for _, sv := range sConf.Adapters {
		sList = append(sList, sv.Endpoint)
	}
	sList = append(sList, sConf.LocalEndpoint)
	for _, sv := range sList {
		network, host, port := "", "", ""
		arr := strings.Fields(sv)
		for i := range arr {
			if i == 0 {
				network = arr[0]
			} else if arr[i] == "-h" && i+1 < len(arr) {
				host = arr[i+1]
			} else if arr[i] == "-p" && i+1 < len(arr) {
				port = arr[i+1]
			}
		}
		if network == "" || host == "" || port == "" {
			return fmt.Errorf("host/port or network can not be empty")
		}
		if err := checkAddr(network, host+":"+port); err != nil {
			return fmt.Errorf("Connect failed %s:%v", sv, err)
		}
	}
	return nil
}

func checkAddr(network, addr string) error {
	conn, err := net.DialTimeout(network, addr, connTimeout)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}
