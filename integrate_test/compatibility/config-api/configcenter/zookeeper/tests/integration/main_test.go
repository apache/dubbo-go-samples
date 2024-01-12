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
	"testing"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	dubbo3pb "github.com/apache/dubbo-go-samples/compatibility/api"
)

const configCenterZKTestClientConfig = `## set in config center, group is 'dubbogo', dataid is 'dubbo-go-samples-configcenter-zookeeper-client', namespace is default
dubbo:
  registries:
    demoZK:
      protocol: zookeeper
      address: 127.0.0.1:2181
  consumer:
    references:
      GreeterClientImpl:
        protocol: tri
`

var greeterProvider = new(dubbo3pb.GreeterClientImpl)

func TestMain(m *testing.M) {
	dynamicConfig, err := config.NewConfigCenterConfigBuilder().
		SetProtocol("zookeeper").
		SetAddress("127.0.0.1:2181").
		Build().GetDynamicConfiguration()
	if err != nil {
		panic(err)
	}

	if err := dynamicConfig.PublishConfig("dubbo-go-samples-configcenter-zookeeper-client", "dubbogo", configCenterZKTestClientConfig); err != nil {
		panic(err)
	}

	config.SetConsumerService(greeterProvider)
	time.Sleep(time.Second * 15)

	rootConfig := config.NewRootConfigBuilder().
		SetConfigCenter(config.NewConfigCenterConfigBuilder().
			SetProtocol("zookeeper").SetAddress("127.0.0.1:2181").
			SetDataID("dubbo-go-samples-configcenter-zookeeper-client").
			SetGroup("dubbogo").
			Build()).
		Build()

	if err := config.Load(config.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)

	os.Exit(m.Run())
}
