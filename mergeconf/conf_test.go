package mergeconf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConf(t *testing.T) {
	inBytes := []byte(`<a>
	<b>
	    d
		c=111
	</b>
</a>`)

	root, err := initFromBytes(inBytes)
	assert.Nil(t, err)
	assert.Equal(t, "111", root.children["a"].children["b"].children["c"].value)
}
