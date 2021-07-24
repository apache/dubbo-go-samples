/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
)

import (
	"github.com/apache/dubbo-go/common/logger"
	"github.com/apache/dubbo-go/config"
)

func init() {
	config.SetProviderService(NewGreeterProvider())
}

type GreeterProvider struct {
	*GreeterProviderBase
}

func NewGreeterProvider() *GreeterProvider {
	return &GreeterProvider{
		GreeterProviderBase: &GreeterProviderBase{},
	}
}

func (g *GreeterProvider) SayHelloTwoSidesStream(server Greeter_SayHelloTwoSidesStreamServer) error {
	req := &HelloRequest{}
	fmt.Printf("in SayHelloTwoSidesStream func\n")
	if err := server.RecvMsg(req); err != nil {
		fmt.Printf("server.RecvMsg err = %s", err.Error())
		return err
	}
	fmt.Printf("req1: %v\n", *req)
	if err := server.RecvMsg(req); err != nil {
		logger.Errorf("server.RecvMsg err = %s", err.Error())
		return err
	}
	fmt.Printf("req2: %v\n", *req)
	if err := server.Send(&HelloReply{Message: "reply1"}); err != nil {
		logger.Errorf("server.Send err = %s", err.Error())
		return err
	}
	if err := server.Send(&HelloReply{Message: "reply2"}); err != nil {
		logger.Errorf("server.Send err = %s", err.Error())
		return err
	}
	return nil
}

func (g *GreeterProvider) SayHelloClientStream(server Greeter_SayHelloClientStreamServer) error {
	req := &HelloRequest{}
	fmt.Printf("in SayHelloClientStream func\n")
	if err := server.RecvMsg(req); err != nil {
		fmt.Printf("server.RecvMsg err = %s", err.Error())
		return err
	}
	fmt.Printf("req1: %v\n", *req)
	if err := server.RecvMsg(req); err != nil {
		logger.Errorf("server.RecvMsg err = %s", err.Error())
		return err
	}
	fmt.Printf("req2: %v\n", *req)
	if err := server.SendMsg(&HelloReply{Message: "unary reply"}); err != nil {
		logger.Errorf("server.Send err = %s", err.Error())
		return err
	}
	return nil
}

func (g *GreeterProvider) SayHelloServerStream(req *HelloRequest, server Greeter_SayHelloServerStreamServer) error {
	fmt.Println("in SayHelloServerStream func")
	fmt.Printf("unary req: %v\n", *req)
	if err := server.Send(&HelloReply{Message: "reply1"}); err != nil {
		logger.Errorf("server.Send err = %s", err.Error())
		return err
	}
	if err := server.Send(&HelloReply{Message: "reply2"}); err != nil {
		logger.Errorf("server.Send err = %s", err.Error())
		return err
	}
	return nil
}
