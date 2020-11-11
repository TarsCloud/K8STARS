package store

import "context"

type ServerConf struct {
	ID          string
	Application string
	Server      string
	NodeName    string
	EnableSet   string
	SetName     string
	SetArea     string
	SetGroup    string
	GridFlag    string
	State       string
	Version     string
}

type AdapterConf struct {
	Application  string
	Server       string
	NodeName     string
	AdapterName  string
	Servant      string
	Protocol     string
	ThreadNum    int
	Endpoint     string
	MaxConns     int
	QueueCap     int
	QueueTimeout int
}

type MetricsTarget struct {
	SetID       string
	Application string
	Server      string
	Address     string
}

type Store interface {
	RegisterNode(ctx context.Context, nodeName string) error
	RegisterServer(ctx context.Context, conf *ServerConf) error
	RegistryAdapter(ctx context.Context, conf []*AdapterConf) error
	DeleteServerConf(ctx context.Context, nodeName, application, server string) error
	KeepAliveNode(ctx context.Context, nodeName string) error
	DeleteAllInactive(ctx context.Context, datetime string, dryRun bool) ([]string, error)
	SetServerState(ctx context.Context, nodeName, application, server, state string) error
	RegisterMetrics(ctx context.Context, nodeName, application, server string, port int) error
	GetMetricTargets(ctx context.Context) ([]MetricsTarget, error)
}
