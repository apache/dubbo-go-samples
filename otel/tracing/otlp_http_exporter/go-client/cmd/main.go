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
	"dubbo.apache.org/dubbo-go/v3/otel/trace"

	"github.com/dubbogo/gost/log/logger"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

import (
	greet "github.com/apache/dubbo-go-samples/otel/tracing/stdout/proto"
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_otel_client"),
		dubbo.WithTracing(
			trace.WithEnabled(),
			trace.WithOtlpHttpExporter(),
			trace.WithW3cPropagator(),
			trace.WithAlwaysMode(),
			trace.WithInsecure(),
			trace.WithEndpoint("127.0.0.1:4318"),
		),
	)
	if err != nil {
		panic(err)
	}

	// Triple
	cli1, err := ins.NewClient(client.WithClientURL("127.0.0.1:20000"))
	if err != nil {
		panic(err)
	}
	svc, err := greet.NewGreetService(cli1)
	if err != nil {
		panic(err)
	}
	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "triple"})
	if err != nil {
		logger.Error(err)
	} else {
		logger.Infof("Triple response: %s", resp.Greeting)
	}

	// Dubbo
	cli2, err := ins.NewClient(
		client.WithClientURL("127.0.0.1:20001"),
		client.WithClientProtocolDubbo(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		panic(err)
	}
	conn2, err := cli2.Dial("GreetProvider")
	if err != nil {
		panic(err)
	}
	var dubboResp string
	if err = conn2.CallUnary(context.Background(), []interface{}{"hello", "new", "dubbo"}, &dubboResp, "SayHello"); err != nil {
		logger.Error(err)
	} else {
		logger.Infof("Dubbo response: %s", dubboResp)
	}

	// JSONRPC
	cli3, err := ins.NewClient(
		client.WithClientURL("127.0.0.1:20002"),
		client.WithClientProtocolJsonRPC(),
		client.WithClientSerialization(constant.JSONSerialization),
	)
	if err != nil {
		panic(err)
	}
	conn3, err := cli3.Dial("GreetProvider")
	if err != nil {
		panic(err)
	}
	var jsonrpcResp string
	if err = conn3.CallUnary(context.Background(), []interface{}{"hello", "new", "jsonrpc"}, &jsonrpcResp, "SayHello"); err != nil {
		logger.Error(err)
	} else {
		logger.Infof("JSONRPC response: %s", jsonrpcResp)
	}

	if tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); ok {
		_ = tp.Shutdown(context.Background())
	}
}
