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
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"github.com/SkyAPM/go2sky"
	dubbo_go "github.com/SkyAPM/go2sky-plugins/dubbo-go"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/apache/dubbo-go-samples/api"
	"github.com/dubbogo/gost/log/logger"
	"log"
)

var grpcGreeterImpl = new(api.GreeterClientImpl)

// export DUBBO_GO_CONFIG_PATH=PATH_TO_SAMPLES/skywalking/go-client/conf/dubbogo.yml
func main() {
	// set dubbogo configs ...
	config.SetConsumerService(grpcGreeterImpl)
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

	// set dubbogo plugin client tracer
	err = dubbo_go.SetClientTracer(tracer)
	if err != nil {
		log.Fatalf("set tracer error: %v \n", err)
	}

	// set extra tags and report tags
	dubbo_go.SetClientExtraTags("extra-tags", "client")
	dubbo_go.SetClientReportTags("release")

	logger.Info("start to test dubbo")
	req := &api.HelloRequest{
		Name: "laurence",
	}

	reply, err := grpcGreeterImpl.SayHello(context.Background(), req)
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("client response result: %v\n", reply)

	select {}
}
