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
)

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

import (
	pb "github.com/apache/dubbo-go-samples/rpc/triple/pb/dubbogo-grpc/protobuf/api"
)

const (
	address = "localhost:20000"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := pb.NewGreeterClient(conn)

	defer func() {
		_ = conn.Close()
	}()

	stream(c)
	unary(c)
}

func stream(c pb.GreeterClient) {
	fmt.Printf(">>>>> gRPC-go client is about to call SayHelloStream\n")

	clientStream, err := c.SayHelloStream(context.Background())
	if err != nil {
		panic(err)
	}

	BigDataReq := &pb.HelloRequest{
		Name: "Laurence",
	}

	for i := 0; i < 2; i++ {
		_ = clientStream.Send(BigDataReq)
	}
	user1, err := clientStream.Recv()
	if err != nil {
		panic(err)
	}
	fmt.Printf("get 1 received user = %+v\n", user1)

	_ = clientStream.Send(BigDataReq)

	user2, err := clientStream.Recv()
	if err != nil {
		panic(err)
	}
	fmt.Printf("get 2 received user = %+v\n", user2)
}

func unary(c pb.GreeterClient) {
	fmt.Printf(">>>>> gRPC-go client is about to call SayHello\n")

	req := &pb.HelloRequest{
		Name: "laurence",
	}
	ctx := context.Background()
	rsp, err := c.SayHello(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("get received user = %+v\n", rsp)
}
