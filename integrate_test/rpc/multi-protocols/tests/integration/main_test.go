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

package integration

import (
	"os"
	"testing"
)

import (
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	greet "github.com/apache/dubbo-go-samples/rpc/multi-protocols/proto"
)

var greeterProvider greet.GreetService
var connDubbo *client.Connection
var connJsonRpc *client.Connection

func TestMain(m *testing.M) {
	//Triple
	cliTriple, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
		client.WithClientClusterFailFast(),
		client.WithClientRetries(0),
	)
	if err != nil {
		panic(err)
	}

	greeterProvider, err = greet.NewGreetService(cliTriple)

	if err != nil {
		panic(err)
	}

	//Dubbo
	cliDubbo, err := client.NewClient(
		client.WithClientURL("dubbo://127.0.0.1:20001"),
		client.WithClientProtocolDubbo(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		panic(err)
	}

	connDubbo, err = cliDubbo.Dial("GreetProvider")
	if err != nil {
		panic(err)
	}

	//JsonRpc
	cliJsonRpc, err := client.NewClient(
		client.WithClientURL("jsonrpc://127.0.0.1:20002"),
		client.WithClientProtocolJsonRPC(),
		client.WithClientSerialization(constant.JSONSerialization),
	)
	if err != nil {
		panic(err)
	}

	connJsonRpc, err = cliJsonRpc.Dial("GreetProvider")

	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
