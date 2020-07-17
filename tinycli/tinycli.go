package tinycli

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Cmd interface {
	Start() error
	InitFlag(EnvFlagSetter)
}

type EnvFlagSetter interface {
	SetInt(string, *int, int, string)
	SetDuration(string, *time.Duration, string, string)
	SetString(string, *string, string, string)
}

type App struct {
	Name  string
	Usage string
	Cmd   Cmd
}

type envFlagSetter struct {
	envs    map[string]string
	envList string
	err     error
}

func (a *App) Run(args, envs []string) error {
	if a.Cmd == nil {
		return fmt.Errorf("cmd is nil")
	}
	efs := &envFlagSetter{envs: make(map[string]string)}
	for i := range envs {
		arr := strings.SplitN(envs[i], "=", 2)
		if len(arr) == 2 {
			efs.envs[arr[0]] = arr[1]
		}
	}
	a.Cmd.InitFlag(efs)
	if len(args) > 1 {
		if args[1] == "help" {
			text := fmt.Sprintf(`%s

Usage:
  %s [help]

Available environment variables:
%s
`, a.Usage, a.Name, efs.envList)
			os.Stdout.Write([]byte(text))
			return nil
		}
	}
	if efs.err != nil {
		return efs.err
	}
	return a.Cmd.Start()
}

func (e *envFlagSetter) SetInt(key string, val *int, defaultVal int, usage string) {
	var defaultStr string
	if defaultVal != 0 {
		defaultStr = fmt.Sprintf(" (default %d)", defaultVal)
	}
	e.envList += fmt.Sprintf(`%s int    %s.%s
`, key, usage, defaultStr)
	exists, ev := e.getEnv(key)
	if err := setFlagInt(exists, ev, val, defaultVal, usage); err != nil {
		e.err = err
	}
}

func (e *envFlagSetter) SetString(key string, val *string, defaultVal string, usage string) {
	var defaultStr string
	if defaultVal != "" {
		defaultStr = fmt.Sprintf(" (default %s)", defaultVal)
	}
	e.envList += fmt.Sprintf(`%s string    %s.%s
`, key, usage, defaultStr)
	exists, ev := e.getEnv(key)
	_ = setFlagString(exists, ev, val, defaultVal, usage)
}

func (e *envFlagSetter) SetDuration(key string, val *time.Duration, defaultVal string, usage string) {
	var defaultStr string
	if defaultVal != "" {
		defaultStr = fmt.Sprintf(" (default %s)", defaultVal)
	}
	e.envList += fmt.Sprintf(`%s time    %s.%s
`, key, usage, defaultStr)
	exists, ev := e.getEnv(key)
	if err := setFlagDuration(exists, ev, val, defaultVal, usage); err != nil {
		e.err = err
	}
}

func (e *envFlagSetter) getEnv(key string) (bool, string) {
	if v, ok := e.envs[key]; ok {
		return true, v
	}
	return false, ""
}
