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
	greet "github.com/apache/dubbo-go-samples/mesh/proto"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("tri://server-demo.dubbo-demo.svc.cluster.local:50052"),
	)
	if err != nil {
		logger.Errorf("new client failed: %v", err)
		panic(err)
	}

	svc, err := greet.NewGreeter(cli)
	if err != nil {
		logger.Errorf("create greeter service failed: %v", err)
		panic(err)
	}

	logger.Info("start to test dubbo")
	req := &greet.HelloRequest{
		Name: "laurence",
	}
	reply, err := svc.SayHello(context.Background(), req)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("client response result: %v\n", reply)
	select {}
}
