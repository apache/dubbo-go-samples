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
	"github.com/apache/dubbo-go-samples/apisix/helloworld/protobuf/helloworld"
)

type GreeterProvider struct {
	helloworld.UnimplementedGreeterServer
}

func main() {
	config.SetProviderService(&GreeterProvider{})

	nacosConfig := config.NewRegistryConfigWithProtocolDefaultPort("nacos")
	nacosConfig.Address = "172.19.0.3:8848"
	rc := config.NewRootConfigBuilder().
		SetProvider(config.NewProviderConfigBuilder().
			AddService("GreeterProvider", config.NewServiceConfigBuilder().Build()).
			Build()).
		AddProtocol("tripleProtocolKey", config.NewProtocolConfigBuilder().
			SetName("tri").
			SetPort("20001").
			Build()).
		AddRegistry("registryKey", nacosConfig).
		Build()

	// start dubbo-go framework with configuration
	if err := config.Load(config.WithRootConfig(rc)); err != nil {
		logger.Infof("init ERR = %s\n", err.Error())
	}

	select {}
}

func (s *GreeterProvider) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.User, error) {
	logger.Infof("SayHello in %s", in.String())
	helloworld := &helloworld.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}
	logger.Infof("SayHello out %s", helloworld.String())
	return helloworld, nil
}
