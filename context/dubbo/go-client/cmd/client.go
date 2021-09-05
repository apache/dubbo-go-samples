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
	"os"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log"
)

type UserProvider struct {
	GetUser func(ctx context.Context, req *User) (rsp *User, err error)
}

func (u *UserProvider) Reference() string {
	return "userProvider"
}

type User struct {
	ID   string
	Name string
	Age  int32
	Time time.Time
}

func (*User) JavaClassName() string {
	return "org.apache.dubbo.User"
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	var userProvider = &UserProvider{}
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&User{})
	config.Load(config.WithPath("C:\\Users\\cachen\\tencent_workspase\\dubbo-go-samples\\context\\dubbo\\go-client\\conf\\dubbogo.yml"))
	time.Sleep(3 * time.Second)

	gxlog.CInfo("\n\n\nstart to test dubbo")
	ctx := context.Background()
	ctx = context.WithValue(ctx, "name", "Alex001")
	user, err := userProvider.GetUser(ctx, &User{Name: "Alex001"})
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("response result: %v\n", user)
}
