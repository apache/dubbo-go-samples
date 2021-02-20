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
	pb "github.com/apache/dubbo-go-samples/general/dubbo3/protobuf/grpc"
	"go.uber.org/atomic"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":20001"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) Dubbo3SayHello2(ctx context.Context, in *pb.Dubbo3HelloRequest) (*pb.Dubbo3HelloReply, error) {
	fmt.Println("######### get server request name :" + in.Myname)
	return &pb.Dubbo3HelloReply{Msg: "Hello " + in.Myname}, nil
}

func (s *server) BigUnaryTest(ctx context.Context , in *pb.BigData) (*pb.BigData, error){
	fmt.Println("server unary recv len = ",  len(in.Data))
	return &pb.BigData{
		Data: make([]byte, in.WantSize),
	}, nil
}
var  finish = 0
var firstSend atomic.Int32
var secondSend atomic.Int32
var thirdRecv atomic.Int32
var firstRecv atomic.Int32
var secondRecv atomic.Int32

func (s *server) BigStreamTest(svr pb.Dubbo3Greeter_BigStreamTestServer ) error{
	c, err := svr.Recv()
	if err != nil {
		return err
	}
	fmt.Println("firstRecv = ", firstRecv.Inc())
	fmt.Println("server server recv 1 len  = ", len(c.Data))
	c2, err := svr.Recv()
	if err != nil {
		return err
	}
	fmt.Println("secondRecv = ", secondRecv.Inc())
	//fmt.Println("server server recv 2 = ", len)
	c3, err := svr.Recv()
	if err != nil {
		return err
	}
	fmt.Println("server server recv 3 len = ",  len(c3.Data))

	fmt.Println("thirdRecv = ", thirdRecv.Inc())
	svr.Send(&pb.BigData{
		Data: make([]byte, c.WantSize),
	})
	fmt.Println("server server send 1 len = ", c.WantSize)
	fmt.Println("firstSend = ", firstSend.Inc())
	svr.Send(&pb.BigData{
		Data: make([]byte, c2.WantSize),
	})
	fmt.Println("server server send 2 len = ", c2.WantSize)
	fmt.Println("secondSend = ", secondSend.Inc())
	fmt.Println(finish)
	return nil
}

func (s *server) Dubbo3SayHello(svr pb.Dubbo3Greeter_Dubbo3SayHelloServer) error {
	c, err := svr.Recv()
	if err != nil {
		return err
	}
	fmt.Println("server server recv 1 = ", c)
	c2, err := svr.Recv()
	if err != nil {
		return err
	}
	fmt.Println("server server recv 2 = ", c2)
	c3, err := svr.Recv()
	fmt.Println("server server recv 3 = ", c3)

	svr.Send(&pb.Dubbo3HelloReply{
		Msg: c.Myname + c2.Myname,
	})
	fmt.Println("server server send 1 = ", c.Myname+c2.Myname)
	svr.Send(&pb.Dubbo3HelloReply{
		Msg: c3.Myname,
	})
	fmt.Println("server server send 2 = ", c3.Myname)
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDubbo3GreeterServer(s, &server{})
	// Register reflection service on gRPC client.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
