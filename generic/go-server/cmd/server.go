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
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"dubbo.apache.org/dubbo-go/v3/server"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-samples/generic/go-server/pkg"
)

func main() {
	hessian.RegisterPOJO(&pkg.User{})

	ins, err := dubbo.NewInstance(
		dubbo.WithName("generic-go-server"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(50052),
		),
	)
	if err != nil {
		panic(err)
	}

	srv, err := ins.NewServer(
		server.WithServerSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		panic(err)
	}

	if err := srv.RegisterService(&pkg.UserProvider{},
		server.WithInterface("org.apache.dubbo.samples.UserProvider"),
		server.WithVersion("1.0.0"),
		server.WithGroup("triple"),
	); err != nil {
		panic(err)
	}

	logger.Info("Generic Go server started on port 50052")
	logger.Info("Registry: zookeeper://127.0.0.1:2181")

	if err := srv.Serve(); err != nil {
		logger.Errorf("server stopped: %v", err)
	}
}
