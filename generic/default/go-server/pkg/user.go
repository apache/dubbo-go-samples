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

package pkg

import (
	"context"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log"
)

func init() {
	config.SetProviderService(new(UserProvider))
	// ------for hessian2------
	hessian.RegisterPOJO(&User{})
	hessian.RegisterPOJO(&UserResponse{})
}

type User struct {
	ID   string
	Name string
	Age  int32
	Time time.Time
}

type UserResponse struct {
	Users []*User
}

type UserProvider struct {
}

func (u *UserProvider) GetUser(ctx context.Context, userID string) (*User, error) {
	gxlog.CInfo("req:%#v", userID)
	rsp := User{"A001", "Alex Stocks", 18, time.Now()}
	gxlog.CInfo("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) QueryUser(ctx context.Context, user *User) (*User, error) {
	gxlog.CInfo("req1:%#v", user)
	rsp := User{user.ID, user.Name, user.Age, time.Now()}
	gxlog.CInfo("rsp1:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) QueryUsers(_ context.Context, users []*User) (*UserResponse, error) {
	return &UserResponse{
		Users: users,
	}, nil
}

func (u *UserProvider) MethodMapper() map[string]string {
	return map[string]string{
		"QueryUser": "queryUser",
	}
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}

func (u *User) JavaClassName() string {
	return "org.apache.dubbo.User"
}

func (u *UserResponse) JavaClassName() string {
	return "org.apache.dubbo.UserResponse"
}
