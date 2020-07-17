// +build !linux

package supervisor

func reapProcess() {
}

func waitShutdown() <-chan bool {
	c := make(chan bool, 100)
	return c
}
