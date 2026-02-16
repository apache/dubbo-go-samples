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
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/streaming/proto"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
	)
	if err != nil {
		panic(err) // fail fast: client not created
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err) // fail fast: service proxy not created
	}
	TestClient(svc)
}

func TestClient(cli greet.GreetService) {
	if err := testUnary(cli); err != nil {
		panic(err)
	}

	if err := testBidiStream(cli); err != nil {
		panic(err)
	}

	if err := testClientStream(cli); err != nil {
		panic(err)
	}

	if err := testServerStream(cli); err != nil {
		panic(err)
	}
}

// testUnary: 1 request -> 1 response
func testUnary(cli greet.GreetService) error {
	logger.Info("start to test TRIPLE unary call")
	resp, err := cli.Greet(context.Background(), &greet.GreetRequest{Name: "triple"})
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("unexpected unary resp: <nil>")
	}
	logger.Infof("TRIPLE unary call resp: %s", resp.Greeting)
	if resp.Greeting != "triple" && resp.Greeting != "Hello triple" {
		return fmt.Errorf("unexpected unary resp: %+v", resp)
	}
	return nil
}

// testBidiStream: N requests -> N responses (1:1)
func testBidiStream(cli greet.GreetService) error {
	logger.Info("start to test TRIPLE bidi stream")
	stream, err := cli.GreetStream(context.Background())
	if err != nil {
		return err
	}
	names := []string{"triple-1", "triple-2", "triple-3"}
	for _, name := range names {
		if sendErr := stream.Send(&greet.GreetStreamRequest{Name: name}); sendErr != nil {
			return sendErr
		}
		resp, recvErr := stream.Recv()
		if recvErr != nil {
			return recvErr
		}
		if resp == nil {
			return fmt.Errorf("unexpected bidi resp: <nil>")
		}
		expectedJava := fmt.Sprintf("Echo from biStream: %s", name)
		if resp.Greeting != name && resp.Greeting != expectedJava {
			return fmt.Errorf("unexpected bidi resp, expect %s or %s got %+v", name, expectedJava, resp)
		}
		logger.Infof("TRIPLE bidi stream resp: %s", resp.Greeting)
	}
	if err := stream.CloseRequest(); err != nil {
		return err
	}
	if err := stream.CloseResponse(); err != nil {
		return err
	}
	return nil
}

// testClientStream: 5 requests -> 1 response
func testClientStream(cli greet.GreetService) error {
	logger.Info("start to test TRIPLE client stream")
	stream, err := cli.GreetClientStream(context.Background())
	if err != nil {
		return err
	}
	for i := 0; i < 5; i++ {
		if sendErr := stream.Send(&greet.GreetClientStreamRequest{Name: "triple"}); sendErr != nil {
			return sendErr
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("unexpected client stream resp: <nil>")
	}
	logger.Infof("TRIPLE client stream resp: %s", resp.Greeting)
	expectedGo := "triple,triple,triple,triple,triple"
	expectedJavaPrefix := "Received 5 names: triple, triple, triple, triple, triple"
	if resp.Greeting != expectedGo && resp.Greeting != expectedJavaPrefix {
		return fmt.Errorf("unexpected client stream resp: %+v", resp)
	}
	return nil
}

// testServerStream: 1 request -> 10 responses
func testServerStream(cli greet.GreetService) error {
	logger.Info("start to test TRIPLE server stream")
	stream, err := cli.GreetServerStream(context.Background(), &greet.GreetServerStreamRequest{Name: "triple"})
	if err != nil {
		return err
	}
	count := 0
	const reqName = "triple"
	for stream.Recv() {
		msg := stream.Msg()
		if msg == nil {
			return fmt.Errorf("unexpected server stream msg: <nil>")
		}
		expectedGo := "triple"
		expectedJava := fmt.Sprintf("Response %d from serverStream for %s", count, reqName)
		legacyExpectedJava := fmt.Sprintf("Response %d from serverStream for StreamingClient", count)
		if msg.Greeting != expectedGo && msg.Greeting != expectedJava && msg.Greeting != legacyExpectedJava {
			return fmt.Errorf("unexpected server stream msg: %+v", msg)
		}
		count++
		logger.Infof("TRIPLE server stream resp #%d: %s", count, msg.Greeting)
	}
	if stream.Err() != nil {
		return stream.Err()
	}
	if count != 10 {
		return fmt.Errorf("unexpected server stream count, expect 10 got %d", count)
	}
	if err := stream.Close(); err != nil {
		return err
	}
	return nil
}
