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

package integration

import (
	"os"
	"strings"
	"testing"
	"time"

	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/config_center"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	greet "github.com/apache/dubbo-go-samples/config_center/nacos/proto"
)

const configCenterNacosClientConfig = `## set in config center, group is 'dubbo', dataid is 'dubbo-go-samples-configcenter-nacos-client', namespace is default
dubbo:
  registries:
    demoZK:
      protocol: zookeeper
      timeout: 3s
      address: 127.0.0.1:2181
  consumer:
    references:
      GreeterClientImpl:
        protocol: tri
        interface: greet.GreetService 
`

var greeterProvider greet.GreetService

func TestMain(m *testing.M) {
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

	success, err := configClient.PublishConfig(vo.ConfigParam{
		DataId:  "dubbo-go-samples-configcenter-nacos-client",
		Group:   "dubbo",
		Content: configCenterNacosClientConfig,
	})
	if err != nil {
		panic(err)
	}
	if !success {
		return
	}

	deadline := time.Now().Add(10 * time.Second)
	for {
		content, err := configClient.GetConfig(vo.ConfigParam{
			DataId: "dubbo-go-samples-configcenter-nacos-client",
			Group:  "dubbo",
		})
		if err == nil && strings.TrimSpace(content) == strings.TrimSpace(configCenterNacosClientConfig) {
			break
		}
		if time.Now().After(deadline) {
			if err != nil {
				panic(err)
			}
			panic("wait for config center timeout")
		}
		time.Sleep(200 * time.Millisecond)
	}

	nacosOption := config_center.WithNacos()
	dataIdOption := config_center.WithDataID("dubbo-go-samples-configcenter-nacos-client")
	addressOption := config_center.WithAddress("127.0.0.1:8848")
	groupOption := config_center.WithGroup("dubbo")
	ins, err := dubbo.NewInstance(
		dubbo.WithConfigCenter(nacosOption, dataIdOption, addressOption, groupOption),
	)
	if err != nil {
		panic(err)
	}
	// configure the params that only client layer cares
	cli, err := ins.NewClient()
	if err != nil {
		panic(err)
	}

	greeterProvider, err = greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
