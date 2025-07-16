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
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/registry"

	greet "github.com/apache/dubbo-go-samples/rpc/triple/pb2/proto"

	"github.com/dubbogo/gost/log/logger"
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_rpc_triple_pb2_client"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
	)
	if err != nil {
		panic(err)
	}

	cli, err := ins.NewClient(
		client.WithClientProtocolTriple(),
		client.WithClientSerialization(constant.ProtobufSerialization),
	)

	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	name_str := "hello world"
	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: &name_str})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet response: %s", resp)
}
