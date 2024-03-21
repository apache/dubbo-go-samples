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
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	greet "github.com/apache/dubbo-go-samples/java_interop/protobuf-triple/go/proto"
)

// export DUBBO_GO_CONFIG_PATH=$PATH_TO_SAMPLES/rpc/triple/pb/dubbogo-java/go-client/conf/dubbogo.yml

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) SayHello(ctx context.Context, req *greet.HelloRequest) (*greet.HelloReply, error) {
	resp := &greet.HelloReply{Message: req.Name}
	return resp, nil
}

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("org.apache.dubbo.sample.GreeterImpl"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(50052),
		),
	)
	// create a server with registry and protocol set above
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	// register a service to server
	if err := greet.RegisterGreeterHandler(srv, &GreetTripleServer{}); err != nil {
		panic(err)
	}
	// start the server
	if err := srv.Serve(); err != nil {
		panic(err)
	}
}
