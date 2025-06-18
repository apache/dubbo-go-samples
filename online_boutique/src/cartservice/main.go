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
	"github.com/apache/dubbo-go-samples/online_boutique_demo/cartservice/handler"
	hipstershop "github.com/apache/dubbo-go-samples/online_boutique_demo/cartservice/proto"
	"github.com/dubbogo/gost/log/logger"
	_ "github.com/dubbogo/gost/log/logger"
	"os"
)

func main() {
	regAddr := os.Getenv("DUBBO_REGISTRY_ADDRESS")
	if regAddr == "" {
		regAddr = "127.0.0.1:2181"
	}

	ins, err := dubbo.NewInstance(
		dubbo.WithName("cartservice"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress(regAddr),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20001),
		),
	)
	if err != nil {
		panic(err)
	}
	//server
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}

	if err := hipstershop.RegisterCartServiceHandler(srv, handler.NewCartService()); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}

}
