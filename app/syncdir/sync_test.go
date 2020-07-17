package syncdir

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncDir(t *testing.T) {
	src := os.TempDir() + "/TestSyncDirSrc"
	dst := os.TempDir() + "/TestSyncDirDst"
	cmd := syncDirCmd{src: src, dst: dst}
	defer func() {
		// fmt.Println(src, dst)
		os.RemoveAll(src)
		os.RemoveAll(dst)
	}()
	os.MkdirAll(src, 0755)
	err := ioutil.WriteFile(src+"/a", []byte("aaaa"), 0755)
	assert.Nil(t, err)
	err = cmd.Start()
	assert.Nil(t, err)
	bs, err := ioutil.ReadFile(dst + "/a")
	assert.Nil(t, err)
	assert.Equal(t, "aaaa", string(bs))
}
