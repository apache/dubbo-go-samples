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
	"strings"
	"time"

	perrors "github.com/pkg/errors"

	"github.com/dubbogo/go-zookeeper/zk"

	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/config_center/zookeeper/proto"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	_ = writeRuleToConfigCenter()

	time.Sleep(time.Second * 10)

	zkOption := config_center.WithZookeeper()
	dataIdOption := config_center.WithDataID("dubbo-go-samples-configcenter-zookeeper-client")
	addressOption := config_center.WithAddress("127.0.0.1:2181")
	groupOption := config_center.WithGroup("dubbogo")
	ins, err := dubbo.NewInstance(
		dubbo.WithConfigCenter(zkOption, dataIdOption, addressOption, groupOption),
	)
	if err != nil {
		panic(err)
	}
	// configure the params that only client layer cares
	cli, err := ins.NewClient()
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet response: %s", resp)
}

func writeRuleToConfigCenter() error {
	c, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*10)
	if err != nil {
		panic(err)
	}

	valueBytes := []byte(configCenterZKClientConfig)
	path := "/dubbo/config/dubbogo/dubbo-go-samples-configcenter-zookeeper-client"
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	paths := strings.Split(path, "/")
	for idx := 2; idx < len(paths); idx++ {
		tmpPath := strings.Join(paths[:idx], "/")
		_, err = c.Create(tmpPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			panic(err)
		}
	}

	_, err = c.Create(path, valueBytes, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		if perrors.Is(err, zk.ErrNodeExists) {
			_, stat, _ := c.Get(path)
			_, setErr := c.Set(path, valueBytes, stat.Version)
			if setErr != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return err
}

const configCenterZKClientConfig = `## set in config center, group is 'dubbogo', dataid is 'dubbo-go-samples-configcenter-zookeeper-client', namespace is default
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
