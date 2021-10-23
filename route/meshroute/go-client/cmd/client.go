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
	"io/ioutil"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	"github.com/apache/dubbo-go-samples/api"
)

var grpcGreeterImpl = api.GreeterClientImpl{}

// export DUBBO_GO_CONFIG_PATH= PATH_TO_SAMPLES/helloworld/go-client/conf/dubbogo.yml
func main() {
	config.SetConsumerService(grpcGreeterImpl)
	if err := config.Load(); err != nil {
		panic(err)
	}

	routeBytes, err := ioutil.ReadFile("../conf/mesh_route.yml")
	if err != nil {
		panic(err)
	}

	dynamicConfiguration, err := config.GetRootConfig().ConfigCenter.GetDynamicConfiguration()
	if err != nil {
		panic(err)
	}

	// 1. publish mesh route config
	err = dynamicConfiguration.PublishConfig("dubbo.io.MESHAPPRULE", "dubbo", string(routeBytes))
	if err != nil {
		return
	}

	logger.Info("start to test dubbo")
	req := &api.HelloRequest{
		Name: "Dong",
	}
	reply, err := grpcGreeterImpl.SayHello(context.Background(), req)
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("client response result: %v\n", reply)
	select {}
}
