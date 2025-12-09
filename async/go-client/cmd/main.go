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
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"
)

import (
	user "github.com/apache/dubbo-go-samples/async/proto"
)

var (
	userProvider   user.UserProvider
	userProviderV2 user.UserProviderV2
)

func main() {
	cli, err := client.NewClient(
		client.WithClientProtocolTriple(),
		client.WithClientURL("tri://127.0.0.1:20000"),
		client.WithClientSerialization(constant.ProtobufSerialization),
	)
	if err != nil {
		panic(err)
	}

	userProvider, err = user.NewUserProvider(cli, client.WithAsync())
	if err != nil {
		panic(err)
	}

	userProviderV2, err = user.NewUserProviderV2(cli, client.WithAsync())
	if err != nil {
		panic(err)
	}

	testAsync()
	testAsyncOneWay()
}

func testAsync() {
	req := &user.GetUserRequest{Id: "003"}
	_, err := userProvider.GetUser(context.Background(), req)
	if err != nil {
		panic(err)
	}

	logger.Info("non-blocking before async callback resp: do something ... ")

	time.Sleep(time.Second)
}

func testAsyncOneWay() {
	req := &user.SayHelloRequest{UserId: "003"}
	_, err := userProviderV2.SayHello(context.Background(), req)
	if err != nil {
		panic(err)
	}
	logger.Info("test end")
}
