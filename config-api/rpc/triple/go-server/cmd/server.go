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
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	"github.com/apache/dubbo-go-samples/api"
)

type GreeterProvider struct {
	api.GreeterProviderBase
}

func (s *GreeterProvider) SayHello(ctx context.Context, in *api.HelloRequest) (*api.User, error) {
	logger.Infof("Dubbo3 GreeterProvider get user name = %s\n", in.Name)
	return &api.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, nil
}

// There is no need to export DUBBO_GO_CONFIG_PATH, as you are using config api to set config
func main() {
	config.SetProviderService(&GreeterProvider{})

	rootConfig := config.NewRootConfigBuilder().
		SetProvider(config.NewProviderConfigBuilder().
			SetRegistries("zk").
			AddService("GreeterProvider", config.NewServiceConfigBuilder().
				SetInterface("com.apache.dubbo.sample.basic.IGreeter").
				SetProtocols("tripleKey").
				Build()).
			Build()).
		AddRegistry("zk", config.NewRegistryConfigWithProtocolDefaultPort("zookeeper")).
		AddProtocol("tripleKey", config.NewProtocolConfigBuilder().
			SetName("tri").
			SetPort("20000").
			Build()).
		Build()

	if err := rootConfig.Init(); err != nil {
		panic(err)
	}
	select {}
}
