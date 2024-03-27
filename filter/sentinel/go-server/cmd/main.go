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
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/isolation"
	"github.com/alibaba/sentinel-golang/util"
	greet "github.com/apache/dubbo-go-samples/filter/proto/sentinel"
	"github.com/dubbogo/gost/log/logger"
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

type stateChangeTestListener struct {
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	logger.Infof("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	logger.Infof("rule.steategy: %+v, From %s to Open, snapshot: %.2f, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	logger.Infof("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, request *greet.GreetRequest) (*greet.GreetResponse, error) {
	time.Sleep(time.Duration(rand.Uint64()%80+20) * time.Millisecond)
	resp := &greet.GreetResponse{Greeting: request.Name}
	return resp, nil
}

func (srv *GreetTripleServer) GreetWithChanceOfError(ctx context.Context, request *greet.GreetRequest) (*greet.GreetResponse, error) {
	if request.Name == "error" {
		return nil, errors.New("greet, error")
	}
	return srv.Greet(ctx, request)
}

func (srv *GreetTripleServer) GreetWithQPSLimit(ctx context.Context, request *greet.GreetRequest) (*greet.GreetResponse, error) {
	return srv.Greet(ctx, request)
}

func (srv *GreetTripleServer) GreetWithConcurrencyLimit(ctx context.Context, request *greet.GreetRequest) (*greet.GreetResponse, error) {
	return srv.Greet(ctx, request)
}

func main() {
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000),
			protocol.WithTriple(),
		),
	)
	if err != nil {
		panic(err)
	}

	if err = greet.RegisterSentinelGreetServiceHandler(srv, &GreetTripleServer{},
		server.WithFilter(constant.SentinelProviderFilterKey),
	); err != nil {
		panic(err)
	}

	// Limit the request flow processed by GreetService to 100QPS
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource: "dubbo:provider:greet.SentinelGreetService:::GreetWithQPSLimit()",
			// MetricType:             flow.QPS,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              100,
			RelationStrategy:       flow.CurrentResource,
		},
	})
	if err != nil {
		panic(err)
	}

	_, err = isolation.LoadRules([]*isolation.Rule{
		{
			Resource:   "dubbo:provider:greet.SentinelGreetService:::GreetWithConcurrencyLimit()",
			MetricType: isolation.Concurrency,
			Threshold:  100,
		},
	})
	if err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
