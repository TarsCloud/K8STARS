package logger

import (
	"os"
	"path/filepath"

	"github.com/tarscloud/k8stars/consts"
	"github.com/TarsCloud/TarsGo/tars/util/rogger"
)

var log *rogger.Logger

// GetLogger returns logger
func GetLogger() *rogger.Logger {
	if log != nil {
		return log
	}
	log = rogger.GetLogger("cmd")
	if _, err := os.Stat(consts.TarsPath); os.IsNotExist(err) {
		return log
	}
	fpath := filepath.Join(consts.TarsPath, "log")
	log.SetFileRoller(fpath, 10, 100)
	return log
}
