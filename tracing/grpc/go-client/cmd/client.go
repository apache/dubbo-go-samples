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

	"github.com/dubbogo/gost/log"

	"github.com/opentracing/opentracing-go"

	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"

	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

import (
	pb "github.com/apache/dubbo-go-samples/tracing/grpc/protobuf"
)

var (
	grpcGreeterImpl     = new(pb.GreeterClientImpl)
	survivalTimeout int = 10e9
)

func init() {
	config.SetConsumerService(grpcGreeterImpl)
}

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}
	// initJaeger() and initZipkin() can only use one at the same time
	initJaeger()
	//initZipkin()
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "Dubbogo-RPC-Client-Service")
	gxlog.CInfo("\n\n\nstart to test dubbo")
	req := &pb.HelloRequest{
		Name: "Mr.Dubbogo",
	}
	reply, err := grpcGreeterImpl.SayHello(ctx, req)
	span.Finish()
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("client response result: %v\n", reply)
	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP,
		syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		gxlog.CInfo("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			time.AfterFunc(time.Duration(survivalTimeout), func() {
				gxlog.CWarn("app exit now by force...")
				os.Exit(1)
			})
			// The program exits normally or timeout forcibly exits.
			fmt.Println("app exit now...")
			return
		}
	}
}

func initJaeger() {
	cfg := jaegerConfig.Configuration{
		ServiceName: "dobbugoJaegerTracingService",
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeRemote,
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LocalAgentHostPort:  "127.0.0.1:6831",
			LogSpans:            true,
			BufferFlushInterval: 5 * time.Second,
		},
	}
	nativeTracer, _, err := cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		gxlog.CError("unable to create jaeger tracer: %+v\n", err)
	}
	opentracing.SetGlobalTracer(nativeTracer)
}

// nolint
func initZipkin() {
	reporter := zipkinhttp.NewReporter("http://localhost:9411/api/v2/spans")
	endpoint, err := zipkin.NewEndpoint("dobbugoZipkinTracingService", "myservice.mydomain.com:80")
	if err != nil {
		gxlog.CError("unable to create local endpoint: %+v\n", err)
	}
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		gxlog.CError("unable to create tracer: %+v\n", err)
	}
	tracer := zipkinot.Wrap(nativeTracer)
	opentracing.SetGlobalTracer(tracer)
}
