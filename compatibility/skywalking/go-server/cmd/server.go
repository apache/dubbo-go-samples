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
	"log"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/SkyAPM/go2sky"
	dubbo_go "github.com/SkyAPM/go2sky-plugins/dubbo-go"
	"github.com/SkyAPM/go2sky/reporter"

	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-samples/compatibility/api"
)

type Greeter struct {
	api.UnimplementedGreeterServer
}

func (s *Greeter) SayHello(ctx context.Context, in *api.HelloRequest) (*api.User, error) {
	logger.Infof("Dubbo3 GreeterProvider get user name = %s\n", in.Name)
	return &api.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, nil
}

func (s *Greeter) SayHelloStream(svr api.Greeter_SayHelloStreamServer) error {
	c, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go3 GreeterProvider recv 1 user, name = %s\n", c.Name)
	c2, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go3 GreeterProvider recv 2 user, name = %s\n", c2.Name)

	err = svr.Send(&api.User{
		Name: "hello " + c.Name,
		Age:  18,
		Id:   "123456789",
	})
	if err != nil {
		return err
	}
	c3, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go3 GreeterProvider recv 3 user, name = %s\n", c3.Name)

	err = svr.Send(&api.User{
		Name: "hello " + c2.Name,
		Age:  19,
		Id:   "123456789",
	})
	if err != nil {
		return err
	}
	return nil
}

// export DUBBO_GO_CONFIG_PATH=PATH_TO_SAMPLES/skywalking/go-server/conf/dubbogo.yml
func main() {
	// set dubbogo configs ...
	config.SetProviderService(&Greeter{})
	if err := config.Load(); err != nil {
		panic(err)
	}

	// setup reporter, use gRPC reporter for production
	report, err := reporter.NewGRPCReporter("YOUR_SKYWALKING_DOMAIN_NAME_OR_IP:11800")
	if err != nil {
		log.Fatalf("new reporter error: %v \n", err)
	}

	// setup tracer
	tracer, err := go2sky.NewTracer("dubbo-go-skywalking-sample-tracer", go2sky.WithReporter(report))
	if err != nil {
		log.Fatalf("crate tracer error: %v \n", err)
	}

	// set dubbogo plugin server tracer
	err = dubbo_go.SetServerTracer(tracer)
	if err != nil {
		log.Fatalf("set tracer error: %v \n", err)
	}

	// set extra tags and report tags
	dubbo_go.SetServerExtraTags("extra-tags", "server")
	dubbo_go.SetServerReportTags("release")

	select {}
}
