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
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/remoting"
	"fmt"
	gxlog "github.com/dubbogo/gost/log"
	"strconv"
	"time"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
)

type Gender hessian.JavaEnum

const (
	MAN Gender = iota
	WOMAN
)

var genderName = map[Gender]string{
	MAN:   "MAN",
	WOMAN: "WOMAN",
}

var genderValue = map[string]Gender{
	"MAN":   MAN,
	"WOMAN": WOMAN,
}

func (g Gender) JavaClassName() string {
	return "org.apache.dubbo.sample.Gender"
}

func (g Gender) String() string {
	s, ok := genderName[g]
	if ok {
		return s
	}

	return strconv.Itoa(int(g))
}

func (g Gender) EnumValue(s string) hessian.JavaEnum {
	v, ok := genderValue[s]
	if ok {
		return hessian.JavaEnum(v)
	}

	return hessian.InvalidJavaEnum
}

type User struct {
	// !!! Cannot define lowercase names of variable
	ID   string `hessian:"id"`
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

func (u *User) JavaClassName() string {
	return "org.apache.dubbo.sample.User"
}

type UserProvider struct {
	GetUser func(ctx context.Context, req *User) (*User, error)
}

// CallBack for async callback
func (u *UserProvider) CallBack(response common.CallbackResponse) {
	resp := response.(remoting.AsyncCallbackResponse).Reply.(*remoting.Response)
	if resp.Error != nil {
		panic(resp.Error)
	}

	user := resp.Result.(*protocol.RPCResult).Rest.(*User)
	gxlog.CInfo("async callback response! got result :%v\n", user)
}

// UserProviderV2 for one way
type UserProviderV2 struct {
	SayHello func(ctx context.Context, name string) error
}
