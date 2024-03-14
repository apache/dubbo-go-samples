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
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/registry"
	greet "github.com/apache/dubbo-go-samples/multirpc/proto"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_multirpc_triple_client"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
	)
	if err != nil {
		panic(err)
	}
	//test triple
	cliTriple, err := ins.NewClient(
		client.WithClientProtocolTriple(),
	)
	if err != nil {
		panic(err)
	}
	svc, err := greet.NewGreetService(cliTriple)
	if err != nil {
		panic(err)
	}

	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet multirpc response: %s", resp)

	//test dubbo
	ins_dubbo, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_multirpc_dubbo_client"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
	)
	if err != nil {
		panic(err)
	}
	cli_dubbo, err := ins_dubbo.NewClient(
		client.WithClientProtocolDubbo(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		panic(err)
	}

	connDubbo, err := cli_dubbo.Dial("GreetProvider")
	if err != nil {
		panic(err)
	}
	var respDubbo string
	if err := connDubbo.CallUnary(context.Background(), []interface{}{"hello", "new", "dubbo"}, &respDubbo, "Greet"); err != nil {
		logger.Errorf("GreetProvider.Greet err: %s", err)
		return
	}
	logger.Infof("Get Response: %s", respDubbo)

	//test json rpc
	insJsonRpc, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_multirpc_jsonrpc_client"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
	)
	if err != nil {
		panic(err)
	}
	cliJsonRpc, err := insJsonRpc.NewClient(
		client.WithClientProtocolJsonRPC(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		panic(err)
	}

	connJsonRpc, err := cliJsonRpc.Dial("GreetProvider")
	if err != nil {
		panic(err)
	}
	var respJsonRpc string
	if err := connJsonRpc.CallUnary(context.Background(), []interface{}{"hello", "new", "dubbo"}, &respJsonRpc, "Greet"); err != nil {
		logger.Errorf("GreetProvider.Greet err: %s", err)
		return
	}
	logger.Infof("Get Response: %s", respJsonRpc)

}
