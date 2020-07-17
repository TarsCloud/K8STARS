package genconf

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	allocated = make(map[string]bool)
)

func getRandomPort(endpoint string, start, end int) (string, error) {
	var host, network string
	fs := strings.Fields(endpoint)
	for i := range fs {
		if i == 0 {
			network = fs[i]
		} else if fs[i] == "-h" && i < len(fs)-1 {
			host = fs[i+1]
		}
	}
	if host == "" || network == "" {
		return "", fmt.Errorf("host or network can not be empty")
	}

	for i := start; i <= end; i++ {
		si := strconv.Itoa(i)
		if allocated[si] {
			continue
		}
		addr := fmt.Sprintf("%s:%s", host, si)
		if portIsAvailable(network, addr) {
			allocated[si] = true
			return si, nil
		}
	}
	return "", fmt.Errorf("no available port")
}

func portIsAvailable(network, addr string) bool {
	ln, err := net.Listen(network, addr)
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}
