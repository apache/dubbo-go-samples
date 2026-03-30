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
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/direct/proto"
)

const (
	providerApplication = "static-tag-provider"
	triPort             = 20000
	serverName          = "server-without-tag"
)

type GreetServer struct{}

func (s *GreetServer) Greet(_ context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	logger.Infof("%s received request: %s", serverName, req.Name)
	return &greet.GreetResponse{
		Greeting: "receive: " + req.Name + ", response from: " + serverName,
	}, nil
}

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName(providerApplication),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(triPort),
		),
	)
	if err != nil {
		panic(err)
	}

	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}

	if err := greet.RegisterGreetServiceHandler(srv, &GreetServer{}); err != nil {
		panic(err)
	}

	logger.Infof("%s started on :%d", serverName, triPort)
	if err := srv.Serve(); err != nil {
		logger.Errorf("%s stopped: %v", serverName, err)
	}
}
