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
	"os"
	"os/signal"
	"syscall"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/common"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"dubbo.apache.org/dubbo-go/v3/server"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log/logger"
)

import (
	pkg2 "github.com/apache/dubbo-go-samples/generic/go-server/pkg"
)

const (
	RegistryAddress = "127.0.0.1:2181"
	ServerName      = "generic-server"
	DubboServerPort = 20004
)

func main() {
	hessian.RegisterPOJO(&pkg2.User{})

	ins := createDubboInstance()

	srv, err := ins.NewServer()
	if err != nil {
		logger.Fatalf("Failed to create server: %v", err)
	}

	registerService(srv, &pkg2.UserProvider{})

	go func() {
		logger.Info("Starting Dubbo Protocol Server...")
		if err := srv.Serve(); err != nil {
			logger.Errorf("Dubbo server failed: %v", err)
		}
	}()

	waitForShutdown()
}

func createDubboInstance() *dubbo.Instance {
	ins, err := dubbo.NewInstance(
		dubbo.WithName(ServerName),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress(RegistryAddress),
		),
		dubbo.WithProtocol(
			protocol.WithID("dubbo"),
			protocol.WithDubbo(),
			protocol.WithPort(DubboServerPort),
		),
	)
	if err != nil {
		logger.Fatalf("Failed to create instance: %v", err)
	}
	return ins
}

func registerService(srv *server.Server, service *pkg2.UserProvider) {
	serviceInfo := &common.ServiceInfo{
		InterfaceName: "org.apache.dubbo.samples.UserProvider",
		ServiceType:   service,
		Methods: []common.MethodInfo{
			{Name: "GetUser1", Type: "normal", Meta: map[string]interface{}{"params": []string{"java.lang.String"}}},
			{Name: "GetUser2", Type: "normal", Meta: map[string]interface{}{"params": []string{"java.lang.String", "java.lang.String"}}},
			{Name: "GetUser3", Type: "normal", Meta: map[string]interface{}{"params": []string{"int"}}},
			{Name: "GetUser4", Type: "normal", Meta: map[string]interface{}{"params": []string{"int", "java.lang.String"}}},
			{Name: "GetOneUser", Type: "normal", Meta: map[string]interface{}{"params": []string{}}},

			{Name: "GetUsers", Type: "normal", Meta: map[string]interface{}{"params": []string{"[Ljava.lang.String;"}}},
			{Name: "GetUsersMap", Type: "normal", Meta: map[string]interface{}{"params": []string{"[Ljava.lang.String;"}}},

			{Name: "QueryUser", Type: "normal", Meta: map[string]interface{}{"params": []string{"org.apache.dubbo.samples.User"}}},
			{Name: "QueryUsers", Type: "normal", Meta: map[string]interface{}{"params": []string{"[]org.apache.dubbo.samples.User"}}},
			{Name: "QueryAll", Type: "normal", Meta: map[string]interface{}{"params": []string{}}},
		},
		Meta: map[string]interface{}{
			"version":  "1.0.0",
			"group":    "dubbo",
			"protocol": "dubbo",
		},
	}

	serviceOpts := []server.ServiceOption{
		server.WithInterface("org.apache.dubbo.samples.UserProvider"),
		server.WithVersion("1.0.0"),
		server.WithGroup("dubbo"),
		server.WithProtocolIDs([]string{"dubbo"}),
	}

	if err := srv.Register(service, serviceInfo, serviceOpts...); err != nil {
		logger.Fatalf("Failed to register Dubbo service: %v", err)
	}
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	sig := <-sigChan
	logger.Infof("Received signal: %s, shutting down...", sig.String())
	os.Exit(0)
}
