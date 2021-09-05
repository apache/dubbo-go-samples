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
)

import (
	triplepb "github.com/apache/dubbo-go-samples/rpc/triple/pb/dubbogo-grpc/protobuf/triple"
)

var greeterProvider = new(triplepb.GreeterClientImpl)

func init() {
	config.SetConsumerService(greeterProvider)
}

// export DUBBO_GO_CONFIG_PATH=$PATH_TO_SAMPLES/rpc/triple/pb/dubbogo-grpc/stream-client/dubbogo-client/conf/dubbogo.yml
func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 3)

	testSayHello()
}

func testSayHello() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "tri-req-id", "triple-request-id-demo")

	req := triplepb.HelloRequest{
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

	rspUser := &triplepb.User{}
	if err := r.RecvMsg(rspUser); err != nil {
		logger.Errorf("Receive 1 SayHelloStream response user error = %v\n", err)
		return
	}
	logger.Infof("Receive 1 user = %+v\n", rspUser)
	if err := r.Send(&req); err != nil {
		logger.Errorf("Send SayHelloStream num %d request error = %v\n", 3, err)
		return
	}
	rspUser2 := &triplepb.User{}
	if err := r.RecvMsg(rspUser2); err != nil {
		logger.Errorf("Receive 2 SayHelloStream response user error = %v\n", err)
		return
	}
	logger.Infof("Receive 2 user = %+v\n", rspUser2)
}
