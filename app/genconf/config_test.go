package genconf

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/tarscloud/k8stars/consts"
	"github.com/stretchr/testify/assert"
)

func TestGenerateConf(t *testing.T) {
	consts.TarsPath = "testtarsdir"
	infile := consts.TarsPath + "/bin/_server_meta.yaml"
	conf := defaultServerConf()
	err := parseServerConf(infile, &conf)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(conf.Adapters))
	assert.Equal(t, 10, conf.Adapters[0].Threads)
	assert.Equal(t, 5, conf.Adapters[1].Threads) // default

	err = generateConf(&conf, conf.Server, "")
	assert.Nil(t, err)

	outfile := filepath.Join(consts.TarsPath, "data", consts.ServerInfoFile)
	_, err = ioutil.ReadFile(outfile)
	assert.Nil(t, err)
}
