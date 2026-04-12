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
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	triple "dubbo.apache.org/dubbo-go/v3/protocol/triple"
	"dubbo.apache.org/dubbo-go/v3/server"
)

import (
	demo "github.com/apache/dubbo-go-samples/rpc/triple/openapi/proto/demo"
	greet "github.com/apache/dubbo-go-samples/rpc/triple/openapi/proto/greet"
)

type GreetTripleServer struct{}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: "Hello, " + req.Name}
	return resp, nil
}

func (srv *GreetTripleServer) GreetServerStream(ctx context.Context, req *greet.GreetServerStreamRequest, stream greet.GreetService_GreetServerStreamServer) error {
	for i := 0; i < 5; i++ {
		if err := stream.Send(&greet.GreetServerStreamResponse{Greeting: "Hello, " + req.Name}); err != nil {
			return err
		}
	}
	return nil
}

func (srv *GreetTripleServer) GreetClientStream(ctx context.Context, stream greet.GreetService_GreetClientStreamServer) (*greet.GreetClientStreamResponse, error) {
	var names []string
	for stream.Recv() {

		msg := stream.Msg()
		names = append(names, msg.Name)
	}
	return &greet.GreetClientStreamResponse{Greeting: "Hello, " + joinNames(names)}, nil
}

func (srv *GreetTripleServer) GreetBidiStream(ctx context.Context, stream greet.GreetService_GreetBidiStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		if err := stream.Send(&greet.GreetBidiStreamResponse{Greeting: "Hello, " + req.Name}); err != nil {
			return err
		}
	}
}

func joinNames(names []string) string {
	result := ""
	for i, n := range names {
		if i > 0 {
			result += ", "
		}
		result += n
	}
	return result
}

type DemoTripleServerV1 struct{}

func (srv *DemoTripleServerV1) Greet(ctx context.Context, req *demo.GreetRequest) (*demo.GreetResponse, error) {
	resp := &demo.GreetResponse{Greeting: "Hello, " + req.Name}
	return resp, nil
}

type DemoTripleServerV2 struct{}

func (srv *DemoTripleServerV2) Greet(ctx context.Context, req *demo.GreetRequest) (*demo.GreetResponse, error) {
	resp := &demo.GreetResponse{Greeting: "Hello, " + req.Name}
	return resp, nil
}

// Non-proto service (non-IDL mode)
type UserService struct{}

type UserRequest struct {
	Id int32 `json:"id"`
}

type UserResponse struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

func (u *UserService) GetUser(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	return &UserResponse{
		Id:   req.Id,
		Name: "Alice",
		Age:  30,
	}, nil
}

func (u *UserService) Reference() string {
	return "com.example.UserService"
}

func main() {
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithTriple(
				triple.WithOpenAPI(
					triple.OpenAPIEnable(),
					triple.OpenAPIInfoTitle("OpenAPI Service"),
					triple.OpenAPIInfoDescription("A service with OpenAPI documentation"),
					triple.OpenAPIInfoVersion("1.0.0"),
				),
			),
			protocol.WithPort(20000),
		),
	)
	if err != nil {
		panic(err)
	}
	if registerErr := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); registerErr != nil {
		panic(registerErr)
	}
	if registerErr := demo.RegisterGreetServiceHandler(srv, &DemoTripleServerV1{},
		server.WithVersion("1.0.0"),
	); registerErr != nil {
		panic(registerErr)
	}
	if registerErr := demo.RegisterGreetServiceHandler(srv, &DemoTripleServerV2{},
		server.WithOpenAPIGroup("demo-2.0.0"),
		server.WithVersion("2.0.0"),
	); registerErr != nil {
		panic(registerErr)
	}

	if err := srv.RegisterService(&UserService{}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		panic(err)
	}
}
