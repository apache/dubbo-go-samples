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
		u.ID, u.Name, u.Age, time.Unix(int64(u.Time), 0).Format("2006-01-02 15:04:05.99999"), u.Sex,
	)
}

type UserProvider struct {
	GetUsers func(ids []interface{}) ([]*JsonRPCUser, error)
	GetUser  func(ctx context.Context, id string) (*JsonRPCUser, error)
	GetUser0 func(id string, name string) (*JsonRPCUser, error)
	GetUser1 func(ctx context.Context, id string) (*JsonRPCUser, error)
	GetUser2 func(ctx context.Context, id string) (*JsonRPCUser, error) `dubbo:"getUser2"`
	GetUser3 func() error
	Echo     func(ctx context.Context, req string) (string, error) // Echo represent EchoFilter will be used
}

type UserProvider1 struct {
	GetUsers func(ids []interface{}) ([]*JsonRPCUser, error)
	GetUser  func(ctx context.Context, id string) (*JsonRPCUser, error)
	GetUser0 func(id string, name string) (*JsonRPCUser, error)
	GetUser1 func(ctx context.Context, id string) (*JsonRPCUser, error)
	GetUser2 func(ctx context.Context, id string) (*JsonRPCUser, error) `dubbo:"getUser2"`
	GetUser3 func() error
	Echo     func(ctx context.Context, req string) (string, error) // Echo represent EchoFilter will be used
}

type UserProvider2 struct {
	GetUsers func(ids []interface{}) ([]*JsonRPCUser, error)
	GetUser  func(ctx context.Context, id string) (*JsonRPCUser, error)
	GetUser0 func(id string, name string) (*JsonRPCUser, error)
	GetUser1 func(ctx context.Context, id string) (*JsonRPCUser, error)
	GetUser2 func(ctx context.Context, id string) (*JsonRPCUser, error) `dubbo:"getUser2"`
	GetUser3 func() error
	Echo     func(ctx context.Context, req string) (string, error) // Echo represent EchoFilter will be used
}
