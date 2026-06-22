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
	"dubbo.apache.org/dubbo-go/v3/global"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
)

import (
	greet "github.com/apache/dubbo-go-samples/rpc/triple/old_triple/proto"
)

type GreetDubbo3Server struct {
	greet.UnimplementedGreetServiceServer
}

func (srv *GreetDubbo3Server) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	return &greet.GreetResponse{Greeting: req.Name}, nil
}

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_rpc_triple_old_triple_server"),
	)
	if err != nil {
		panic(err)
	}

	provider := global.DefaultProviderConfig()
	// make sure using old triple
	provider.Services["GreetDubbo3Server"] = &global.ServiceConfig{
		Interface: "greet.GreetService",
	}

	srv, err := ins.NewServer(
		server.SetServerProvider(provider),
		server.WithServerProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20000),
		),
	)
	if err != nil {
		panic(err)
	}

	if err := srv.Register(&GreetDubbo3Server{}, nil, server.WithInterface("greet.GreetService")); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		panic(err)
	}
}
