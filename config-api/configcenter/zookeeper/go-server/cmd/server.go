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

const configCenterZKServerConfig = `# set in config center, group is 'dubbogo', dataid is 'dubbo-go-samples-configcenter-zookeeper-server', namespace is default
dubbo:
  registries:
    demoZK:
      protocol: nacos
      timeout: 3s
      address: 127.0.0.1:8848
  protocols:
    triple:
      name: tri
      port: 20000
  provider:
    registry-ids:
      - demoZK
    services:
      GreeterProvider:
        protocol-ids: triple
        interface: com.apache.dubbo.sample.basic.IGreeter # must be compatible with grpc or dubbo-java`

type GreeterProvider struct {
	api.UnimplementedGreeterServer
}

func (s *GreeterProvider) SayHello(ctx context.Context, in *api.HelloRequest) (*api.User, error) {
	logger.Infof("Dubbo3 GreeterProvider get user name = %s\n", in.Name)
	return &api.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, nil
}

// There is no need to export DUBBO_GO_CONFIG_PATH, as you are using config api to set config
func main() {
	dynamicConfig, err := config.NewConfigCenterConfigBuilder().
		SetProtocol("zookeeper").
		SetAddress("127.0.0.1:2181").
		Build().GetDynamicConfiguration()
	if err != nil {
		panic(err)
	}
	if err := dynamicConfig.PublishConfig("dubbo-go-samples-configcenter-zookeeper-server", "dubbogo", configCenterZKServerConfig); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 10)

	config.SetProviderService(&GreeterProvider{})

	rootConfig := config.NewRootConfigBuilder().
		SetConfigCenter(config.NewConfigCenterConfigBuilder().
			SetProtocol("zookeeper").SetAddress("127.0.0.1:2181").
			SetDataID("dubbo-go-samples-configcenter-zookeeper-server").
			SetGroup("dubbogo").
			Build()).
		Build()

	if err := rootConfig.Init(); err != nil {
		panic(err)
	}
	select {}
}
