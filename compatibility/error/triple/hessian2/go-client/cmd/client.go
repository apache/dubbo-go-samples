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

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log/logger"

	tripleCommon "github.com/dubbogo/triple/pkg/common"
)

var errorResponseProvider = new(ErrorResponseProvider)

func init() {
	config.SetConsumerService(errorResponseProvider)
	hessian.RegisterPOJO(&User{})
}

// need to setup environment variable "DUBBO_GO_CONFIG_PATH" to "conf/dubbogo.yml" before run
func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}
	testErrorService()
	testService()
}

func testErrorService() {
	if user, err := errorResponseProvider.GetUser(context.TODO(), &User{Name: "laurence"}); err != nil {
		logger.Infof("response result: %v, error = %s", user, err)
		logger.Infof("error details = %+v", err.(tripleCommon.TripleError).Stacks())
	}
}

func testService() {
	if user, err := errorResponseProvider.GetUserWithoutError(context.TODO(), &User{Name: "laurence"}); err != nil {
		logger.Infof("response result: %v, error = %s", user, err)
		logger.Infof("error details = %+v", err.(tripleCommon.TripleError).Stacks())
	}
}
