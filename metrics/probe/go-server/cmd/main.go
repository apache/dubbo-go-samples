/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/metrics"
	"dubbo.apache.org/dubbo-go/v3/metrics/probe"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"

	"github.com/dubbogo/gost/log/logger"

	"github.com/pkg/errors"
)

import (
	greet "github.com/apache/dubbo-go-samples/direct/proto"
)

const (
	triplePort         = 20000
	probePort          = 22222
	probeLivenessPath  = "/live"
	probeReadinessPath = "/ready"
	probeStartupPath   = "/startup"
	warmupSeconds      = 15
)

type ProbeGreetServer struct{}

func (srv *ProbeGreetServer) Greet(_ context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: "hello " + req.Name}
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	if r.Intn(101) > 99 {
		return nil, errors.New("random error")
	}
	time.Sleep(10 * time.Millisecond)
	return resp, nil
}

func main() {
	var warmupDone atomic.Bool
	var dependencyReady atomic.Bool
	dependencyReady.Store(true)

	go func() {
		logger.Infof("Warmup started, readiness/startup will fail for %ds", warmupSeconds)
		time.Sleep(warmupSeconds * time.Second)
		warmupDone.Store(true)
		logger.Info("Warmup completed, readiness/startup should be healthy")
	}()

	probe.RegisterLiveness("process", func(ctx context.Context) error {
		return nil
	})
	probe.RegisterReadiness("dependency", func(ctx context.Context) error {
		if !dependencyReady.Load() {
			return errors.New("dependency not ready")
		}
		return nil
	})
	probe.RegisterReadiness("warmup", func(ctx context.Context) error {
		if !warmupDone.Load() {
			return errors.New("warmup not complete")
		}
		return nil
	})
	probe.RegisterStartup("warmup", func(ctx context.Context) error {
		if !warmupDone.Load() {
			return errors.New("startup warmup not complete")
		}
		return nil
	})

	ins, err := dubbo.NewInstance(
		dubbo.WithMetrics(
			metrics.WithEnabled(),
			metrics.WithProbeEnabled(),
			metrics.WithProbePort(probePort),
			metrics.WithProbeLivenessPath(probeLivenessPath),
			metrics.WithProbeReadinessPath(probeReadinessPath),
			metrics.WithProbeStartupPath(probeStartupPath),
			metrics.WithProbeUseInternalState(true),
		),
	)
	if err != nil {
		panic(err)
	}

	srv, err := ins.NewServer(
		server.WithServerProtocol(
			protocol.WithTriple(),
			protocol.WithPort(triplePort),
		),
	)
	if err != nil {
		panic(err)
	}

	if err := greet.RegisterGreetServiceHandler(srv, &ProbeGreetServer{}); err != nil {
		panic(err)
	}

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Infof("Probe sample server started, triple=%d probe=%d", triplePort, probePort)
		if err := srv.Serve(); err != nil {
			logger.Error("Server error:", err)
		}
	}()

	<-stopCh
	logger.Info("Received shutdown signal, mark probe states to not ready")
	dependencyReady.Store(false)
	probe.SetReady(false)
	probe.SetStartupComplete(false)
}
