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
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/direct/proto"
)

const (
	RegistryAddress = "127.0.0.1:8848"
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("condition-client"),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress(RegistryAddress),
		),
		dubbo.WithConfigCenter( // configure config center to enable condition router
			config_center.WithNacos(),
			config_center.WithAddress(RegistryAddress),
		),
	)

	if err != nil {
		logger.Errorf("new instance failed: %v", err)
		panic(err)
	}

	cli, err := ins.NewClient()

	if err != nil {
		logger.Errorf("new client failed: %v", err)
		panic(err)
	}

	srv, err := greet.NewGreetService(cli)

	if err != nil {
		logger.Errorf("new service failed: %v", err)
		panic(err)
	}

	for {
		time.Sleep(5 * time.Second) // sleep 5 seconds
		rep, err := srv.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
		printRes(rep, err)
	}

}

func printRes(rep *greet.GreetResponse, err error) {
	if err != nil {
		logger.Errorf("call greet method failed: %v", err)
	} else {
		logger.Infof("receive: %s", rep.GetGreeting())
	}
}
