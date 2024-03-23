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
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"dubbo.apache.org/dubbo-go/v3/server"
	greet "github.com/apache/dubbo-go-samples/multirpc/proto"
	"github.com/dubbogo/gost/log/logger"
)

type GreetMultiRPCServer struct {
}

func (srv *GreetMultiRPCServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	return resp, nil
}

type GreetProvider struct {
}

func (*GreetProvider) SayHello(req string, req1 string, req2 string) (string, error) {
	return req + req1 + req2, nil
}

func main() {
	//triple
	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_multirpc_triple_server"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
		//dubbo.WithProtocol(protocol.WithDubbo()),
		dubbo.WithProtocol(protocol.WithTriple()),
		//dubbo.WithProtocol(protocol.WithJSONRPC()),
	)
	if err != nil {
		panic(err)
	}
	//var protocols []string
	//protocols = append(protocols, "tri", "dubbo", "jsonrpc")
	srv, err := ins.NewServer(
		server.WithServerProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20000),
		),
		//server.WithServerProtocolIDs(protocols),
	)
	if err != nil {
		panic(err)
	}
	// 利用生成代码注册业务逻辑(GreetTripleServer)
	// service配置，可以在此处覆盖server注入的默认配置
	// 若观察RegisterGreetServiceHandler的代码，会发现本质上是调用Server.Register
	if err = greet.RegisterGreetServiceHandler(srv, &GreetMultiRPCServer{}); err != nil {
		panic(err)
	}
	// 运行
	if err = srv.Serve(); err != nil {
		logger.Error(err)
	}

	//dubbo
	insDubbo, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_multirpc_dubbo_server"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
		dubbo.WithProtocol(protocol.WithDubbo()),
		//dubbo.WithProtocol(protocol.WithTriple()),
		//dubbo.WithProtocol(protocol.WithJSONRPC()),
	)
	if err != nil {
		panic(err)
	}

	srvDubbo, err := insDubbo.NewServer(
		server.WithServerProtocol(
			protocol.WithDubbo(),
			protocol.WithPort(20001),
		),
		//server.WithServerProtocolIDs(protocols),
	)
	if err != nil {
		panic(err)
	}
	if err = srvDubbo.Register(&GreetProvider{}, nil, server.WithInterface("GreetProvider")); err != nil {
		panic(err)
	}

	// 运行
	if err = srvDubbo.Serve(); err != nil {
		logger.Error(err)
	}

	//JsonRpc
	insJsonRpc, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_multirpc_jsonrpc_server"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
		dubbo.WithProtocol(protocol.WithJSONRPC()),
		//dubbo.WithProtocol(protocol.WithTriple()),
		//dubbo.WithProtocol(protocol.WithJSONRPC()),
	)
	if err != nil {
		panic(err)
	}

	srvJsonRpc, err := insJsonRpc.NewServer(
		server.WithServerProtocol(
			protocol.WithJSONRPC(),
			protocol.WithPort(20002),
		),
		//server.WithServerProtocolIDs(protocols),
	)
	if err != nil {
		panic(err)
	}
	if err = srvJsonRpc.Register(&GreetProvider{}, nil, server.WithInterface("GreetProvider")); err != nil {
		panic(err)
	}

	// 运行
	if err := srvJsonRpc.Serve(); err != nil {
		logger.Error(err)
	}
}
