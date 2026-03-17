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
	"flag"
	"fmt"
	"strings"
	"time"

	"dubbo.apache.org/dubbo-go/v3/graceful_shutdown"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/grpc"
	"dubbo.apache.org/dubbo-go/v3/server"

	greet "github.com/apache/dubbo-go-samples/graceful_shutdown/proto"

	"github.com/dubbogo/gost/log/logger"
)

type GreetProvider struct {
	fixedDelay time.Duration
}

func (p *GreetProvider) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	start := time.Now()
	logger.Infof("Handling greet request, name=%s delay=%s", req.Name, p.fixedDelay)

	if p.fixedDelay > 0 {
		timer := time.NewTimer(p.fixedDelay)
		defer timer.Stop()

		select {
		case <-timer.C:
		case <-ctx.Done():
			logger.Warnf("Greet request canceled before completion, name=%s err=%v", req.Name, ctx.Err())
			return nil, ctx.Err()
		}
	}

	resp := &greet.GreetResponse{
		Greeting: fmt.Sprintf("%s response after %s", req.Name, time.Since(start).Truncate(time.Millisecond)),
	}
	logger.Infof("Greet request finished, name=%s cost=%s", req.Name, time.Since(start).Truncate(time.Millisecond))
	return resp, nil
}

func main() {
	protocols := flag.String("protocols", "tri", "protocols to expose: tri,dubbo,grpc (comma separated)")
	portBase := flag.Int("port-base", 20000, "base port number")
	notifyTimeout := flag.Duration("notify-timeout", 5*time.Second, "overall timeout budget for active long-connection notices")
	stepTimeout := flag.Duration("step-timeout", 3*time.Second, "timeout for waiting provider and consumer in-flight requests")
	consumerUpdateWait := flag.Duration("consumer-update-wait", 3*time.Second, "time to wait for consumers to observe instance changes")
	offlineWindow := flag.Duration("offline-window", 3*time.Second, "time window for observing late requests after offline")
	requestDelay := flag.Duration("delay", 0, "artificial delay added to each greet request")
	flag.Parse()

	graceful_shutdown.Init(
		graceful_shutdown.WithTimeout(60*time.Second),
		graceful_shutdown.WithStepTimeout(*stepTimeout),
		graceful_shutdown.WithNotifyTimeout(*notifyTimeout),
		graceful_shutdown.WithConsumerUpdateWaitTime(*consumerUpdateWait),
		graceful_shutdown.WithOfflineRequestWindowTimeout(*offlineWindow),
	)
	logger.Infof("Graceful shutdown initialized, notify-timeout=%s step-timeout=%s consumer-update-wait=%s offline-window=%s request-delay=%s",
		notifyTimeout.String(), stepTimeout.String(), consumerUpdateWait.String(), offlineWindow.String(), requestDelay.String())

	var serverOpts []server.ServerOption
	port := *portBase

	protoList := map[string]func(int){
		"tri": func(p int) {
			serverOpts = append(serverOpts,
				server.WithServerProtocol(
					protocol.WithProtocol("tri"),
					protocol.WithPort(p),
					protocol.WithID("tri"),
				),
			)
		},
		"dubbo": func(p int) {
			serverOpts = append(serverOpts,
				server.WithServerProtocol(
					protocol.WithProtocol("dubbo"),
					protocol.WithPort(p),
					protocol.WithID("dubbo"),
				),
			)
		},
		"grpc": func(p int) {
			serverOpts = append(serverOpts,
				server.WithServerProtocol(
					protocol.WithProtocol("grpc"),
					protocol.WithPort(p),
					protocol.WithID("grpc"),
				),
			)
		},
	}

	for _, p := range strings.Split(*protocols, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if fn, ok := protoList[p]; ok {
			fn(port)
			logger.Infof("Exposing protocol %s on port %d", p, port)
			port += 10000
		} else {
			logger.Warnf("Unknown protocol: %s", p)
		}
	}

	if len(serverOpts) == 0 {
		logger.Error("No valid protocols specified")
		return
	}

	srv, err := server.NewServer(serverOpts...)
	if err != nil {
		logger.Fatalf("failed to create server: %v", err)
	}

	provider := &GreetProvider{fixedDelay: *requestDelay}
	if err := greet.RegisterGreetServiceHandler(srv, provider); err != nil {
		logger.Fatalf("failed to register greet service handler: %v", err)
	}

	logger.Info("Server started, press Ctrl+C to trigger graceful shutdown")

	if err := srv.Serve(); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
