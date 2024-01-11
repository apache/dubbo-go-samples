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
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-samples/compatibility/api"
)

const configCenterZKClientConfig = `## set in config center, group is 'dubbogo', dataid is 'dubbo-go-samples-configcenter-zookeeper-client', namespace is default
dubbo:
  registries:
    demoZK:
      protocol: nacos
      address: 127.0.0.1:8848
  consumer:
    references:
      GreeterClientImpl:
        protocol: tri
`

var grpcGreeterImpl = new(api.GreeterClientImpl)

// There is no need to export DUBBO_GO_CONFIG_PATH, as you are using config api to set config
func main() {
	dynamicConfig, err := config.NewConfigCenterConfigBuilder().
		SetProtocol("zookeeper").
		SetAddress("127.0.0.1:2181").
		Build().GetDynamicConfiguration()
	if err != nil {
		panic(err)
	}

	if err = dynamicConfig.PublishConfig("dubbo-go-samples-configcenter-zookeeper-client", "dubbogo", configCenterZKClientConfig); err != nil {
		panic(err)
	}

	config.SetConsumerService(grpcGreeterImpl)

	rootConfig := config.NewRootConfigBuilder().
		SetConfigCenter(config.NewConfigCenterConfigBuilder().
			SetProtocol("nacos").SetAddress("127.0.0.1:2182").
			SetDataID("dubbo-go-samples-configcenter-zookeeper-client").
			Build()).
		Build()

	if err = config.Load(config.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}
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
