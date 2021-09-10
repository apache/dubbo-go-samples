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
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	"github.com/apache/dubbo-go-samples/api"
)

const configCenterZKClientConfig = `## set in config center, group is 'dubbogo', dataid is 'dubbo-go-samples-configcenter-zookeeper-client', namespace is default
dubbo:
  registries:
    demoZK:
      protocol: nacos
      timeout: 3s
      address: 127.0.0.1:8848
  consumer:
    registry:
      - demoZK
    references:
      GreeterClientImpl:
        protocol: tri
        interface: com.apache.dubbo.sample.basic.IGreeter # must be compatible with grpc or dubbo-java`

var grpcGreeterImpl = new(api.GreeterClientImpl)

// There is no need to export DUBBO_GO_CONFIG_PATH, as you are using config api to set config
func main() {
	dynamicConfig, err := config.NewConfigCenterConfig(
		config.WithConfigCenterProtocol("zookeeper"),
		config.WithConfigCenterAddress("127.0.0.1:2181")).GetDynamicConfiguration()
	if err != nil {
		panic(err)
	}
	if err := dynamicConfig.PublishConfig("dubbo-go-samples-configcenter-zookeeper-client", "dubbogo", configCenterZKClientConfig); err != nil {
		panic(err)
	}

	config.SetConsumerService(grpcGreeterImpl)

	centerConfig := config.NewConfigCenterConfig(
		config.WithConfigCenterProtocol("zookeeper"),
		config.WithConfigCenterAddress("localhost:2181"),
		config.WithConfigCenterDataID("dubbo-go-samples-configcenter-zookeeper-client"),
		config.WithConfigCenterGroup("dubbogo"),
	)

	rootConfig := config.NewRootConfig(
		config.WithRootCenterConfig(centerConfig),
	)

	if err := rootConfig.Init(); err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)

	logger.Info("start to test dubbo")
	req := &api.HelloRequest{
		Name: "laurence",
	}
	reply, err := grpcGreeterImpl.SayHello(context.Background(), req)
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("client response result: %v\n", reply)
}
