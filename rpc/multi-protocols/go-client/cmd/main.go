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
	"fmt"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet2 "github.com/apache/dubbo-go-samples/rpc/multi-protocols/proto"
)

func main() {
	// 初始化实例
	ins, err := initDubboInstance()
	if err != nil {
		panic(err)
	}

	// 调用不同协议的服务
	if err := callTripleService(ins); err != nil {
		logger.Error(err)
	}

	if err := callDubboService(ins); err != nil {
		logger.Error(err)
	}

	if err := callJsonRpcService(ins); err != nil {
		logger.Error(err)
	}
}

// 初始化 Dubbo 实例
func initDubboInstance() (*dubbo.Instance, error) {
	return dubbo.NewInstance(
		dubbo.WithName("dubbo_multirpc_client"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
	)
}

// 调用 Triple 协议服务
func callTripleService(ins *dubbo.Instance) error {
	cli, err := ins.NewClient(client.WithClientProtocolTriple())
	if err != nil {
		return fmt.Errorf("failed to create triple client: %w", err)
	}

	svc, err := greet2.NewGreetService(cli)
	if err != nil {
		return fmt.Errorf("failed to create greet service: %w", err)
	}

	resp, err := svc.Greet(context.Background(), &greet2.GreetRequest{Name: "hello world"})
	if err != nil {
		return fmt.Errorf("triple service call failed: %w", err)
	}

	logger.Infof("Greet triple response: %s", resp.Greeting)
	return nil
}

// 调用 Dubbo 协议服务
func callDubboService(ins *dubbo.Instance) error {
	cli, err := ins.NewClient(
		client.WithClientProtocolDubbo(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		return fmt.Errorf("failed to create dubbo client: %w", err)
	}

	conn, err := cli.Dial("GreetProvider")
	if err != nil {
		return fmt.Errorf("failed to dial dubbo service: %w", err)
	}

	var resp string
	if err := conn.CallUnary(context.Background(), []interface{}{"hello", "new", "dubbo"}, &resp, "SayHello"); err != nil {
		return fmt.Errorf("dubbo service call failed: %w", err)
	}

	logger.Infof("Get dubbo Response: %s", resp)
	return nil
}

// 调用 JsonRPC 协议服务
func callJsonRpcService(ins *dubbo.Instance) error {
	cli, err := ins.NewClient(
		client.WithClientProtocolJsonRPC(),
		client.WithClientSerialization(constant.JSONSerialization),
	)
	if err != nil {
		return fmt.Errorf("failed to create jsonrpc client: %w", err)
	}

	conn, err := cli.Dial("GreetProvider")
	if err != nil {
		return fmt.Errorf("failed to dial jsonrpc service: %w", err)
	}

	var resp string
	if err := conn.CallUnary(context.Background(), []interface{}{"hello", "new", "jsonrpc"}, &resp, "SayHello"); err != nil {
		return fmt.Errorf("jsonrpc service call failed: %w", err)
	}

	logger.Infof("Get jsonrpc Response: %s", resp)
	return nil
}
