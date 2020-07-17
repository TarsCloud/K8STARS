package main

import (
	"context"
	"time"

	"github.com/tarscloud/k8stars/tarsregistry/autogen/Tars"
)

func (s *registryImp) keepAlive(startReq *Tars.OnStartupReq) {
	hasReg := false
	ctx := context.Background()
	keepReq := &Tars.KeepAliveReq{
		NodeName: startReq.NodeName,
		State:    "active",
	}
	for range time.NewTicker(time.Second * 10).C {
		if !hasReg {
			if err := s.OnStartup(ctx, startReq); err != nil {
				logger.Errorf("Register error %v", err)
				continue
			}
			hasReg = true
		}

		if err := s.KeepAlive(ctx, keepReq); err != nil {
			logger.Errorf("Keep alive error %v", err)
		}

	}
}
