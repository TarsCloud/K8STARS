package recentuse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRecentUse(t *testing.T) {
	ru := NewRecentUse(time.Millisecond * 5)
	ret := ru.KeepAlive("a")
	assert.True(t, !ret)
	ret = ru.KeepAlive("a")
	assert.True(t, ret)

	time.Sleep(time.Millisecond * 3)
	ret = ru.KeepAlive("a")
	assert.True(t, ret)

	time.Sleep(time.Millisecond * 3)
	ret = ru.KeepAlive("a")
	assert.True(t, !ret)

	ret = ru.KeepAlive("a")
	assert.True(t, ret)

	ret = ru.KeepAlive("b")
	assert.True(t, !ret)

	ret = ru.KeepAlive("c")
	ret = ru.KeepAlive("b")

	time.Sleep(time.Millisecond * 6)
	ret = ru.KeepAlive("b")
	assert.True(t, !ret)
	assert.Equal(t, 1, len(ru.m))
	assert.Equal(t, ru.tail, ru.head)
	assert.Nil(t, ru.head.next)
	assert.Nil(t, ru.tail.prev)
}

func TestListNode(t *testing.T) {
	n := &node{
		value: "1",
	}
	assert.Equal(t, "1", n.String())
	n.next = &node{
		value: "1",
	}
	assert.Equal(t, "1->1", n.String())
}
