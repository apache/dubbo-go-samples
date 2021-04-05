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
	"fmt"
	"log"
)
import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

import (
	pb "github.com/apache/dubbo-go-samples/general/dubbo3/pb/protobuf/grpc"
)

const (
	address = "localhost:20001"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	BigDataReq := &pb.HelloRequest{
		Name: "Laurence",
	}
	clientStream, err := c.SayHelloStream(context.Background())
	if err != nil {
		panic(err)
	}
	for i := 0; i < 3; i++ {
		clientStream.Send(BigDataReq)
	}
	user1, err := clientStream.Recv()
	if err != nil {
		panic(err)
	}
	fmt.Printf("get 1 received user = %+v\n", user1)

	user2, err := clientStream.Recv()
	if err != nil {
		panic(err)
	}
	fmt.Printf("get 2 received user = %+v\n", user2)
}
