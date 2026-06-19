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
	"net/http"
)

import (
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	triple "dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/triple_header_trailer/proto"
)

const (
	tokenHeader       = "X-Sample-Token"
	modeHeader        = "X-Sample-Mode"
	streamResponseKey = "X-Stream-Response"
	streamTrailerKey  = "X-Stream-Trailer"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	if err := testUnary(svc); err != nil {
		panic(err)
	}
	if err := testBidiStream(svc); err != nil {
		panic(err)
	}
	if err := testClientStream(svc); err != nil {
		panic(err)
	}
	if err := testServerStream(svc); err != nil {
		panic(err)
	}
}

func testUnary(cli greet.GreetService) error {
	var responseHeader http.Header
	var responseTrailer http.Header
	resp, err := cli.Greet(
		outgoingContext("unary-token"),
		&greet.GreetRequest{Name: "unary"},
		client.WithResponseHeader(&responseHeader),
		client.WithResponseTrailer(&responseTrailer),
	)
	if err != nil {
		return err
	}
	if err := requireGreeting(resp.Greeting, "hello unary, token=unary-token, modes=metadata,trailer-demo"); err != nil {
		return err
	}
	logger.Infof("unary response: %s", resp.Greeting)
	logger.Infof("unary response header content-type=%q", responseHeader.Get("Content-Type"))
	logger.Infof("unary response trailer grpc-status=%q", responseTrailer.Get("Grpc-Status"))
	if responseHeader.Get("Content-Type") == "" {
		return fmt.Errorf("missing unary response content-type header in %v", responseHeader)
	}
	return nil
}

func testBidiStream(cli greet.GreetService) error {
	stream, err := cli.GreetStream(outgoingContext("bidi-token"))
	if err != nil {
		return err
	}

	if err := stream.Send(&greet.GreetStreamRequest{Name: "bidi"}); err != nil {
		return err
	}
	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("unexpected empty bidi response")
	}
	if err := requireGreeting(resp.Greeting, "bidi hello bidi, token=bidi-token"); err != nil {
		return err
	}
	logger.Infof("bidi response: %s", resp.Greeting)
	logger.Infof("bidi response header %s=%v", streamResponseKey, stream.ResponseHeader().Values(streamResponseKey))

	if err := stream.CloseRequest(); err != nil {
		return err
	}
	if err := stream.CloseResponse(); err != nil {
		return err
	}
	logger.Infof("bidi response trailer %s=%v", streamTrailerKey, stream.ResponseTrailer().Values(streamTrailerKey))
	return requireHeader(stream.ResponseTrailer(), streamTrailerKey, "bidi-trailer")
}

func testClientStream(cli greet.GreetService) error {
	stream, err := cli.GreetClientStream(outgoingContext("client-stream-token"))
	if err != nil {
		return err
	}

	for _, name := range []string{"client", "stream"} {
		if err := stream.Send(&greet.GreetClientStreamRequest{Name: name}); err != nil {
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	if err := requireGreeting(resp.Greeting, "client-stream hello client,stream, token=client-stream-token"); err != nil {
		return err
	}
	logger.Infof("client-stream response: %s", resp.Greeting)
	return nil
}

func testServerStream(cli greet.GreetService) error {
	stream, err := cli.GreetServerStream(outgoingContext("server-stream-token"), &greet.GreetServerStreamRequest{Name: "server-stream"})
	if err != nil {
		return err
	}

	var count int
	expectedGreetings := []string{
		"server-stream hello server-stream #1, token=server-stream-token",
		"server-stream hello server-stream #2, token=server-stream-token",
		"server-stream hello server-stream #3, token=server-stream-token",
	}
	for stream.Recv() {
		if count >= len(expectedGreetings) {
			return fmt.Errorf("unexpected extra server-stream response: %s", stream.Msg().Greeting)
		}
		if err := requireGreeting(stream.Msg().Greeting, expectedGreetings[count]); err != nil {
			return err
		}
		count++
		logger.Infof("server-stream response #%d: %s", count, stream.Msg().Greeting)
	}
	if stream.Err() != nil {
		return stream.Err()
	}
	if count != len(expectedGreetings) {
		return fmt.Errorf("unexpected server-stream response count: %d", count)
	}
	if err := stream.Close(); err != nil {
		return err
	}

	logger.Infof("server-stream response header %s=%v", streamResponseKey, stream.ResponseHeader().Values(streamResponseKey))
	logger.Infof("server-stream response trailer %s=%v", streamTrailerKey, stream.ResponseTrailer().Values(streamTrailerKey))
	if err := requireHeader(stream.ResponseHeader(), streamResponseKey, "server-stream-header"); err != nil {
		return err
	}
	return requireHeader(stream.ResponseTrailer(), streamTrailerKey, "server-stream-trailer")
}

func outgoingContext(token string) context.Context {
	ctx := triple.NewOutgoingContext(context.Background(), http.Header{
		tokenHeader: []string{token},
	})
	return triple.AppendToOutgoingContext(ctx, modeHeader, "metadata", modeHeader, "trailer-demo")
}

func requireGreeting(got string, want string) error {
	if got != want {
		return fmt.Errorf("unexpected response: got %q, want %q", got, want)
	}
	return nil
}

func requireHeader(headers http.Header, key string, want string) error {
	for _, got := range headers.Values(key) {
		if got == want {
			return nil
		}
	}
	return fmt.Errorf("missing header %s=%q in %v", key, want, headers)
}
