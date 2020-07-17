package store

import "context"

type ServerConf struct {
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

type Store interface {
	RegisterNode(ctx context.Context, nodeName string) error
	RegisterServer(ctx context.Context, conf *ServerConf) error
	RegistryAdapter(ctx context.Context, conf []*AdapterConf) error
	DeleteNodeConf(ctx context.Context, nodeName string) error
	KeepAliveNode(ctx context.Context, nodeName string) error
	SetServerState(ctx context.Context, nodeName, state string) error
}
