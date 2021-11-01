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
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"
)

import (
	"github.com/apache/dubbo-go-samples/filter/tpslimit/go-client/pkg"
)

var userProvider = &pkg.UserProvider{}

func init() {
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&pkg.User{})
}

func main() {
	err := config.Load()
	if err != nil {
		panic(err)
	}

	var successCount, failCount int64
	logger.Infof("\n\n\nstart to test dubbo")
	for i := 0; i < 60; i++ {
		time.Sleep(200 * time.Millisecond)
		user, err := userProvider.GetUser(context.TODO(), "A001")
		if err != nil {
			failCount++
			logger.Infof("error: %v\n", err)
		} else {
			successCount++
		}
		logger.Infof("response: %v\n", user)
	}
	logger.Infof("failCount=%v, failCount=%v\n", successCount, failCount)
}
