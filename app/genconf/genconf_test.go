package genconf

import (
	"os"
	"testing"

	"github.com/tarscloud/k8stars/consts"
	"github.com/tarscloud/k8stars/tinycli"
	"github.com/stretchr/testify/assert"
)

func TestLaunchCmd(t *testing.T) {
	consts.TarsPath = "testtarsdir"
	err := os.MkdirAll(consts.TarsPath+"/bin/../conf/../data", 0755)
	assert.Nil(t, err)

	app := &tinycli.App{
		Name:  "launch",
		Usage: "Launch is an launcher for tars server",
		Cmd:   NewCmd(),
	}
	envs := []string{"TARS_MERGE_CONF=testtarsdir/bin/_for_merge.conf"}
	err = app.Run(os.Args, envs)
	assert.Nil(t, err)
}
