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
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"dubbo.apache.org/dubbo-go/v3/server"
	"github.com/apache/dubbo-go-samples/helloworld/go-server/handler"
	"github.com/apache/dubbo-go-samples/helloworld/proto/greettriple"
	"github.com/dubbogo/gost/log/logger"
	"time"
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
			registry.WithTimeout(3*time.Second),
			registry.WithGroup("myGroup"),
			registry.WithRegisterInterface(),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20000),
		),
	)
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	if err := greettriple.RegisterGreetServiceHandler(srv, &handler.GreetTripleServer{}); err != nil {
		panic(err)
	}
	if err := greettriple.RegisterGreetServiceHandler(srv, &handler.GreetTripleServer{},
		server.WithGroup("myInterfaceGroup"),
		server.WithVersion("myInterfaceVersion"),
	); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
