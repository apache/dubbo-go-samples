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
	// reference ctx to avoid unused parameter warning
	_ = ctx
	logger.Info("Received request:" + req.Name)
	resp := &greet.GreetResponse{Greeting: "Hello, this is dubbo go server!" + " I received: " + req.Name}
	return resp, nil
}

func main() {
	// Write configuration to the config center
	if err := writeRuleToConfigCenter(); err != nil {
		logger.Errorf("Failed to write config to config center: %v", err)
		panic(err)
	}

	ins, err := dubbo.NewInstance(
		dubbo.WithConfigCenter(
			config_center.WithZookeeper(),
			config_center.WithDataID("dubbo-go-samples-configcenter-zookeeper-go-server"),
			config_center.WithAddress("127.0.0.1:2181"),
			config_center.WithGroup("dubbogo"),
		),
	)
	if err != nil {
		logger.Errorf("Failed to create dubbo instance: %v", err)
		panic(err)
	}

	srv, err := ins.NewServer()
	if err != nil {
		logger.Errorf("Failed to create server: %v", err)
		panic(err)
	}

	if err = greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
		logger.Errorf("Failed to register service: %v", err)
		panic(err)
	}

	logger.Info("Starting Dubbo-Go server...")
	if err = srv.Serve(); err != nil {
		logger.Errorf("Server failed to serve: %v", err)
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
        interface: greet.GreetService
`

// ensurePath Ensure the path exists (removed because it is unused)

func writeRuleToConfigCenter() error {
	// Connect to Zookeeper
	c, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*10)
	if err != nil {
		return perrors.Wrap(err, "failed to connect to zookeeper")
	}
	defer c.Close() // Ensure the connection is closed

	valueBytes := []byte(configCenterZKServerConfig)
	path := "/dubbo/config/dubbogo/dubbo-go-samples-configcenter-zookeeper-go-server"

	// Ensure the path starts with '/'
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Create parent paths
	if err := createParentPaths(c, path); err != nil {
		return perrors.Wrap(err, "failed to create parent paths")
	}

	// Create or update configuration node
	_, err = c.Create(path, valueBytes, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		if perrors.Is(err, zk.ErrNodeExists) {
			// Node already exists, update configuration
			_, stat, getErr := c.Get(path)
			if getErr != nil {
				return perrors.Wrap(getErr, "failed to get existing node")
			}
			_, setErr := c.Set(path, valueBytes, stat.Version)
			if setErr != nil {
				return perrors.Wrap(setErr, "failed to update existing node")
			}
			logger.Info("Updated existing configuration in config center")
		} else {
			return perrors.Wrap(err, "failed to create configuration node")
		}
	} else {
		logger.Info("Created new configuration in config center")
	}

	if err := waitForConfigReady(c, path, valueBytes, 10*time.Second); err != nil {
		return perrors.Wrap(err, "wait for config ready")
	}

	return nil
}

func waitForConfigReady(c *zk.Conn, path string, expected []byte, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	expectedStr := strings.TrimSpace(string(expected))
	for {
		data, _, err := c.Get(path)
		if err == nil && strings.TrimSpace(string(data)) == expectedStr {
			return nil
		}
		if time.Now().After(deadline) {
			if err != nil {
				return perrors.Wrap(err, "wait for config timeout")
			}
			return perrors.New("wait for config timeout")
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// createParentPaths Create parent paths
func createParentPaths(c *zk.Conn, path string) error {
	paths := strings.Split(path, "/")
	for idx := 2; idx < len(paths); idx++ {
		tmpPath := strings.Join(paths[:idx], "/")
		_, err := c.Create(tmpPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil && !perrors.Is(err, zk.ErrNodeExists) {
			return perrors.Wrapf(err, "failed to create parent path: %s", tmpPath)
		}
	}
	return nil
}
