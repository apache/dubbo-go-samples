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
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/config_center/apollo/proto"
)

// Apollo Configuration Center Parameters
const (
	apolloMetaAddress = "tony2c4g:8080"
	apolloAppID       = "SampleApp"
	apolloCluster     = "default"
	apolloNamespace   = "dubbo.yml"
)

const configCenterApolloClientConfig = `## set in config center, namespace is 'dubbo.yml', appId is 'SampleApp'
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
        interface: com.apache.dubbo.sample.basic.IGreeter 
`

func main() {
	// Initialize client using configuration center
	ins, err := dubbo.NewInstance(
		dubbo.WithConfigCenter(
			config_center.WithApollo(),
			config_center.WithAddress(apolloMetaAddress),
			config_center.WithNamespace(apolloNamespace),
			config_center.WithDataID(apolloNamespace),
			config_center.WithAppID(apolloAppID),
			config_center.WithCluster(apolloCluster),
			//config_center.WithFileExtProperties(),
		),
	)
	if err != nil {
		panic(err)
	}

	// Configure client parameters
	cli, err := ins.NewClient()
	if err != nil {
		panic(err)
	}

	// Create service client
	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	// Call remote service
	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "Apollo"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet response: %s", resp)
}
