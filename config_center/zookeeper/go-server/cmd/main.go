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
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/config_center/zookeeper/proto"
	"github.com/dubbogo/gost/log/logger"
	"time"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	return resp, nil
}

const configCenterZKServerConfig = `## set in config center, group is 'dubbogo', dataid is 'dubbo-go-samples-configcenter-zookeeper-server', namespace is default
dubbo:
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
        interface: com.apache.dubbo.sample.basic.IGreeter
`

func main() {
	dynamicConfig, err := config.NewConfigCenterConfigBuilder().
		SetProtocol("zookeeper").
		SetAddress("127.0.0.1:2181").
		Build().GetDynamicConfiguration()
	if err != nil {
		panic(err)
	}

	if err = dynamicConfig.PublishConfig("dubbo-go-samples-configcenter-zookeeper-server", "dubbogo", configCenterZKServerConfig); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 10)

	zkOption := config_center.WithZookeeper()
	dataIdOption := config_center.WithDataID("dubbo-go-samples-configcenter-zookeeper-server")
	addressOption := config_center.WithAddress("127.0.0.1:2181")
	groupOption := config_center.WithGroup("dubbogo")
	ins, err := dubbo.NewInstance(
		dubbo.WithConfigCenter(zkOption, dataIdOption, addressOption, groupOption),
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
