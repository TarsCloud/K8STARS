package main

import (
	"context"
	"time"

	"github.com/tarscloud/k8stars/consts"

	tars "github.com/tarscloud/k8stars/tarsregistry/autogen/tars"
)

func (s *registryImp) keepAlive(startReq *tars.OnStartupReq) {
	hasReg := false
	ctx := context.Background()
	keepReq := &tars.KeepAliveReq{
		NodeName:    startReq.NodeName,
		State:       consts.StateActive,
		Application: startReq.Application,
		Server:      startReq.Server,
	}
	for range time.NewTicker(time.Second * 10).C {
		if !hasReg {
			if err := s.OnStartup(ctx, startReq); err != nil {
				logger.Errorf("Register error %v", err)
				continue
			}
			hasReg = true
			go s.registryMetrics()
		}

		if err := s.KeepAlive(ctx, keepReq); err != nil {
			logger.Errorf("Keep alive error %v", err)
		}

	}
}
