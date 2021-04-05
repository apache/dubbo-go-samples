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
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	"github.com/apache/dubbo-go/common/logger"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo3"
	_ "github.com/apache/dubbo-go/protocol/grpc"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

import (
	dubbo3pb "github.com/apache/dubbo-go-samples/general/dubbo3/pb/protobuf/dubbo3"
	"github.com/apache/dubbo-go-samples/general/dubbo3/pb/stream-client/dubbogo-client/pkg"
)

var greeterProvider = new(pkg.GreeterProvider)

func init() {
	config.SetConsumerService(greeterProvider)
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	config.Load()
	time.Sleep(time.Second * 3)

	testSayHello()
}

func testSayHello() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "tri-req-id", "triple-request-id-demo")

	req := dubbo3pb.HelloRequest{
		Name: "laurence",
	}

	r, err := greeterProvider.SayHelloStream(ctx)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		if err := r.Send(&req); err != nil {
			logger.Errorf("Send SayHelloStream num %d request error = %v\n", i+1, err)
			return
		}
	}

	rspUser := &dubbo3pb.User{}
	if err := r.RecvMsg(rspUser); err != nil {
		logger.Errorf("Receive 1 SayHelloStream response user error = %v\n", err)
		return
	}
	logger.Infof("Receive 1 user = %+v\n", rspUser)
	rspUser2 := &dubbo3pb.User{}
	if err := r.RecvMsg(rspUser2); err != nil {
		logger.Errorf("Receive 2 SayHelloStream response user error = %v\n", err)
		return
	}
	logger.Infof("Receive 2 user = %+v\n", rspUser2)
}
