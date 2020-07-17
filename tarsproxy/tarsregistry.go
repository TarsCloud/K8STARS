package tarsproxy

import (
	"context"
	"strings"
	"time"

	"github.com/tarscloud/k8stars/algorithm/retry"
	"github.com/tarscloud/k8stars/tarsregistry/autogen/Tars"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/endpointf"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/queryf"
)

// RegistryClient is client of tars registry
type RegistryClient interface {
	OnStartup(ctx context.Context, Req *Tars.OnStartupReq) (err error)
	OnPrestop(ctx context.Context, Req *Tars.OnPrestopReq) (err error)
	KeepAlive(ctx context.Context, Req *Tars.KeepAliveReq) (err error)
}

// QueryClient is client of tars registry for query
type QueryClient interface {
	FindObjectByIdInSameGroup(Id string, ActiveEp *[]endpointf.EndpointF, InactiveEp *[]endpointf.EndpointF, _opt ...map[string]string) (ret int32, err error)
}

// GetRegistryClient returns client of tars registry
func GetRegistryClient(locator string) RegistryClient {
	if mockClient != nil {
		return mockClient
	}
	if impClient != nil {
		return impClient
	}
	client := &Tars.Tarsregistry{}
	if err := StringToProxy(locator, "tars.tarsregistry.Registry", client); err != nil {
		return nil
	}
	client.TarsSetTimeout(rpcTimeout)
	impClient = &registryClientImp{
		client: client,
		retry:  retry.New(retry.MaxTimeoutOpt(time.Second*100, time.Second*3)),
	}
	return impClient
}

// GetQueryClient returns client of tars registry for query
func GetQueryClient(locator string) QueryClient {
	if mockQuery != nil {
		return mockQuery
	}
	if impQuery != nil {
		return impQuery
	}
	if !strings.Contains(locator, "@") {
		return nil
	}
	q := &queryf.QueryF{}
	comm.StringToProxy(locator, q)
	q.TarsSetTimeout(rpcTimeout)
	impQuery = q
	return q
}

type registryClientImp struct {
	client *Tars.Tarsregistry
	retry  retry.Func
}

var mockClient RegistryClient
var impClient RegistryClient

var mockQuery QueryClient
var impQuery QueryClient

func (r *registryClientImp) OnStartup(ctx context.Context, Req *Tars.OnStartupReq) (err error) {
	return r.retry(func() error {
		return r.client.OnStartupWithContext(ctx, Req)
	})
}
func (r *registryClientImp) OnPrestop(ctx context.Context, Req *Tars.OnPrestopReq) (err error) {
	return r.retry(func() error {
		return r.client.OnPrestopWithContext(ctx, Req)
	})
}
func (r *registryClientImp) KeepAlive(ctx context.Context, Req *Tars.KeepAliveReq) (err error) {
	return r.retry(func() error {
		return r.client.KeepAliveWithContext(ctx, Req)
	})
}

type registryClientMock struct {
}

func (r *registryClientMock) OnStartup(ctx context.Context, Req *Tars.OnStartupReq) (err error) {
	return nil
}
func (r *registryClientMock) OnPrestop(ctx context.Context, Req *Tars.OnPrestopReq) (err error) {
	return nil
}
func (r *registryClientMock) KeepAlive(ctx context.Context, Req *Tars.KeepAliveReq) (err error) {
	return nil
}
