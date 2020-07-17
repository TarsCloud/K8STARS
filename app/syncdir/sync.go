package syncdir

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tarscloud/k8stars/consts"
	"github.com/tarscloud/k8stars/logger"
	"github.com/tarscloud/k8stars/tinycli"
)

var (
	log = logger.GetLogger()
)

type syncDirCmd struct {
	src, dst string
}

// NewCmd returns an instances of syncDirCmd
func NewCmd() tinycli.Cmd {
	return &syncDirCmd{
		dst: filepath.Join(consts.TarsPath, "bin"),
	}
}

// InitFlag initializes options from environment variables
func (c *syncDirCmd) InitFlag(setter tinycli.EnvFlagSetter) {
	setter.SetString("TARS_SYNC_DIRECTORY", &c.src, c.src, "Source directory for syncing files")
	setter.SetString("TARS_SYNC_TARGET_DIRECTORY", &c.dst, c.dst, "Destination directory for syncing files")
}

// Start starts the command
func (c *syncDirCmd) Start() error {
	if c.src == "" {
		return nil
	}
	src, dst := c.src, c.dst
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	os.MkdirAll(dst, 0755)

	for _, f := range files {
		srcFile := filepath.Join(src, f.Name())
		if !f.Mode().IsRegular() {
			if srcFile, err = os.Readlink(srcFile); err != nil {
				// log.Debugf("skip %s", f.Name()) // skip not regular file
				continue
			}
			srcFile = filepath.Clean(filepath.Join(src, srcFile))
			if f, err = os.Stat(srcFile); f == nil || !f.Mode().IsRegular() {
				// log.Debugf("skip %s", srcFile) // skip not regular file
				continue
			}
		}
		dstFile := filepath.Join(dst, f.Name())
		st, _ := os.Stat(dstFile)
		if st != nil && st.Size() == f.Size() && st.ModTime() == f.ModTime() {
			// log.Debugf("file not change %s", f.Name())
			continue
		}
		log.Debugf("sync file %s", f.Name())
		bs, err := ioutil.ReadFile(srcFile)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(dstFile, bs, 0755); err != nil {
			return err
		}
		if err := os.Chtimes(dstFile, f.ModTime(), f.ModTime()); err != nil {
			return err
		}
	}
	return nil
}
