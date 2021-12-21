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

	hessian "github.com/apache/dubbo-go-hessian2"
)

import (
	"github.com/apache/dubbo-go-samples/rpc/dubbo/go-client/pkg"
)

var (
	userProvider = &pkg.UserProvider{}
)

// need to setup environment variable "DUBBO_GO_CONFIG_PATH" to "conf/dubbogo.yml" before run
func main() {
	hessian.RegisterJavaEnum(pkg.Gender(pkg.MAN))
	hessian.RegisterJavaEnum(pkg.Gender(pkg.WOMAN))
	hessian.RegisterPOJO(&pkg.User{})

	config.SetConsumerService(userProvider)

	err := config.Load()
	if err != nil {
		panic(err)
	}

	logger.Infof("\n\ntest")
	test()
}

func test() {
	logger.Infof("\n\n\nstart to test dubbo")
	reqUser := &pkg.User{}
	reqUser.ID = "003"
	user, err := userProvider.GetUser(context.TODO(), reqUser)
	if err != nil {
		panic(err)
	}
	logger.Infof("response result: %v", user)

	logger.Infof("\n\n\nstart to test dubbo - enum")
	gender, err := userProvider.GetGender(context.TODO(), 1)
	if err != nil {
		panic(err)
	}
	logger.Infof("response result: %v", gender)

	logger.Infof("\n\n\nstart to test dubbo - GetUser0")
	ret, err := userProvider.GetUser0("003", "Moorse")
	if err != nil {
		panic(err)
	}
	logger.Infof("response result: %v", ret)

	logger.Infof("\n\n\nstart to test dubbo - GetUsers")
	ret1, err := userProvider.GetUsers([]string{"002", "003"})
	if err != nil {
		panic(err)
	}
	logger.Infof("response result: %v", ret1)

	logger.Infof("\n\n\nstart to test dubbo - getUser")

	var i int32 = 1
	user, err = userProvider.GetUser2(context.TODO(), i)
	if err != nil {
		panic(err)
	}
	logger.Infof("response result: %v", user)

	logger.Infof("\n\n\nstart to test dubbo - getErr")
	reqUser.ID = "003"
	_, err = userProvider.GetErr(context.TODO(), reqUser)
	if err == nil {
		panic("err is nil")
	}
	logger.Infof("getErr - error: %v", err)
}
