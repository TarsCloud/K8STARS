package tinycli

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testApp struct {
	started bool

	intv     int
	intv1    int
	intv2    int
	stringv  string
	stringv1 string
	stringv2 string
	timev    time.Duration
	timev1   time.Duration
	timev2   time.Duration
}

func (t *testApp) Start() error {
	t.started = true
	return nil
}
func (t *testApp) InitFlag(st EnvFlagSetter) {
	st.SetInt("INT_X", &t.intv, 1, "int val")
	st.SetInt("INT_VAL", &t.intv1, 1, "int val")
	st.SetInt("INT", &t.intv2, 0, "int val")
	st.SetString("STRING", &t.stringv, "", "string val")
	st.SetDuration("TIME", &t.timev, "", "time val")
	st.SetString("STRING1", &t.stringv1, "", "string val")
	st.SetString("STRING2", &t.stringv2, "xxx", "string val")
	st.SetDuration("TIME1", &t.timev1, "", "time val")
	st.SetDuration("TIME2", &t.timev2, "10ms", "time val")
}

func TestAppNormal(t *testing.T) {
	tp := &testApp{}
	envs := []string{"INT=2", "STRING=XXX", "TIME=100ms"}
	app := App{
		Name:  "test",
		Usage: "usage",
	}
	err := app.Run(nil, envs)
	assert.NotNil(t, err)

	app.Cmd = tp
	err = app.Run(nil, envs)
	assert.Nil(t, err)

	assert.Equal(t, 1, tp.intv1)
	assert.Equal(t, 2, tp.intv2)
	assert.Equal(t, "XXX", tp.stringv)
	assert.Equal(t, 100*time.Millisecond, tp.timev)
}

func TestAppHelp(t *testing.T) {
	tp := &testApp{}
	envs := []string{"INT_X=1mm", "TIME2=20mm"}
	app := App{
		Name:  "test",
		Usage: "usage",
		Cmd:   tp,
	}
	tf, err := ioutil.TempFile("", "")
	assert.Nil(t, err)
	os.Stdout = tf
	app.Run([]string{"cmd", "help"}, envs)

	app = App{
		Name:  "test",
		Usage: "usage",
		Cmd:   tp,
	}
	err = app.Run(nil, envs)
	assert.NotNil(t, err)
	assert.Equal(t, tp.timev1, time.Duration(0))
}
