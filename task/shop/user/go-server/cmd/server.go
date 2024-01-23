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
	"time"

	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"github.com/dubbogo/gost/log/logger"

	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/apache/dubbo-go-samples/task/shop/user/api"
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("shop-user"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20013),
		),
	)
	if err != nil {
		panic(err)
	}

	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	if err = api.RegisterUserServiceHandler(srv, &UserProvider{}); err != nil {
		panic(err)
	}
	if err = srv.Serve(); err != nil {
		logger.Error(err)
	}
}

// UserProvider is the provider of user service
type UserProvider struct {
	count int
}

// Register registers a user
func (u *UserProvider) Register(ctx context.Context, req *api.User) (*api.RegisterResp, error) {
	return &api.RegisterResp{
		Success: true,
	}, nil
}

// Login gets the user
func (u *UserProvider) Login(ctx context.Context, req *api.LoginReq) (*api.User, error) {
	return &api.User{
		Username: req.Username,
		Password: req.Password,
		Phone:    "11111111111",
		Mail:     "dubbo@dubbo",
		RealName: "dubbo_test",
	}, nil
}

func (u *UserProvider) TimeoutLogin(ctx context.Context, req *api.LoginReq) (*api.User, error) {
	time.Sleep(3 * time.Second)
	return &api.User{
		Username: req.Username,
		Password: req.Password,
		Phone:    "11111111111",
		Mail:     "dubbo@dubbo",
		RealName: "dubbo_test",
	}, nil
}

func (u *UserProvider) GetInfo(ctx context.Context, req *api.GetInfoReq) (*api.User, error) {
	fmt.Println("Received getInfo request......")
	u.count++
	if u.count%3 == 0 {
		time.Sleep(3 * time.Second)
	}
	return &api.User{
		Username: req.Username,
		Password: "password",
		Phone:    "11111111111",
		Mail:     "dubbo@dubbo",
		RealName: "dubbo_test",
	}, nil
}
