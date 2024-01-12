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
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log"

	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
)

import (
	pb "github.com/apache/dubbo-go-samples/compatibility/api"
)

var grpcGreeterImpl = new(pb.GreeterClientImpl)

func init() {
	config.SetConsumerService(grpcGreeterImpl)
}

// export DUBBO_GO_CONFIG_PATH=$PATH_TO_SAMPLES/rpc/triple/pb/dubbogo-java/go-client/conf/dubbogo.yml
func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	gxlog.CInfo("\n\n\nstart to test dubbo")
	req := &pb.HelloRequest{
		Name: "laurence",
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, tripleConstant.TripleCtxKey("tri-req-id"), "test_value_XXXXXXXX")

	reply, err := grpcGreeterImpl.SayHello(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("client response result: %v\n", reply)
}
