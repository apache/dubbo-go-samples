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
	"strings"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	triple "dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
	"dubbo.apache.org/dubbo-go/v3/server"

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

type GreetTripleServer struct{}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	token, modes := incomingMetadata(ctx)
	logger.Infof("unary request metadata: token=%q modes=%v", token, modes)

	return &greet.GreetResponse{
		Greeting: fmt.Sprintf("hello %s, token=%s, modes=%s", req.Name, token, strings.Join(modes, ",")),
	}, nil
}

func (srv *GreetTripleServer) GreetStream(ctx context.Context, stream greet.GreetService_GreetStreamServer) error {
	token := firstValue(stream.RequestHeader(), tokenHeader)
	if token == "" {
		token = firstIncomingValue(ctx, tokenHeader)
	}
	logger.Infof("bidi request header %s=%q", tokenHeader, token)

	stream.ResponseHeader().Set(streamResponseKey, "bidi-header")
	stream.ResponseTrailer().Set(streamTrailerKey, "bidi-trailer")

	for {
		req, err := stream.Recv()
		if err != nil {
			if triple.IsEnded(err) {
				return nil
			}
			return err
		}
		if err := stream.Send(&greet.GreetStreamResponse{
			Greeting: fmt.Sprintf("bidi hello %s, token=%s", req.Name, token),
		}); err != nil {
			return err
		}
	}
}

func (srv *GreetTripleServer) GreetClientStream(ctx context.Context, stream greet.GreetService_GreetClientStreamServer) (*greet.GreetClientStreamResponse, error) {
	token := firstValue(stream.RequestHeader(), tokenHeader)
	if token == "" {
		token = firstIncomingValue(ctx, tokenHeader)
	}
	logger.Infof("client-stream request header %s=%q", tokenHeader, token)

	var names []string
	for stream.Recv() {
		names = append(names, stream.Msg().Name)
	}
	if stream.Err() != nil && !triple.IsEnded(stream.Err()) {
		return nil, stream.Err()
	}

	return &greet.GreetClientStreamResponse{
		Greeting: fmt.Sprintf("client-stream hello %s, token=%s", strings.Join(names, ","), token),
	}, nil
}

func (srv *GreetTripleServer) GreetServerStream(ctx context.Context, req *greet.GreetServerStreamRequest, stream greet.GreetService_GreetServerStreamServer) error {
	token := firstIncomingValue(ctx, tokenHeader)
	logger.Infof("server-stream request metadata %s=%q", tokenHeader, token)

	stream.ResponseHeader().Set(streamResponseKey, "server-stream-header")
	stream.ResponseTrailer().Set(streamTrailerKey, "server-stream-trailer")

	for i := 0; i < 3; i++ {
		if err := stream.Send(&greet.GreetServerStreamResponse{
			Greeting: fmt.Sprintf("server-stream hello %s #%d, token=%s", req.Name, i+1, token),
		}); err != nil {
			return err
		}
	}
	return nil
}

func incomingMetadata(ctx context.Context) (string, []string) {
	headers, ok := triple.FromIncomingContext(ctx)
	if !ok {
		return "", nil
	}
	return firstValue(headers, tokenHeader), headers.Values(modeHeader)
}

func firstIncomingValue(ctx context.Context, key string) string {
	headers, ok := triple.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	return firstValue(headers, key)
}

func firstValue(headers http.Header, key string) string {
	values := headers.Values(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func main() {
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000),
			protocol.WithTriple(),
		),
	)
	if err != nil {
		panic(err)
	}

	if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
