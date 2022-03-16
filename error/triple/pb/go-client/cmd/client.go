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

	tripleCommon "github.com/dubbogo/triple/pkg/common"
)

import (
	triplepb "github.com/apache/dubbo-go-samples/api"
)

var greeterProvider = new(triplepb.GreeterClientImpl)

func init() {
	config.SetConsumerService(greeterProvider)
}

// export DUBBO_GO_CONFIG_PATH=$PATH_TO_SAMPLES/error/triple/pb/go-client/conf/dubbogo.yml
func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	req := triplepb.HelloRequest{
		Name: "laurence",
	}

	if user, err := greeterProvider.SayHello(context.TODO(), &req); err != nil {
		logger.Infof("response result: %v, error = %s", user, err)
		logger.Infof("error details = %+v", err.(tripleCommon.TripleError).Stacks())
		logger.Infof("error code = %+v", err.(tripleCommon.TripleError).Code())
		logger.Infof("error message = %+v", err.(tripleCommon.TripleError).Message())
	}
}
