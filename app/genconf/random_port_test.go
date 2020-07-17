package genconf

import (
	"fmt"
	"testing"

	"github.com/tarscloud/k8stars/consts"
	"github.com/stretchr/testify/assert"
)

func TestRandPort(t *testing.T) {
	ep := fmt.Sprintf("tcp -h %s -p ${rrr}", consts.LocalIP)
	p, err := getRandomPort(ep, 10000, 10002)
	assert.Nil(t, err)
	assert.Equal(t, "10000", p)
}
