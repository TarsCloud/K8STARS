package supervisor

import (
	"fmt"
	"net"
	"testing"
)

func TestCheckAddr(t *testing.T) {

	fmt.Println("begin dial...")
	conn, err := net.Dial("tcp", "localhost:8888")
	if opErr, ok := err.(*net.OpError); ok {
		fmt.Printf("op err %T %+v %v\n", opErr.Err, opErr.Err, opErr.Timeout())
	}
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()
	fmt.Println("dial ok")

	err = checkAddr("tcp", "baidu.com:80")
	fmt.Println("1", err)
	err = checkAddr("tcp", "google.com:80")
	fmt.Println("2", err)
	err = checkAddr("tcp", "34.55.66.6:8099")
	if opErr, ok := err.(*net.OpError); ok {
		fmt.Printf("op err %T %+v\n", opErr.Err, opErr.Err)
	}
}
