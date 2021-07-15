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
	"fmt"
	"time"
)

import (
	"github.com/apache/dubbo-go/config"
)

var (
	UserProviderVar  = new(UserProvider)
	UserProvider1Var = new(UserProvider1)
	UserProvider2Var = new(UserProvider2)
)

func init() {
	config.SetConsumerService(UserProviderVar)
	config.SetConsumerService(UserProvider1Var)
	config.SetConsumerService(UserProvider2Var)
}

type User struct {
	ID   string
	Name string
	Age  int64
	Time int64
	Sex  string
}

func (u User) String() string {
	return fmt.Sprintf(
		"User{ID:%s, Name:%s, Age:%d, Time:%s, Sex:%s}",
		u.ID, u.Name, u.Age, time.Unix(u.Time, 0).Format("2006-01-02 15:04:05.99999"), u.Sex,
	)
}

type UserProvider struct {
	GetUsers func(req []interface{}) ([]User, error)
	GetUser  func(ctx context.Context, req []interface{}, rsp *User) error
	GetUser0 func(id string, name string, age int) (*User, error)
	GetUser1 func(ctx context.Context, req []interface{}, rsp *User) error
	GetUser2 func(ctx context.Context, req []interface{}, rsp *User) error `dubbo:"getUser"`
	GetUser3 func() error
	Echo     func(ctx context.Context, req interface{}) (interface{}, error) // Echo represent EchoFilter will be used
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}

type UserProvider1 struct {
	GetUsers func(req []interface{}) ([]User, error)
	GetUser  func(ctx context.Context, req []interface{}, rsp *User) error
	GetUser0 func(id string, name string, age int) (User, error)
	GetUser1 func(ctx context.Context, req []interface{}, rsp *User) error
	GetUser2 func(ctx context.Context, req []interface{}, rsp *User) error `dubbo:"getUser"`
	GetUser3 func() error
	Echo     func(ctx context.Context, req interface{}) (interface{}, error) // Echo represent EchoFilter will be used
}

func (u *UserProvider1) Reference() string {
	return "UserProvider1"
}

type UserProvider2 struct {
	GetUsers func(req []interface{}) ([]User, error)
	GetUser  func(ctx context.Context, req []interface{}, rsp *User) error
	GetUser0 func(id string, name string, age int) (User, error)
	GetUser1 func(ctx context.Context, req []interface{}, rsp *User) error
	GetUser2 func(ctx context.Context, req []interface{}, rsp *User) error `dubbo:"getUser"`
	GetUser3 func() error
	Echo     func(ctx context.Context, req interface{}) (interface{}, error) // Echo represent EchoFilter will be used
}

func (u *UserProvider2) Reference() string {
	return "UserProvider2"
}
