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
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo3"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/grpc"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
)

import (
	dubbo3pb "github.com/apache/dubbo-go-samples/general/dubbo3/pb/dubbogo-grpc/protobuf/dubbo3"
)

var greeterProvider = new(dubbo3pb.GreeterClientImpl)

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
	ctx = context.WithValue(ctx, tripleConstant.TripleCtxKey(tripleConstant.TripleRequestID), "triple-request-id-demo")

	req := dubbo3pb.HelloRequest{
		Name: "laurence",
	}
	user := dubbo3pb.User{}
	err := greeterProvider.SayHello(ctx, &req, &user)
	if err != nil {
		panic(err)
	}

	logger.Infof("Receive user = %+v\n", user)
}
