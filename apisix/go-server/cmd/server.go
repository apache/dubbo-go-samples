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
	"os"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"

	greet "github.com/apache/dubbo-go-samples/apisix/proto"

	"github.com/dubbogo/gost/log/logger"
)

type GreeterProvider struct{}

func (s *GreeterProvider) SayHello(ctx context.Context, in *greet.HelloRequest) (*greet.User, error) {
	logger.Infof("Dubbo3 GreeterProvider get user name = %s\n", in.Name)
	return &greet.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, nil
}

func (s *GreeterProvider) SayHelloStream(ctx context.Context, svr greet.Greeter_SayHelloStreamServer) error {
	c, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go3 GreeterProvider recv 1 user, name = %s\n", c.Name)
	c2, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go3 GreeterProvider recv 2 user, name = %s\n", c2.Name)

	err = svr.Send(&greet.User{
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
	logger.Infof("Dubbo-go3 GreeterProvider recv 3 user, name = %s\n", c3.Name)

	err = svr.Send(&greet.User{
		Name: "hello " + c2.Name,
		Age:  19,
		Id:   "123456789",
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Get Nacos address from environment variable, default to localhost:8848
	nacosAddr := os.Getenv("NACOS_ADDR")
	if nacosAddr == "" {
		nacosAddr = "localhost:8848"
	}

	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_apisix_server"),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress(nacosAddr),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20001),
		),
	)
	if err != nil {
		logger.Errorf("new dubbo instance failed: %v", err)
		panic(err)
	}

	srv, err := ins.NewServer()
	if err != nil {
		logger.Errorf("new server failed: %v", err)
		panic(err)
	}

	if err := greet.RegisterGreeterHandler(srv, &GreeterProvider{}); err != nil {
		logger.Errorf("register greeter handler failed: %v", err)
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Errorf("server serve failed: %v", err)
		panic(err)
	}
}
