package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/tarscloud/k8stars/tarsregistry/store"

	"github.com/tarscloud/k8stars/tarsregistry/autogen/Tars"
)

type registryImp struct {
	driver store.Store
}

// OnStartup is a reentrant function
func (r *registryImp) OnStartup(ctx context.Context, Req *Tars.OnStartupReq) (err error) {
	// register node
	if err := r.driver.RegisterNode(ctx, Req.NodeName); err != nil {
		return fmt.Errorf("RegisterNode error %v", err)
	}

	// register server
	sConf := &store.ServerConf{
		Application: Req.Application,
		Server:      Req.Server,
		NodeName:    Req.NodeName,
		State:       Req.State,
		Version:     Req.Version,
		EnableSet:   "N",
		GridFlag:    "NORMAL",
	}
	if Req.DisableFlow {
		sConf.GridFlag = "NO_FLOW"
	}
	tmps := strings.Split(Req.SetID, ".")
	if len(tmps) == 3 {
		sConf.SetName = tmps[0]
		sConf.SetGroup = tmps[1]
		sConf.SetArea = tmps[2]
		sConf.EnableSet = "Y"
	}
	if err := r.driver.RegisterServer(ctx, sConf); err != nil {
		return fmt.Errorf("RegisterServer error %v", err)
	}

	// register adapter
	adapters := make([]*store.AdapterConf, len(Req.Adapters))
	for i := range Req.Adapters {
		ad := Req.Adapters[i]
		adapters[i] = &store.AdapterConf{
			Application:  Req.Application,
			Server:       Req.Server,
			NodeName:     Req.NodeName,
			Servant:      ad.Servant,
			AdapterName:  ad.Servant + "Adapter",
			Endpoint:     ad.Endpoint,
			Protocol:     ad.Protocol,
			QueueCap:     int(ad.QueueCap),
			QueueTimeout: int(ad.QueueTimeout),
			MaxConns:     int(ad.MaxConns),
			ThreadNum:    int(ad.ThreadNum),
		}
	}
	if err := r.driver.RegistryAdapter(ctx, adapters); err != nil {
		return fmt.Errorf("RegistryAdapter error %v", err)
	}
	return nil
}

// OnPrestop is a reentrant function
func (r *registryImp) OnPrestop(ctx context.Context, Req *Tars.OnPrestopReq) (err error) {
	logger.Debugf("NodeName:%s", Req.NodeName)
	return r.driver.DeleteNodeConf(ctx, Req.NodeName)
}

// KeepAlive is a reentrant function
func (r *registryImp) KeepAlive(ctx context.Context, Req *Tars.KeepAliveReq) (err error) {
	logger.Debugf("NodeName:%s, State:%s, Application:%s, Server:%s, SetID:%s", 
		Req.NodeName, Req.State, Req.Application, Req.Server, Req.SetID)
	if err := r.driver.KeepAliveNode(ctx, Req.NodeName); err != nil {
		logger.Errorf("KeepAliveNode error %v", err)
		return err
	}
	if Req.State == "" {
		Req.State = "active"
	}
	return r.driver.SetServerState(ctx, Req.NodeName, Req.Application, Req.Server, Req.State)
}
