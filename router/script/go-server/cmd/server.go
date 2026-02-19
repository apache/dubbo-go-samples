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
	"strconv"
	"strings"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/direct/proto"
)

const (
	RegistryAddress = "127.0.0.1:8848"
	TriPort         = 20001
)

type GreetServer struct {
	srvPort int
}

func (srv *GreetServer) Greet(_ context.Context, req *greet.GreetRequest) (rep *greet.GreetResponse, err error) {
	rep = &greet.GreetResponse{Greeting: req.Name + " from: " + strconv.Itoa(srv.srvPort)}
	return rep, nil
}

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("script-server"),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress(RegistryAddress),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(TriPort),
		),
	)

	if err != nil {
		logger.Errorf("new instance failed: %v", err)
		panic(err)
	}

	srv, err := ins.NewServer()

	if err != nil {
		logger.Errorf("new server failed: %v", err)
		panic(err)
	}

	if err := greet.RegisterGreetServiceHandler(srv, &GreetServer{srvPort: TriPort}); err != nil {
		logger.Errorf("register service failed: %v", err)
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Errorf("server serve failed: %v", err)
		if strings.Contains(err.Error(), "client not connected") {
			logger.Errorf("hint: Nacos client not connected (gRPC). Check %s is reachable and gRPC port %d is open (Nacos 2.x default).", RegistryAddress, TriPort)
		}
		panic(err)
	}
}
