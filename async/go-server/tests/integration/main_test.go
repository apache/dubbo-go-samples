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
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/metadata/service/inmemory"
	"github.com/apache/dubbo-go/protocol"
	"github.com/apache/dubbo-go/remoting"
	gxlog "github.com/dubbogo/gost/log"
)

import (
	"context"
	"os"
	"testing"
	"time"
)

var userProvider = &UserProvider{
	ch: make(chan *User),
}

func TestMain(m *testing.M) {
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&User{})
	config.Load()
	time.Sleep(3 * time.Second)

	os.Exit(m.Run())
}

type User struct {
	Id   string
	Name string
	Age  int32
	Time time.Time
}

type UserProvider struct {
	GetUser func(ctx context.Context, req []interface{}, rsp *User) error
	ch      chan *User
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}

// to enable async call:
// 1. need to implement AsyncCallbackService
// 2. need to specify references -> UserProvider -> async in conf/client.yml
func (u *UserProvider) CallBack(res common.CallbackResponse) {
	gxlog.CInfo("CallBack res: %v", res)
	if r, ok := res.(remoting.AsyncCallbackResponse); ok {
		if reply, ok := r.Reply.(*remoting.Response); ok {
			if result, ok := reply.Result.(*protocol.RPCResult); ok {
				if user, ok := result.Rest.(*User); ok {
					u.ch <- user
				}
			}
		}
	}
	u.ch <- nil
}

func (User) JavaClassName() string {
	return "org.apache.dubbo.User"
}
