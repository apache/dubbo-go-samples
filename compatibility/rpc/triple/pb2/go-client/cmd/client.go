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
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
)

import (
	"github.com/apache/dubbo-go-samples/compatibility/rpc/triple/pb2/api"
	"github.com/apache/dubbo-go-samples/compatibility/rpc/triple/pb2/models"
)

var greeterProvider = new(api.GreeterClientImpl)

func init() {
	config.SetConsumerService(greeterProvider)
}

// export DUBBO_GO_CONFIG_PATH=$PATH_TO_SAMPLES/rpc/triple/pb2/go-client/conf/dubbogo.yml
func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	stream()
	unary()
}

func stream() {
	logger.Infof(">>>>> Dubbo-go client is about to call to SayHelloStream")

	ctx := context.Background()
	ctx = context.WithValue(ctx, tripleConstant.TripleCtxKey("tri-req-id"), "triple-request-id-demo")

	req := models.HelloRequest{
		Name: "laurence",
	}

	r, err := greeterProvider.SayHelloStream(ctx)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 2; i++ {
		if err := r.Send(&req); err != nil {
			logger.Errorf("Send SayHelloStream num %d request error = %v\n", i+1, err)
			return
		}
	}

	rspUser := &models.User{}
	if err := r.RecvMsg(rspUser); err != nil {
		logger.Errorf("Receive 1 SayHelloStream response user error = %v\n", err)
		return
	}
	logger.Infof("Receive 1 user = %+v\n", rspUser)
	if err := r.Send(&req); err != nil {
		logger.Errorf("Send SayHelloStream num %d request error = %v\n", 3, err)
		return
	}
	rspUser2 := &models.User{}
	if err := r.RecvMsg(rspUser2); err != nil {
		logger.Errorf("Receive 2 SayHelloStream response user error = %v\n", err)
		return
	}
	logger.Infof("Receive 2 user = %+v\n", rspUser2)
}

func unary() {
	logger.Infof(">>>>> Dubbo-go client is about to call to SayHello")

	ctx := context.Background()
	ctx = context.WithValue(ctx, tripleConstant.TripleCtxKey(tripleConstant.TripleRequestID), "triple-request-id-demo")

	req := models.HelloRequest{
		Name: "laurence",
	}
	user, err := greeterProvider.SayHello(ctx, &req)
	if err != nil {
		panic(err)
	}

	logger.Infof("Receive user = %+v\n", user)
}
