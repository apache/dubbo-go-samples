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
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	pb "github.com/apache/dubbo-go-samples/tls/grpc/protobuf"
)

var (
	survivalTimeout = int(3 * time.Second)
)

type GreeterProvider struct {
	pb.GreeterProviderBase
}

func (g *GreeterProvider) SayHello(ctx context.Context, req *pb.HelloRequest) (reply *pb.HelloReply, err error) {
	fmt.Printf("req: %v", req)
	return &pb.HelloReply{Message: "this is message from reply"}, nil
}

// need to setup environment variable "CONF_PROVIDER_FILE_PATH" to "conf/server.yml" before run
func main() {
	config.SetProviderService(&GreeterProvider{})
	if err := config.Load(); err != nil {
		panic(err)
	}
	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)
	for {
		sig := <-signals
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			time.AfterFunc(time.Duration(survivalTimeout), func() {
				os.Exit(1)
			})

			// The program exits normally or timeout forcibly exits.
			fmt.Println("provider app exit now...")
			return
		}
	}
}
