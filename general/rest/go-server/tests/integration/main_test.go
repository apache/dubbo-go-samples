// +build integration

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

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go-samples/general/rest/go-client/pkg"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/service/local"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/rest"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
)

var (
	userProvider = new(pkg.UserProvider)
	userProvider1 = new(pkg.UserProvider1)
	userProvider2 = new(pkg.UserProvider2)
	)

func TestMain(m *testing.M) {
	config.SetConsumerService(userProvider)
	config.SetConsumerService(userProvider1)
	config.SetConsumerService(userProvider2)
	hessian.RegisterPOJO(&pkg.UserProvider{})
	config.Load()
	time.Sleep(3 * time.Second)

	os.Exit(m.Run())
}


type JsonRPCUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int64  `json:"age"`
	Time int64  `json:"time"`
	Sex  string `json:"sex"`
}

func (u JsonRPCUser) String() string {
	return fmt.Sprintf(
		"User{ID:%s, Name:%s, Age:%d, Time:%s, Sex:%s}",
		u.ID, u.Name, u.Age, time.Unix(u.Time, 0).Format("2006-01-02 15:04:05.99999"), u.Sex,
	)
}

type UserProvider struct {
	GetUsers func(req []interface{}) ([]JsonRPCUser, error)
	GetUser  func(ctx context.Context, req []interface{}, rsp *JsonRPCUser) error
	GetUser0 func(id string, name string) (JsonRPCUser, error)
	GetUser1 func(ctx context.Context, req []interface{}, rsp *JsonRPCUser) error
	GetUser2 func(ctx context.Context, req []interface{}, rsp *JsonRPCUser) error `dubbo:"getUser"`
	GetUser3 func() error
	Echo     func(ctx context.Context, req interface{}) (interface{}, error) // Echo represent EchoFilter will be used
}

func (u *UserProvider) Reference() string {
	return "com.ikurento.UserProvider"
}

type UserProvider1 struct {
	GetUsers func(req []interface{}) ([]JsonRPCUser, error)
	GetUser  func(ctx context.Context, req []interface{}, rsp *JsonRPCUser) error
	GetUser0 func(id string, name string) (JsonRPCUser, error)
	GetUser1 func(ctx context.Context, req []interface{}, rsp *JsonRPCUser) error
	GetUser2 func(ctx context.Context, req []interface{}, rsp *JsonRPCUser) error `dubbo:"getUser"`
	GetUser3 func() error
	Echo     func(ctx context.Context, req interface{}) (interface{}, error) // Echo represent EchoFilter will be used
}

func (u *UserProvider1) Reference() string {
	return "com.ikurento.UserProvider1"
}

type UserProvider2 struct {
	GetUsers func(req []interface{}) ([]JsonRPCUser, error)
	GetUser  func(ctx context.Context, req []interface{}, rsp *JsonRPCUser) error
	GetUser0 func(id string, name string) (JsonRPCUser, error)
	GetUser1 func(ctx context.Context, req []interface{}, rsp *JsonRPCUser) error
	GetUser2 func(ctx context.Context, req []interface{}, rsp *JsonRPCUser) error `dubbo:"getUser"`
	GetUser3 func() error
	Echo     func(ctx context.Context, req interface{}) (interface{}, error) // Echo represent EchoFilter will be used
}

func (u *UserProvider2) Reference() string {
	return "com.ikurento.UserProvider2"
}
