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
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/rpc/triple/openapi/proto/greet"
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

	// Unary
	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "openapi"})
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("Greet response: %s", resp.Greeting)

	// Server Stream
	serverStream, err := svc.GreetServerStream(context.Background(), &greet.GreetServerStreamRequest{Name: "openapi"})
	if err != nil {
		logger.Error(err)
		return
	}
	for serverStream.Recv() {
		msg := serverStream.Msg()
		logger.Infof("GreetServerStream response: %s", msg.Greeting)
	}

	// Client Stream
	clientStream, err := svc.GreetClientStream(context.Background())
	if err != nil {
		logger.Error(err)
		return
	}
	for _, name := range []string{"alice", "bob", "charlie"} {
		if sendErr := clientStream.Send(&greet.GreetClientStreamRequest{Name: name}); sendErr != nil {
			logger.Error(sendErr)
			return
		}
	}
	clientResp, err := clientStream.CloseAndRecv()
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("GreetClientStream response: %s", clientResp.Greeting)

	// Bidi Stream
	bidiStream, err := svc.GreetBidiStream(context.Background())
	if err != nil {
		logger.Error(err)
		return
	}
	for _, name := range []string{"dave", "eve"} {
		if err := bidiStream.Send(&greet.GreetBidiStreamRequest{Name: name}); err != nil {
			logger.Error(err)
			return
		}
	}
	if err := bidiStream.CloseRequest(); err != nil {
		logger.Error(err)
		return
	}
	for {
		msg, err := bidiStream.Recv()
		if err != nil {
			break
		}
		logger.Infof("GreetBidiStream response: %s", msg.Greeting)
	}
}
