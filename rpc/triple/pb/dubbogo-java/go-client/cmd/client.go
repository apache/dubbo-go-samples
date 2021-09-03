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
	"fmt"
	"time"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo3"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/grpc"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"

	"github.com/dubbogo/gost/log"
)

import (
	pb "github.com/apache/dubbo-go-samples/rpc/triple/pb/dubbogo-java/protobuf"
)

var grpcGreeterImpl = new(pb.GreeterClientImpl)

func init() {
	config.SetConsumerService(grpcGreeterImpl)
}

// export DUBBO_GO_CONFIG_PATH=$PATH_TO_SAMPLES/rpc/triple/pb/dubbogo-java/go-client/conf/dubbogo.yml
func main() {
	config.Load()
	time.Sleep(3 * time.Second)

	gxlog.CInfo("\n\n\nstart to test dubbo")
	req := &pb.HelloRequest{
		Name: "laurence",
	}

	reply := &pb.User{}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "tri-req-id", "test_value_XXXXXXXX")

	err := grpcGreeterImpl.SayHello(ctx, req, reply)
	if err != nil {
		panic(err)
	}
	fmt.Printf("client response result: %v\n", reply)
}
