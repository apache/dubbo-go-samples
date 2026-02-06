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
	"fmt"
	"strings"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

import (
	greet "github.com/apache/dubbo-go-samples/config_center/nacos/proto"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	logger.Info("Received request:" + req.Name)
	resp := &greet.GreetResponse{Greeting: "Hello, this is dubbo go server!" + " I received: " + req.Name}
	return resp, nil
}

const configCenterNacosServerConfig = `## set in config center, group is 'dubbo', dataid is 'dubbo-go-samples-configcenter-nacos-server', namespace is default
dubbo:
  application:
    name: dubbo-go-provider
  registries:
    demoZK:
      protocol: zookeeper
      timeout: 3s
      address: '127.0.0.1:2181'
  protocols:
    triple:
      name: tri
      port: 20000
  provider:
    services:
      GreeterProvider:
        interface: greet.GreetService
`

func main() {
	clientConfig := constant.ClientConfig{}
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			"127.0.0.1",
			8848,
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		),
	}
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}

	if err := publishAndWaitConfig(
		configClient,
		"dubbo-go-samples-configcenter-nacos-go-server",
		"dubbo",
		configCenterNacosServerConfig,
		10*time.Second,
		200*time.Millisecond,
	); err != nil {
		panic(err)
	}

	nacosOption := config_center.WithNacos()
	dataIdOption := config_center.WithDataID("dubbo-go-samples-configcenter-nacos-go-server")
	addressOption := config_center.WithAddress("127.0.0.1:8848")
	groupOption := config_center.WithGroup("dubbo")
	ins, err := dubbo.NewInstance(
		dubbo.WithConfigCenter(nacosOption, dataIdOption, addressOption, groupOption),
	)
	if err != nil {
		panic(err)
	}
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}

	if err = greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
		panic(err)
	}

	if err = srv.Serve(); err != nil {
		logger.Error(err)
	}
}

func publishAndWaitConfig(
	configClient config_client.IConfigClient,
	dataID string,
	group string,
	content string,
	timeout time.Duration,
	pollInterval time.Duration,
) error {
	success, err := configClient.PublishConfig(vo.ConfigParam{
		DataId:  dataID,
		Group:   group,
		Content: content,
	})
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("publish config failed")
	}

	deadline := time.Now().Add(timeout)
	for {
		current, err := configClient.GetConfig(vo.ConfigParam{
			DataId: dataID,
			Group:  group,
		})
		if err == nil && strings.TrimSpace(current) == strings.TrimSpace(content) {
			return nil
		}
		if time.Now().After(deadline) {
			if err != nil {
				return err
			}
			return fmt.Errorf("wait for config center timeout")
		}
		time.Sleep(pollInterval)
	}
}
