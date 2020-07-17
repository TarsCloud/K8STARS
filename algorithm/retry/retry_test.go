package retry

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMaxRetiesOpt(t *testing.T) {
	retry := New(MaxRetiesOpt(2))

	var count int
	err := retry(func() error {
		count++
		if count <= 1 {
			return fmt.Errorf("ee")
		}
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, count)

	count = 0
	err = retry(func() error {
		count++
		if count <= 5 {
			return fmt.Errorf("ee")
		}
		return nil
	})
	assert.NotNil(t, err)
	assert.Equal(t, 3, count)
}

func TestMaxTimeoutOpt(t *testing.T) {
	retry := New(MaxTimeoutOpt(time.Millisecond*30, time.Millisecond*10))

	var count int
	err := retry(func() error {
		count++
		if count <= 1 {
			time.Sleep(time.Microsecond * 20)
			return fmt.Errorf("ee")
		}
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, count)

	count = 0
	err = retry(func() error {
		count++
		if count <= 1 {
			return fmt.Errorf("ee")
		} else if count <= 2 {
			time.Sleep(time.Millisecond * 10)
			return fmt.Errorf("ee")
		} else if count <= 3 {
			time.Sleep(time.Millisecond * 20)
			return fmt.Errorf("ee")
		}
		return nil
	})
	assert.NotNil(t, err)
	assert.Equal(t, 3, count)
}
