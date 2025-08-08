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
	apolloMetaAddress = "127.0.0.1:8080"
	apolloAppID       = "SampleApp"
	apolloCluster     = "default"
	apolloNamespace   = "dubbo.yml"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: "Hello, " + req.Name}
	return resp, nil
}

func main() {
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
