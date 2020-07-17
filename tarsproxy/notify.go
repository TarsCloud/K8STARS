package tarsproxy

import (
	"context"
	"time"

	"github.com/tarscloud/k8stars/algorithm/retry"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/notifyf"
)

var mockNotifyClient NotifyClient
var impNoitfyClient NotifyClient

// NotifyClient is client of tars registry
type NotifyClient interface {
	ReportNotify(ctx context.Context, Req *notifyf.ReportInfo) (err error)
}

// GetNotifyClient returns client of tars registry
func GetNotifyClient(locator string) NotifyClient {
	if mockNotifyClient != nil {
		return mockNotifyClient
	}
	if impNoitfyClient != nil {
		return impNoitfyClient
	}
	client := &notifyf.Notify{}
	StringToProxy(locator, "tars.tarsnotify.NotifyObj", client)
	client.TarsSetTimeout(rpcTimeout)
	impNoitfyClient = &notifyClientImp{
		client: client,
		retry:  retry.New(retry.MaxTimeoutOpt(time.Second*100, time.Second*3)),
	}
	return impNoitfyClient
}

type notifyClientImp struct {
	client *notifyf.Notify
	retry  retry.Func
}

func (r *notifyClientImp) ReportNotify(ctx context.Context, Req *notifyf.ReportInfo) (err error) {
	return r.retry(func() error {
		return r.client.ReportNotifyInfoWithContext(ctx, Req)
	})
}

type notifyClientMock struct {
}

func (r *notifyClientMock) ReportNotify(ctx context.Context, Req *notifyf.ReportInfo) (err error) {
	return nil
}
