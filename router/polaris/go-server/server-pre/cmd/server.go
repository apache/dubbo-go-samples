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
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	user "github.com/apache/dubbo-go-samples/router/polaris/proto"
)

type UserServiceHandler struct {
}

func (s *UserServiceHandler) GetUser(ctx context.Context, req *user.User) (*user.GetUserResponse, error) {
	logger.Infof("req:%#v", req)
	resp := &user.GetUserResponse{
		User: &user.User{
			Id:   "A001",
			Name: "[Pre] User",
			Age:  18,
			Time: time.Now().Unix(),
		},
	}
	logger.Infof("resp:%#v", resp)
	return resp, nil
}

func main() {
	polarisAddr := "127.0.0.1:8091"

	ins, err := dubbo.NewInstance(
		dubbo.WithName("myApp"),
		dubbo.WithEnvironment("pre"),
		dubbo.WithRegistry(
			registry.WithPolaris(),
			registry.WithAddress(polarisAddr),
			registry.WithNamespace("dubbogo"),
			registry.WithRegisterInterface(),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(21000),
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

	if err := user.RegisterUserServiceHandler(srv, &UserServiceHandler{}); err != nil {
		logger.Errorf("register user service handler failed: %v", err)
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Errorf("server serve failed: %v", err)
		panic(err)
	}
}
