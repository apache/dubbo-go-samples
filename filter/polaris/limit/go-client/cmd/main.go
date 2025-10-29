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
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/helloworld/proto"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli,
		client.WithInterface("org.apache.dubbo.UserProvider.Test"),
	)
	if err != nil {
		panic(err)
	}

	var (
		successCount int
		failCount    int
	)

	for i := 0; i < 10; i++ {
		time.Sleep(50 * time.Millisecond)
		user, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
		if err != nil {
			failCount++
			continue
		}

		logger.Info(user)
		successCount++
	}
	logger.Infof("successCount=%v, failCount=%v\n", successCount, failCount)

	if !(successCount == 1 && failCount == 9) {
		panic("ratelimit expect 1 success and 9 fail")
	}

	successCount = 0
	failCount = 0

	svc, err = greet.NewGreetService(cli,
		client.WithInterface("org.apache.dubbo.UserProvider.Test2"),
		client.WithGroup("myInterfaceGroup"),
		client.WithVersion("myInterfaceVersion"),
	)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		time.Sleep(50 * time.Millisecond)
		user, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
		if err != nil {
			failCount++
			continue
		}

		logger.Info(user)
		successCount++
	}
	logger.Infof("successCount=%v, failCount=%v\n", successCount, failCount)

	if !(successCount == 10 && failCount == 0) {
		panic("ratelimit expect 10 success and 0 fail")
	}
}
