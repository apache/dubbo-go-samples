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
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"

	gxlog "github.com/dubbogo/gost/log"
	"github.com/dubbogo/gost/log/logger"
)

import (
	user "github.com/apache/dubbo-go-samples/async/proto"
)

// UserProviderHandler implements proto.UserProviderHandler
type UserProviderHandler struct {
}

func (h *UserProviderHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	gxlog.CInfo("req:%#v", req)

	u, ok := user.GetUserByID(req.Id)
	if !ok {
		return nil, fmt.Errorf("invalid user id:%s", req.Id)
	}

	gxlog.CInfo("rsp:%#v", u)
	return &user.GetUserResponse{User: u}, nil
}

// UserProviderV2Handler implements proto.UserProviderV2Handler
type UserProviderV2Handler struct {
}

func (h *UserProviderV2Handler) SayHello(ctx context.Context, req *user.SayHelloRequest) (*user.SayHelloResponse, error) {
	u, ok := user.GetUserByID(req.UserId)
	if !ok {
		return nil, fmt.Errorf("invalid user id:%s", req.UserId)
	}

	gxlog.CInfo("hello, %s", u.Name)
	return &user.SayHelloResponse{}, nil
}

func main() {
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20000),
		),
	)
	if err != nil {
		panic(err)
	}

	if err := user.RegisterUserProviderHandler(srv, &UserProviderHandler{}); err != nil {
		panic(err)
	}

	if err := user.RegisterUserProviderV2Handler(srv, &UserProviderV2Handler{}); err != nil {
		panic(err)
	}

	logger.Infof("async dubbo server started on port 20000")

	if err := srv.Serve(); err != nil {
		panic(err)
	}
}
