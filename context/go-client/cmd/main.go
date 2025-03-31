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
	"dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
	"net/http"

	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/context/proto"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	header := http.Header{"testKey1": []string{"testVal1"}, "testKey2": []string{"testVal2"}}
	// to store outgoing data ,and reserve the location for the receive field.
	// header will be copy , and header's key will change to be lowwer.
	ctx := triple_protocol.NewOutgoingContext(context.Background(), header)
	ctx = triple_protocol.AppendToOutgoingContext(ctx, "testKey3", "testVal3")

	resp, err := svc.Greet(ctx, &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error(err)
	}
	extractedHeader, _ := triple_protocol.FromIncomingContext(ctx)

	var value1, value2 string
	if values, ok := extractedHeader["outgoingcontextkey1"]; ok && len(values) > 0 {
		value1 = values[0]
		logger.Infof("OutgoingContextKey1: %s", value1)
	}
	if values, ok := extractedHeader["outgoingcontextkey2"]; ok && len(values) > 0 {
		value2 = values[0]
		logger.Infof("OutgoingContextKey2: %s", value2)
	}

	logger.Infof("Greet response: %s", resp.Greeting)
}
