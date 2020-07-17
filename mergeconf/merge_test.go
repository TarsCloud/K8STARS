package mergeconf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	src1 := []byte(`<a><b>c=1
d
</b></a>`)
	src2 := []byte(`<a>c=1<b>c=2</b>
	</a>
	`)

	out, err := MergeConf(src1, src2)
	assert.Nil(t, err)
	root, err := initFromBytes(out)
	assert.Nil(t, err)
	assert.Equal(t, "2", root.children["a"].children["b"].children["c"].value)
	assert.Equal(t, "", root.children["a"].children["b"].children["d"].value)
	assert.Equal(t, "1", root.children["a"].children["c"].value)
}
