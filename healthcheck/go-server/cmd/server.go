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
	"errors"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
)

import (
	greet "github.com/apache/dubbo-go-samples/healthcheck/proto"
	"github.com/dubbogo/gost/log/logger"
)

// GreetTripleServer implements the greet.GreetServiceServer interface for the Triple protocol.
type GreetTripleServer struct {
}

// Greet handles the Greet request and returns a greeting message.
func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	if req == nil {
		return nil, errors.New("GreetRequest is nil")
	}
	resp := &greet.GreetResponse{Greeting: req.Name}
	return resp, nil
}

// runServer starts and runs the Dubbo-go Triple server.
func runServer() error {
	ins, err := dubbo.NewInstance(
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20000),
		),
	)
	if err != nil {
		return err
	}
	srv, err := ins.NewServer()
	if err != nil {
		return err
	}
	if err = greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
		return err
	}
	return srv.Serve()
}

func main() {
	// Start the server and handle errors
	if err := runServer(); err != nil {
		logger.Errorf("failed to start server: %v", err)
	}
}
