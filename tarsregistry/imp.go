package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	rtars "github.com/TarsCloud/TarsGo/tars"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tarscloud/k8stars/consts"
	tars "github.com/tarscloud/k8stars/tarsregistry/autogen/tars"
	"github.com/tarscloud/k8stars/tarsregistry/store"
)

type registryImp struct {
	driver store.Store
}

// OnStartup is a reentrant function
func (r *registryImp) OnStartup(ctx context.Context, Req *tars.OnStartupReq) (err error) {
	logger.Debugf("OnStartup: %+v", Req)
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
func (r *registryImp) OnPrestop(ctx context.Context, Req *tars.OnPrestopReq) (err error) {
	logger.Debugf("OnPrestop: %+v", Req)
	// compatible with old version
	if Req.Application == "" && Req.Server == "" {
		return r.driver.SetServerState(ctx, Req.NodeName, Req.Application, Req.Server, consts.StateDestroyed)
	}
	return r.driver.DeleteServerConf(ctx, Req.NodeName, Req.Application, Req.Server)
}

// KeepAlive is a reentrant function
func (r *registryImp) KeepAlive(ctx context.Context, Req *tars.KeepAliveReq) (err error) {
	logger.Debugf("KeepAlive: %+v", Req)
	if err := r.driver.KeepAliveNode(ctx, Req.NodeName); err != nil {
		logger.Errorf("KeepAliveNode error %v", err)
		return err
	}
	if Req.State == "" {
		Req.State = consts.StateActive
	}
	return r.driver.SetServerState(ctx, Req.NodeName, Req.Application, Req.Server, Req.State)
}

func (r *registryImp) RegisterMetrics(ctx context.Context, Req *tars.RegisterMetricsReq) (err error) {
	logger.Debugf("RegisterMetrics: %+v", Req)
	if Req.MetricsPort == 0 {
		Req.MetricsPort = int32(consts.MetricsPort)
	}
	return r.driver.RegisterMetrics(ctx, Req.NodeName, Req.Application, Req.Server, int(Req.MetricsPort))
}

func (r *registryImp) GetMetricsAdapters(ctx context.Context, Req *tars.GetMetricsAdaptersReq, Rsp *[]tars.MetricsAdapterInfo) (err error) {
	targets, err := r.driver.GetMetricTargets(ctx)
	if err != nil {
		return err
	}

	ret := make([]tars.MetricsAdapterInfo, 0)
	retMap := make(map[string]int)
	for _, t := range targets {
		key := fmt.Sprintf("%s|%s|%s", t.Application, t.Server, t.SetID)
		if idx, ok := retMap[key]; ok {
			ret[idx].Targets = append(ret[idx].Targets, t.Address)
		} else {
			info := tars.MetricsAdapterInfo{
				Targets: []string{t.Address},
				Labels: map[string]string{
					"application": t.Application,
					"server":      t.Server,
					"set":         t.SetID,
				},
			}
			retMap[key] = len(ret)
			ret = append(ret, info)
		}
	}
	*Rsp = ret
	return nil
}

func (r *registryImp) registryMetrics() {
	cfg := rtars.GetServerConfig()
	rReq := tars.RegisterMetricsReq{
		NodeName:    cfg.LocalIP,
		Application: cfg.App,
		Server:      cfg.Server,
	}

	for i := 0; i < 100; i++ {
		if err := r.RegisterMetrics(context.Background(), &rReq); err != nil {
			logger.Errorf("RegisterMetrics error %v", err)
			time.Sleep(time.Minute)
			continue
		}
		break
	}

	http.Handle("/metrics", promhttp.Handler())
	addr := cfg.LocalIP + ":" + fmt.Sprint(consts.MetricsPort)
	logger.Debugf("Listen metrics %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Errorf("ListenAndServe error %v", err)
	}
}

func (r *registryImp) cleanAll(cmd string) (string, error) {
	args := strings.Fields(cmd)
	var datetime string
	var dryRun bool
	if len(args) == 1 {
		dryRun = true
		datetime = time.Now().Add(-time.Hour * 24).Format("2006-01-02 15:04:05")
	} else if len(args) == 2 {
		dryRun = args[1] != "n" && args[1] != "N"
		datetime = time.Now().Add(-time.Hour * 24).Format("2006-01-02 15:04:05")
	} else if len(args) == 3 {
		dryRun = args[1] != "n" && args[1] != "N"
		datetime = args[2]
	} else {
		return "", fmt.Errorf("usage: %s [dryRun] [datetime]", args[0])
	}
	ret, err := r.driver.DeleteAllInactive(context.Background(), datetime, dryRun)
	if err != nil {
		return "", err
	}
	res := "\n" + strings.Join(ret, "\n")
	if dryRun {
		res += "\nDry run.\n"
	} else {
		res += "\nDeleted\n"
	}
	return res, nil
}
