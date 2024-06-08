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
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	// 	"github.com/apache/dubbo-go-hessian2/java8_time"
	// 	"github.com/apache/dubbo-go-hessian2/java_exception"
	// 	"github.com/apache/dubbo-go-hessian2/java_sql_time"
	// 	"github.com/apache/dubbo-go-hessian2/java_util"
	greet "github.com/apache/dubbo-go-samples/java_interop/non-protobuf-triple/proto"
	"github.com/dubbogo/gost/log/logger"
	//	java_math "github.com/dubbogo/gost/math/big"
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_interop_non_idl_triple_client"),
	)
	if err != nil {
		panic(err)
	}

	cli, err := ins.NewClient(
		client.WithClientProtocolTriple(),
		client.WithClientSerialization(constant.Hessian2Serialization),
		client.WithClientURL("127.0.0.1:50052"),
	)

	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetingsService(cli)
	if err != nil {
		panic(err)
	}

	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{
		Name: "dubbo-go",
		// 		Way:              greet.GreetEnum_GREET_ENUM_1,
		// 		Time:             &java_sql_time.Time{},
		// 		Duration:         &java8_time.Duration{},
		// 		RuntimeException: &java_exception.RuntimeException{},
		// 		Uuid:             &java_util.UUID{},
		// 		BigInteger:       &java_math.Integer{},
	})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet response: %s", resp.Greeting)
}
