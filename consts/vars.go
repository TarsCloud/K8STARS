package consts

import (
	"net"
	"os"
	"strconv"
)

// ServerMetaFile is filename of server meta file
const ServerMetaFile = "_server_meta.yaml"

// ServerInfoFile is filename of server info file
const ServerInfoFile = "server_info.yaml"

// LastCheckTime is the last time of checking server status
const CheckStatusFile = "check_status"

var (
	// LocalIP is the ipv4 address
	LocalIP string
	// NameSpace is the namepace of cluster
	NameSpace = "tars-system"
	// TarsPath is the work directory for tars server
	TarsPath = "/tars"
	// RandPortMin is min random port
	RandPortMin = 13000
	// RandPortMax is max random port
	RandPortMax = 16000
)

func init() {
	if e := os.Getenv("POD_IP"); e != "" {
		LocalIP = e
	} else {
		LocalIP = getIP()
	}
	if e := os.Getenv("TARS_NAMESPACE"); e != "" {
		NameSpace = e
	}
	if e := os.Getenv("TARS_PATH"); e != "" {
		TarsPath = e
	}
	if e := os.Getenv("TARS_RAND_PORT_MIN"); e != "" {
		if port, err := strconv.Atoi(e); err == nil {
			RandPortMin = port
		}
	}
	if e := os.Getenv("TARS_RAND_PORT_MAX"); e != "" {
		if port, err := strconv.Atoi(e); err == nil {
			RandPortMax = port
		}
	}
}

func getIP() string {
	conn, err := net.Dial("udp", "10.0.0.1:80")
	if err != nil {
		return "0.0.0.0"
	}
	defer conn.Close()

	addr := conn.LocalAddr().(*net.UDPAddr)
	if addr == nil {
		return "0.0.0.0"
	}
	return addr.IP.String()
}
