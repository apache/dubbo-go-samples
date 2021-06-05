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
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	"github.com/apache/dubbo-go/common/logger"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/protocol/grpc"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
	gxlog "github.com/dubbogo/gost/log"
)

var (
	survivalTimeout int = 10e9
)
var grpcGreeterImpl = new(GreeterClientImpl)

func init() {
	config.SetConsumerService(grpcGreeterImpl)
}

// they are necessary:
// 		export CONF_CONSUMER_FILE_PATH="xxx"
// 		export APP_LOG_CONF_FILE="xxx"
func main() {
	config.Load()
	time.Sleep(time.Second)

	gxlog.CInfo("\n\n\n===== start to test SayHelloTwoSidesStream ======")
	reply := &HelloReply{}
	stream, err := grpcGreeterImpl.SayHelloTwoSidesStream(context.TODO())
	if err != nil {
		logger.Errorf("stream get err = %s", err.Error())
		return
	}

	if err := stream.Send(&HelloRequest{Name: "request 1"}); err != nil {
		logger.Errorf("send stream req err = %s", err.Error())
	}
	if err := stream.Send(&HelloRequest{Name: "request 2"}); err != nil {
		logger.Errorf("send stream req err = %s", err.Error())
	}
	reply, err = stream.Recv()
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("client response result: %v\n", reply)
	reply, err = stream.Recv()
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("client response result: %v\n", reply)

	gxlog.CInfo("\n\n\n===== start to test SayHelloClientStream =====")
	clientStream, err := grpcGreeterImpl.SayHelloClientStream(context.TODO())
	if err != nil {
		logger.Errorf("stream get err = %s", err.Error())
		return
	}

	if err := clientStream.Send(&HelloRequest{Name: "request 1"}); err != nil {
		logger.Errorf("send stream req err = %s", err.Error())
	}
	if err := clientStream.Send(&HelloRequest{Name: "request 2"}); err != nil {
		logger.Errorf("send stream req err = %s", err.Error())
	}
	err = clientStream.RecvMsg(reply)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("client response result: %v\n", reply)

	gxlog.CInfo("\n\n\n===== start to test SayHelloServerStream =====")
	req := &HelloRequest{}
	serverStream, err := grpcGreeterImpl.SayHelloServerStream(context.TODO(), req)
	if err != nil {
		logger.Errorf("stream get err = %s", err.Error())
		return
	}
	reply, err = serverStream.Recv()
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("client response result: %v\n", reply)
	reply, err = serverStream.Recv()
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("client response result: %v\n", reply)

	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP,
		syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			time.AfterFunc(time.Duration(survivalTimeout), func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})

			// The program exits normally or timeout forcibly exits.
			fmt.Println("app exit now...")
			return
		}
	}
}
