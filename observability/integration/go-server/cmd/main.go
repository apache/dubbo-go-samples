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
	"os"
	"strconv"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	dubbolog "dubbo.apache.org/dubbo-go/v3/logger"
	"dubbo.apache.org/dubbo-go/v3/metrics"
	"dubbo.apache.org/dubbo-go/v3/otel/trace"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-samples/observability/integration/internal/tracefields"
	observability "github.com/apache/dubbo-go-samples/observability/integration/proto"
)

const (
	serverPort        = 20000
	serverMetricsPort = 9099
)

type GreetServer struct{}

func (s *GreetServer) Greet(ctx context.Context, req *observability.GreetRequest) (*observability.GreetResponse, error) {
	if req == nil {
		err := errors.New("observability sample received nil request")
		logger.Errorf("greet request failed: %v%s", err, tracefields.Fields(ctx))
		return nil, err
	}

	logger.Infof("greet request received: name=%s%s", req.GetName(), tracefields.Fields(ctx))

	if req.GetName() == "error" {
		err := triple_protocol.NewError(triple_protocol.CodeBizError, errors.New("observability sample forced error"))
		logger.Errorf("observability sample forced error: %v%s", err, tracefields.Fields(ctx))
		return nil, err
	}

	time.Sleep(20 * time.Millisecond)
	logger.Infof("greet request completed%s", tracefields.Fields(ctx))
	return &observability.GreetResponse{Greeting: "hello " + req.GetName()}, nil
}

func traceOptions() []trace.Option {
	return []trace.Option{
		trace.WithEnabled(),
		trace.WithOtlpHttpExporter(),
		trace.WithW3cPropagator(),
		trace.WithAlwaysMode(),
		trace.WithEndpoint(getEnv("DUBBO_OBSERVABILITY_OTLP_ENDPOINT", "127.0.0.1:4318")),
		trace.WithInsecure(),
	}
}

func metricsOptions() []metrics.Option {
	port, err := strconv.Atoi(getEnv("DUBBO_OBSERVABILITY_METRICS_PORT", strconv.Itoa(serverMetricsPort)))
	if err != nil {
		port = serverMetricsPort
	}
	return []metrics.Option{
		metrics.WithEnabled(),
		metrics.WithPrometheus(),
		metrics.WithPrometheusExporterEnabled(),
		metrics.WithRegistryEnabled(),
		metrics.WithMetadataEnabled(),
		metrics.WithPort(port),
		metrics.WithPath("/prometheus"),
	}
}

func registryOptions() []registry.Option {
	return []registry.Option{
		registry.WithNacos(),
		registry.WithAddress(getEnv("DUBBO_OBSERVABILITY_REGISTRY_ADDRESS", "127.0.0.1:8848")),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo-observability-server"),
		dubbo.WithRegistry(registryOptions()...),
		dubbo.WithProtocol(protocol.WithPort(serverPort), protocol.WithTriple()),
		dubbo.WithTracing(traceOptions()...),
		dubbo.WithMetrics(metricsOptions()...),
		dubbo.WithLogger(
			dubbolog.WithZap(),
			dubbolog.WithLevel("debug"),
		),
	)
	if err != nil {
		panic(err)
	}

	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}

	if err := observability.RegisterGreetServiceHandler(srv, &GreetServer{}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Fatalf("server serve failed: %v", err)
	}
}
