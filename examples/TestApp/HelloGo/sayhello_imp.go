package main

import (
	"context"
)

// SayHelloImp servant implementation
type SayHelloImp struct {
}

// Init servant init
func (imp *SayHelloImp) Init() (error) {
	//initialize servant here:
	//...
	return nil
}

// Destroy servant destory
func (imp *SayHelloImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *SayHelloImp) Add(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	*c = a + b
	return 0, nil
}
func (imp *SayHelloImp) Sub(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
        *c = a - b
	return 0, nil
}
