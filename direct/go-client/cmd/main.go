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
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/direct/proto"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
	)
	if err != nil {
		panic(err) // fail fast if client cannot be created
	}

	greetService, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err) // fail fast if service proxy cannot be created
	}

	// Direct call: 1 request -> 1 response
	name := "Golang Client dubbo-go"
	req := &greet.GreetRequest{Name: name}
	resp, err := greetService.Greet(context.Background(), req)
	if err != nil {
		panic(err)
	}
	if resp == nil {
		panic("direct call failed: empty response")
	}

	// Go 服务返回 "hello {name}"，Java 示例服务返回 "hello from java server, {name}"
	expectGo := "hello " + name
	expectJava := "hello from java server, " + name
	if resp.Greeting != expectGo && resp.Greeting != expectJava {
		panic("unexpected greeting: " + resp.Greeting)
	}

	logger.Infof("direct call response: %s", resp.Greeting)
}
