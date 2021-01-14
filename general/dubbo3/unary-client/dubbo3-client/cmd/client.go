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
	dubbo3 "github.com/apache/dubbo-go-samples/general/dubbo3/protobuf/dubbo3"
	"github.com/apache/dubbo-go-samples/general/dubbo3/unary-client/dubbo3-client/pkg"
	"github.com/apache/dubbo-go/common/logger"
	"github.com/dubbogo/gost/log"
	"sync"
	"time"
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo3"
	_ "github.com/apache/dubbo-go/protocol/grpc"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

var grpcGreeterImpl = new(pkg.GrpcGreeterImpl)

func init() {
	config.SetConsumerService(grpcGreeterImpl)
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	config.Load()
	time.Sleep(3 * time.Second)

	gxlog.CInfo("\n\n\nstart to test dubbo")
	reply := &dubbo3.Dubbo3HelloReply{}
	req := &dubbo3.Dubbo3HelloRequest{
		Myname: "jifeng",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "tri-req-id", "test_value_XXXXXXXX")
	wg := sync.WaitGroup{}
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			err := grpcGreeterImpl.Dubbo3SayHello2(ctx, req, reply)
			if err != nil {
				logger.Error(err)
			}
			fmt.Printf("client response result: %v\n", reply)
			//wg.Done()
			wg.Done()
		}()
	}
	wg.Wait()
}
