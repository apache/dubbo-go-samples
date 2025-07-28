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
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/go-zookeeper/zk"

	"github.com/dubbogo/gost/log/logger"

	perrors "github.com/pkg/errors"
)

import (
	greet "github.com/apache/dubbo-go-samples/config_center/zookeeper/proto"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	return resp, nil
}

func main() {
	_ = writeRuleToConfigCenter()
	time.Sleep(time.Second * 10)

	ins, err := dubbo.NewInstance(
		dubbo.WithConfigCenter(
			config_center.WithZookeeper(),
			config_center.WithDataID("dubbo-go-samples-configcenter-zookeeper-server"),
			config_center.WithAddress("127.0.0.1:2181"),
			config_center.WithGroup("dubbogo"),
		),
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
      port: 50000
  provider:
    services:
      GreeterProvider:
        interface: com.apache.dubbo.sample.basic.IGreeter
`

func ensurePath(c *zk.Conn, path string, data []byte, flags int32, acl []zk.ACL) error {
	_, err := c.Create(path, data, flags, acl)
	return err
}

func writeRuleToConfigCenter() error {
	c, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*10)
	if err != nil {
		panic(err)
	}

	valueBytes := []byte(configCenterZKServerConfig)
	path := "/dubbo/config/dubbogo/dubbo-go-samples-configcenter-zookeeper-server"
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
