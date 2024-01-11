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
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"
)

import (
	_ "github.com/apache/dubbo-go-samples/compatibility/tls/triple/codec"
)

type User struct {
	ID   string
	Name string
	Age  int32
}

type UserProvider struct {
	GetUser func(context.Context, *User, *User, string) (*User, error)
}

var userProvider = new(UserProvider)

// export DUBBO_GO_CONFIG_PATH=PATH_TO_SAMPLES/rpc/triple/codec-extension/go-client/conf/dubbogo.yml
func main() {
	config.SetConsumerService(userProvider)
	if err := config.Load(); err != nil {
		panic(err)
	}

	user, err := userProvider.GetUser(context.TODO(), &User{Name: "zlber"}, &User{Name: "zlber2"}, "testName")
	if err != nil {
		panic(err)
	}
	logger.Infof("response result: %v\n", user)
}
