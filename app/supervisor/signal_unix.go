// +build linux

package supervisor

import (
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

func waitShutdown() <-chan bool {
	c := make(chan bool, 100)
	go onShutdown(c)
	return c
}

func onShutdown(c chan bool) {
	cc := make(chan os.Signal, 1)
	signal.Notify(cc, unix.SIGCHLD)
	for range c {
		c <- true
	}
}

func reapProcess() {
	cc := make(chan os.Signal, 1)
	signal.Notify(cc, unix.SIGCHLD)
	for range cc {
		func() {
		POLL:
			var status unix.WaitStatus
			pid, err := unix.Wait4(-1, &status, unix.WNOHANG, nil)
			switch err {
			case nil:
				if pid > 0 {
					goto POLL
				}
				return
			case unix.ECHILD:
				return
			case unix.EINTR:
				goto POLL
			default:
				return
			}
		}()
	}
}
