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
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-samples/compatibility/async/go-client/pkg"
)

var (
	userProvider   = &pkg.UserProvider{}
	userProviderV2 = &pkg.UserProviderV2{}
)

// need to setup environment variable "DUBBO_GO_CONFIG_PATH" to "conf/dubbogo.yml" before run
func main() {
	hessian.RegisterJavaEnum(pkg.MAN)
	hessian.RegisterJavaEnum(pkg.WOMAN)
	hessian.RegisterPOJO(&pkg.User{})

	config.SetConsumerService(userProvider)
	config.SetConsumerService(userProviderV2)

	err := config.Load()
	if err != nil {
		panic(err)
	}

	testAsync()
	testAsyncOneWay()
}

// two-way
func testAsync() {
	reqUser := &pkg.User{}
	reqUser.ID = "003"
	_, err := userProvider.GetUser(context.TODO(), reqUser)
	if err != nil {
		panic(err)
	}

	// Mock do something else
	logger.Info("non-blocking before async callback resp: do something ... ")

	// Wait for Callback
	time.Sleep(time.Second)
}

// one-way
func testAsyncOneWay() {
	err := userProviderV2.SayHello(context.TODO(), "003")
	if err != nil {
		panic(err)
	}
	logger.Info("test end")
}
