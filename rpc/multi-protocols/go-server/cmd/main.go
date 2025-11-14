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
	"dubbo.apache.org/dubbo-go/v3/common"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"dubbo.apache.org/dubbo-go/v3/server"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet2 "github.com/apache/dubbo-go-samples/rpc/multi-protocols/proto"
)

type GreetMultiRPCServer struct {
}

func (srv *GreetMultiRPCServer) Greet(ctx context.Context, req *greet2.GreetRequest) (*greet2.GreetResponse, error) {
	resp := &greet2.GreetResponse{Greeting: req.Name}
	return resp, nil
}

type GreetProvider struct {
}

func (*GreetProvider) SayHello(req string, req1 string, req2 string) (string, error) {
	return req + req1 + req2, nil
}

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_multirpc_server"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20000)),
		dubbo.WithProtocol(
			protocol.WithDubbo(),
			protocol.WithPort(20001)),
		dubbo.WithProtocol(
			protocol.WithJSONRPC(),
			protocol.WithPort(20002)),
	)
	if err != nil {
		panic(err)
	}
	//Triple
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	if err = greet2.RegisterGreetServiceHandler(srv, &GreetMultiRPCServer{}); err != nil {
		panic(err)
	}
	if err = srv.Register(&GreetProvider{}, &common.ServiceInfo{}, server.WithInterface("GreetProvider")); err != nil {
		panic(err)
	}
	if err = srv.Serve(); err != nil {
		logger.Error(err)
	}
}
