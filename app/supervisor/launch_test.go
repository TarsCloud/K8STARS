package supervisor

import (
	"os"
	"testing"
)

func TestLaunch(t *testing.T) {
	os.Setenv("TARS_PATH", "/tmp/launch_test/")
}
