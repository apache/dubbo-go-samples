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
	"os"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	dubbolog "dubbo.apache.org/dubbo-go/v3/logger"
	"dubbo.apache.org/dubbo-go/v3/metrics"
	"dubbo.apache.org/dubbo-go/v3/otel/trace"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"dubbo.apache.org/dubbo-go/v3/server"

	"github.com/dubbogo/gost/log/logger"
)

import (
	adminsample "github.com/apache/dubbo-go-samples/admin/proto"
)

const (
	appName       = "admin-sample-provider"
	serverPort    = 20000
	metricsPort   = 9099
	serviceVer    = "1.0.0"
	serviceGroup  = "admin-sample"
	simulatedWork = 20 * time.Millisecond
)

type GreetServer struct {
	instance string
}

func (s *GreetServer) Greet(ctx context.Context, req *adminsample.GreetRequest) (*adminsample.GreetResponse, error) {
	logger.GetLogger().(dubbolog.CtxLogger).CtxInfof(ctx, "greet name=%s", req.GetName())
	time.Sleep(simulatedWork)
	return &adminsample.GreetResponse{Greeting: "hello " + req.GetName(), Instance: s.instance}, nil
}

func main() {
	nacos := getEnv("NACOS_ADDRESS", "127.0.0.1:8848")
	instance := getEnv("INSTANCE_NAME", "")
	if instance == "" {
		instance, _ = os.Hostname()
	}

	ins, err := dubbo.NewInstance(
		dubbo.WithName(appName),
		dubbo.WithRegistry(registry.WithNacos(), registry.WithAddress(nacos)),
		dubbo.WithProtocol(protocol.WithTriple(), protocol.WithPort(serverPort)),
		dubbo.WithMetrics(
			metrics.WithEnabled(),
			metrics.WithPrometheus(),
			metrics.WithPrometheusExporterEnabled(),
			metrics.WithRegistryEnabled(),
			metrics.WithMetadataEnabled(),
			metrics.WithHistogramEnabled(),
			metrics.WithPort(metricsPort),
			metrics.WithPath("/metrics"),
		),
		dubbo.WithTracing(
			trace.WithEnabled(),
			trace.WithOtlpHttpExporter(),
			trace.WithEndpoint(getEnv("OTLP_ENDPOINT", "127.0.0.1:4318")),
			trace.WithInsecure(),
			trace.WithW3cPropagator(),
			trace.WithAlwaysMode(),
		),
		dubbo.WithLogger(dubbolog.WithZap(), dubbolog.WithTraceIntegration(true)),
	)
	if err != nil {
		panic(err)
	}
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	if err := adminsample.RegisterGreetServiceHandler(srv, &GreetServer{instance: instance}, server.WithVersion(serviceVer), server.WithGroup(serviceGroup)); err != nil {
		panic(err)
	}
	logger.Infof("provider %s serving on port %d", instance, serverPort)
	if err := srv.Serve(); err != nil {
		logger.Fatalf("serve failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
