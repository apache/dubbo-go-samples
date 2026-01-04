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
	"os"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/registry/nacos/proto"
)

func main() {
	// global conception
	// configure global configurations and common modules
	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_registry_nacos_client"),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress("127.0.0.1:8848"),
		),
	)
	if err != nil {
		logger.Errorf("new dubbo instance failed: %v", err)
		os.Exit(1)
	}
	// configure the params that only client layer cares
	cli, err := ins.NewClient()
	if err != nil {
		logger.Errorf("new client failed: %v", err)
		os.Exit(1)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		logger.Errorf("create greet service failed: %v", err)
		os.Exit(1)
	}

	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet response: %s", resp)
	select {}
}
