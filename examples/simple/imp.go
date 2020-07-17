package main

import (
	"context"
	"time"
)

type simpleServerImp struct{}

func (s *simpleServerImp) Sum(ctx context.Context, A int32, B int32) (ret int32, err error) {
	ret = A + B
	time.Sleep(time.Millisecond * 100)
	return
}
