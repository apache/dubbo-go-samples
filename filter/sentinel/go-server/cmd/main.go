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
	"github.com/alibaba/sentinel-golang/core/flow"
	greet "github.com/apache/dubbo-go-samples/filter/proto"
	"github.com/dubbogo/gost/log/logger"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	return resp, nil
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

	if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{},
		server.WithFilter(constant.SentinelProviderFilterKey),
	); err != nil {
		panic(err)
	}

	// Limit the request flow processed by GreetService to 100QPS
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource: greet.GreetService_ServiceInfo.InterfaceName + "::",
			// MetricType:             flow.QPS,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              100,
			RelationStrategy:       flow.CurrentResource,
		},
	})
	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
