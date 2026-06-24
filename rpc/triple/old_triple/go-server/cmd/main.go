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
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
)

import (
	greet "github.com/apache/dubbo-go-samples/rpc/triple/old_triple/proto"
)

type GreetDubbo3Server struct {
	greet.UnimplementedGreetServiceServer
}

func (srv *GreetDubbo3Server) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	return &greet.GreetResponse{Greeting: req.Name}, nil
}

var greetServiceInfo = server.ServiceInfo{
	InterfaceName: "greet.GreetService",
	ServiceType:   (*greet.GreetServiceServer)(nil),
	Methods: []server.MethodInfo{
		{
			Name: "Greet",
			Type: constant.CallUnary,
			ReqInitFunc: func() any {
				return new(greet.GreetRequest)
			},
			MethodFunc: func(ctx context.Context, args []any, handler any) (any, error) {
				req := args[0].(*greet.GreetRequest)
				res, err := handler.(greet.GreetServiceServer).Greet(ctx, req)
				if err != nil {
					return nil, err
				}
				return triple_protocol.NewResponse(res), nil
			},
		},
	},
}

func main() {
	dubbo.SetProviderServiceWithInfo(&GreetDubbo3Server{}, &greetServiceInfo)
	if err := dubbo.Load(dubbo.WithPath("./go-server/conf/dubbogo.yml")); err != nil {
		panic(err)
	}
	select {}
}
