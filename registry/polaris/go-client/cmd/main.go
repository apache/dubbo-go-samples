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
	"github.com/dubbogo/gost/log/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"
)

type UserProviderWithCustomGroupAndVersion struct {
	GetUser func(ctx context.Context, req *User) (rsp *User, err error)
}

type UserProvider struct {
	GetUser func(ctx context.Context, req *User) (rsp *User, err error)
}

type User struct {
	ID   string
	Name string
	Age  int32
	Time time.Time
}

func (u *User) JavaClassName() string {
	return "org.apache.dubbo.User"
}

func main() {
	var userProvider = &UserProvider{}
	var userProviderWithCustomRegistryGroupAndVersion = &UserProviderWithCustomGroupAndVersion{}
	config.SetConsumerService(userProvider)
	config.SetConsumerService(userProviderWithCustomRegistryGroupAndVersion)
	hessian.RegisterPOJO(&User{})
	err := config.Load()
	if err != nil {
		panic(err)
	}

	logger.Infof("\n\n\nstart to test dubbo")
	user, err := userProvider.GetUser(context.TODO(), &User{Name: "Alex001"})
	if err != nil {
		logger.Errorf("error: %v\n", err)
		os.Exit(1)
		return
	}
	logger.Infof("response result: %v\n", user)

	user, err = userProviderWithCustomRegistryGroupAndVersion.GetUser(context.TODO(), &User{Name: "Alex001"})
	if err != nil {
		logger.Errorf("error: %v\n", err)
		os.Exit(1)
		return
	}
	logger.Infof("response result: %v\n", user)
}
