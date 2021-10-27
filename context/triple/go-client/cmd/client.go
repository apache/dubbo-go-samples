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

	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
)

import (
	"github.com/apache/dubbo-go-samples/api"
)

var grpcGreeterImpl = new(api.GreeterClientImpl)

func init() {
	config.SetConsumerService(grpcGreeterImpl)
}

// export DUBBO_GO_CONFIG_PATH= PATH_TO_SAMPLES/helloworld/go-client/conf/dubbogo.yml
func main() {
	err := config.Load()
	if err != nil {
		panic(err)
	}

	logger.Info("start to test dubbo")
	req := &api.HelloRequest{
		Name: "laurence",
	}
	ctx := context.Background()
	// set user defined context attachment
	ctx = context.WithValue(ctx, tripleConstant.CtxAttachmentKey, "user-defined-value")
	reply, err := grpcGreeterImpl.SayHello(ctx, req)
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("client response result: %v\n", reply)
}
