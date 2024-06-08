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
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	"fmt"
	greet "github.com/apache/dubbo-go-samples/java_interop/non-protobuf-triple/proto"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: fmt.Sprintf("The name of the request is:	 %s", req.Name)}
	return resp, nil
}

func main() {
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(50052),
			protocol.WithTriple(),
		),
		server.WithServerSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		panic(err) // 这里也有错误检查，确保每次赋值后都检查了 err
	}

	if err := greet.RegisterGreetingsServiceHandler(srv, &GreetTripleServer{}); err != nil {
		panic(err) // 这里的错误处理也是正确的
	}

	if err := srv.Serve(); err != nil {
		panic(err) // 启动服务器的错误处理也没有问题
	}
}
