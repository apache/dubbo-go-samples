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
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	"github.com/apache/dubbo-go-samples/rpc/jsonrpc/go-client/pkg"
)

var (
	// nolint
	survivalTimeout int = 10e9
	userProvider        = &pkg.UserProvider{}
)

func init() {
	config.SetConsumerService(userProvider)
}

// Do some checking before the system starts up:
// 1. env config
// 		`export DUBBO_GO_CONFIG_PATH= ROOT_PATH/conf/dubbogo.yml` or `dubbogo.yaml`
func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	logger.Infof("\n\ntest")
	test()
}

func test() {
	logger.Infof("\n\n\nstart to test jsonrpc")
	user, err := userProvider.GetUser(context.TODO(), "A003")
	if err != nil {
		panic(err)
	}
	logger.Infof("response result: %v", user)

	logger.Infof("\n\n\nstart to test jsonrpc - getUser")
	rep2, err := userProvider.GetUser2(context.TODO(), "A001")
	if err != nil {
		panic(err)
	}
	logger.Infof("response result: %v", rep2)
}
