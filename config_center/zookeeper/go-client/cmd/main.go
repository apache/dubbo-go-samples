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

func main() {
	// write configuration to config center
	if err := writeRuleToConfigCenter(); err != nil {
		logger.Errorf("Failed to write config to config center: %v", err)
		panic(err)
	}
	logger.Info("Successfully wrote config to ZooKeeper")

	// wait for config write to finish
	time.Sleep(time.Second * 3)

	// configure Dubbo instance
	zkOption := config_center.WithZookeeper()
	dataIdOption := config_center.WithDataID("dubbo-go-samples-configcenter-zookeeper-go-client")
	addressOption := config_center.WithAddress("127.0.0.1:2181")
	groupOption := config_center.WithGroup("dubbogo")

	ins, err := dubbo.NewInstance(
		dubbo.WithConfigCenter(zkOption, dataIdOption, addressOption, groupOption),
	)
	if err != nil {
		logger.Errorf("Failed to create Dubbo instance: %v", err)
		panic(err)
	}

	// create client
	cli, err := ins.NewClient()
	if err != nil {
		logger.Errorf("Failed to create Dubbo client: %v", err)
		panic(err)
	}

	// create service proxy
	svc, err := greet.NewGreetService(cli)
	if err != nil {
		logger.Errorf("Failed to create GreetService: %v", err)
		panic(err)
	}

	// call remote service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := svc.Greet(ctx, &greet.GreetRequest{Name: "Hello, this is dubbo go client!"})
	if err != nil {
		logger.Errorf("Failed to call Greet service: %v", err)
		return
	}
	logger.Infof("Server response: %s", resp)
}

func writeRuleToConfigCenter() error {
	// connect to ZooKeeper
	c, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*10)
	if err != nil {
		return perrors.Wrap(err, "failed to connect to ZooKeeper")
	}
	defer c.Close() // ensure resource cleanup

	valueBytes := []byte(configCenterZKClientConfig)
	path := "/dubbo/config/dubbogo/dubbo-go-samples-configcenter-zookeeper-go-client"

	// ensure path starts with '/'
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// create parent paths
	if err := createParentPaths(c, path); err != nil {
		return perrors.Wrap(err, "failed to create parent paths")
	}

	// create or update config node
	_, err = c.Create(path, valueBytes, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		if perrors.Is(err, zk.ErrNodeExists) {
			// node exists, update its content
			_, stat, getErr := c.Get(path)
			if getErr != nil {
				return perrors.Wrap(getErr, "failed to get existing node")
			}
			_, setErr := c.Set(path, valueBytes, stat.Version)
			if setErr != nil {
				return perrors.Wrap(setErr, "failed to update existing node")
			}
			logger.Info("Updated existing config node")
		} else {
			return perrors.Wrap(err, "failed to create config node")
		}
	} else {
		logger.Info("Created new config node")
	}

	return nil
}

// helper function to create parent paths
func createParentPaths(c *zk.Conn, path string) error {
	paths := strings.Split(path, "/")
	for idx := 2; idx < len(paths); idx++ {
		tmpPath := strings.Join(paths[:idx], "/")
		_, err := c.Create(tmpPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil && !perrors.Is(err, zk.ErrNodeExists) {
			return perrors.Wrapf(err, "failed to create path: %s", tmpPath)
		}
	}
	return nil
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
        interface: greet.GreetService 
`
