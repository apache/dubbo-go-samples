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
	"log"
	"net"
)

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

import (
	pb "github.com/apache/dubbo-go-samples/compatibility/rpc/triple/pb/dubbogo-grpc/protobuf/api"
)

const (
	port = ":20000"
)

// GreeterProvider is used as provider
type GreeterProvider struct {
	pb.UnimplementedGreeterServer
}

func (s *GreeterProvider) SayHelloStream(svr pb.Greeter_SayHelloStreamServer) error {
	c, err := svr.Recv()
	if err != nil {
		return err
	}
	fmt.Printf("grpc GreeterProvider recv 1 user, name = %s\n", c.Name)
	c2, err := svr.Recv()
	if err != nil {
		return err
	}
	fmt.Printf("grpc GreeterProvider recv 2 user, name = %s\n", c2.Name)

	err = svr.Send(&pb.User{
		Name: "hello " + c.Name,
		Age:  18,
		Id:   "123456789",
	})
	if err != nil {
		return err
	}
	c3, err := svr.Recv()
	if err != nil {
		return err
	}
	fmt.Printf("grpc GreeterProvider recv 3 user, name = %s\n", c3.Name)
	err = svr.Send(&pb.User{
		Name: "hello " + c2.Name,
		Age:  19,
		Id:   "123456789",
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *GreeterProvider) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.User, error) {
	fmt.Printf("Dubbo3 GreeterProvider get user name = %s\n", in.Name)
	return &pb.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &GreeterProvider{})
	// Register reflection service on gRPC client.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
