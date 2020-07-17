package tinycli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeDuration(t *testing.T) {
	inputs := []string{"0", "1ns", "1us", "1ms", "-1s", "1m", "2h", "1d"}
	output := []time.Duration{0, time.Duration(1), time.Microsecond, time.Millisecond, -time.Second, time.Minute, time.Hour * 2, time.Hour * 24}

	for i := range inputs {
		ret, err := GetTimeDuration(inputs[i])
		assert.Nil(t, err)
		assert.Equal(t, output[i], ret)
	}

	_, err := GetTimeDuration("a1")
	assert.NotNil(t, err)
	_, err = GetTimeDuration("a1s")
	assert.NotNil(t, err)
	_, err = GetTimeDuration("a1ms")
	assert.NotNil(t, err)

	var ii int
	err = setFlagInt(true, "-1a1", &ii, 0, "")
	assert.NotNil(t, err)
	var dd time.Duration
	err = setFlagDuration(true, "100", &dd, "", "")
	assert.NotNil(t, err)
	err = setFlagDuration(true, "100as", &dd, "", "")
	assert.NotNil(t, err)
	err = setFlagDuration(true, "100ms", &dd, "10mm", "")
	assert.NotNil(t, err)
}
