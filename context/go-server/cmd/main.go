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
	triple "dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
	"fmt"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	greet "github.com/apache/dubbo-go-samples/context/proto"
	"github.com/dubbogo/gost/log/logger"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	data, _ := triple.FromIncomingContext(ctx)
	ctx = triple.AppendToOutgoingContext(ctx, "OutgoingContextKey1", "OutgoingDataVal1", "OutgoingContextKey2", "OutgoingDataVal2")
	var value1, value2, value3 string
	if values, ok := data["testkey1"]; ok && len(values) > 0 {
		value1 = values[0]
		logger.Infof("testkey1: %s", value1)
	}
	if values, ok := data["testkey2"]; ok && len(values) > 0 {
		value2 = values[0]
		logger.Infof("testkey2: %s", value2)
	}
	if values, ok := data["testkey3"]; ok && len(values) > 0 {
		value3 = values[0]
		logger.Infof("testkey3: %s", value3)
	}

	respStr := fmt.Sprintf("name: %s, testKey1: %s, testKey2: %s", req.Name, value1, value2)
	resp := &greet.GreetResponse{Greeting: respStr}
	return resp, nil
}

func main() {
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000),
			protocol.WithTriple(),
		),
	)
	if err != nil {
		panic(err)
	}

	if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
