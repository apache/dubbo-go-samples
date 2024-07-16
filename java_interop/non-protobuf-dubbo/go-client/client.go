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
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
)

import (
	greet_gen "github.com/apache/dubbo-go-samples/java_interop/non-protobuf-dubbo/proto"
)

func main() {
	ctx := context.Background()

	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_rpc_hessian2_client"),
	)
	if err != nil {
		panic(err)
	}

	cli, err := ins.NewClient(
		client.WithClientURL("127.0.0.1:50052"),
		client.WithClientProtocolDubbo(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet_gen.NewGreetingsService(cli)
	if err != nil {
		panic(err)
	}

	resp, err := svc.Greet(ctx, &greet_gen.GreetRequest{
		Name: "dubbo-go",
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Greeting)
}
