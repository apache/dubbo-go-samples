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
	"dubbo.apache.org/dubbo-go/v3/config"

	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/service/local"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
)

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

var userProvider = new(UserProvider)

func TestMain(m *testing.M) {
	config.SetConsumerService(userProvider)
	hessian.RegisterJavaEnum(Gender(MAN))
	hessian.RegisterJavaEnum(Gender(WOMAN))
	hessian.RegisterPOJO(&User{})
	config.Load()
	time.Sleep(3 * time.Second)

	os.Exit(m.Run())
}

type Gender hessian.JavaEnum

func init() {
	config.SetConsumerService(userProvider)
}

const (
	MAN hessian.JavaEnum = iota
	WOMAN
)

var genderName = map[hessian.JavaEnum]string{
	MAN:   "MAN",
	WOMAN: "WOMAN",
}

var genderValue = map[string]hessian.JavaEnum{
	"MAN":   MAN,
	"WOMAN": WOMAN,
}

func (g Gender) JavaClassName() string {
	return "org.apache.dubbo.Gender"
}

func (g Gender) String() string {
	s, ok := genderName[hessian.JavaEnum(g)]
	if ok {
		return s
	}

	return strconv.Itoa(int(g))
}

func (g Gender) EnumValue(s string) hessian.JavaEnum {
	v, ok := genderValue[s]
	if ok {
		return v
	}

	return hessian.InvalidJavaEnum
}

type User struct {
	// !!! Cannot define lowercase names of variable
	ID   string
	Name string
	Age  int32
	Time time.Time
	Sex  Gender // notice: java enum Object <--> go string
}

func (u User) String() string {
	return fmt.Sprintf(
		"User{ID:%s, Name:%s, Age:%d, Time:%s, Sex:%s}",
		u.ID, u.Name, u.Age, u.Time, u.Sex,
	)
}

func (User) JavaClassName() string {
	return "org.apache.dubbo.User"
}

type UserProvider struct {
	GetUsers  func(req []interface{}) ([]interface{}, error)
	GetErr    func(ctx context.Context, req []interface{}, rsp *User) error
	GetUser   func(ctx context.Context, req []interface{}, rsp *User) error
	GetUser0  func(id string, name string) (User, error)
	GetUser1  func(ctx context.Context, req []interface{}, rsp *User) error
	GetUser2  func(ctx context.Context, req []interface{}, rsp *User) error `dubbo:"getUser"`
	GetUser3  func() error
	GetGender func(i int32) (Gender, error)
	Echo      func(ctx context.Context, req interface{}) (interface{}, error) // Echo represent EchoFilter will be used
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}
