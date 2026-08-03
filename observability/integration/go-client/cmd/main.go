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
	"flag"
	"fmt"
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
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

import (
	"github.com/apache/dubbo-go-samples/observability/integration/internal/tracefields"
	"github.com/apache/dubbo-go-samples/observability/integration/internal/verify"
	observability "github.com/apache/dubbo-go-samples/observability/integration/proto"
)

var (
	requests    = flag.Int("requests", 5, "number of requests to send; 0 runs until interrupted")
	interval    = flag.Duration("interval", time.Second, "delay between requests")
	timeout     = flag.Duration("timeout", 0, "per-request timeout; 0 disables the timeout")
	cancelAfter = flag.Duration("cancel-after", 0, "cancel each request after this duration; 0 disables cancellation")
)

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
	port, err := strconv.Atoi(getEnv("DUBBO_OBSERVABILITY_METRICS_PORT", "9098"))
	if err != nil {
		port = 9098
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
	flag.Parse()

	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo-observability-client"),
		dubbo.WithRegistry(registryOptions()...),
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

	cli, err := ins.NewClient()
	if err != nil {
		panic(err)
	}

	svc, err := observability.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	tracer := otel.Tracer("dubbo-observability-client")
	strict := *timeout == 0 && *cancelAfter == 0
	contractViolations := 0

	for i := 1; *requests == 0 || i <= *requests; i++ {
		requestCtx := context.Background()
		cancel := context.CancelFunc(func() {})
		if *timeout > 0 {
			requestCtx, cancel = context.WithTimeout(requestCtx, *timeout)
		}
		var cancelTimer *time.Timer
		if *cancelAfter > 0 {
			parentCancel := cancel
			var cancelRequest context.CancelFunc
			requestCtx, cancelRequest = context.WithCancel(requestCtx)
			cancel = func() {
				cancelRequest()
				parentCancel()
			}
			cancelTimer = time.AfterFunc(*cancelAfter, cancelRequest)
		}
		ctx, span := tracer.Start(requestCtx, "observability.request")
		name := fmt.Sprintf("request-%d", i)
		if i%5 == 0 {
			name = "error"
		}
		logger.Infof("sending greet request: name=%s%s", name, tracefields.Fields(ctx))

		resp, callErr := svc.Greet(ctx, &observability.GreetRequest{Name: name})
		if strict {
			if err := verify.GreetRequestExpected(name, resp, callErr); err != nil {
				contractViolations++
				if callErr != nil {
					span.RecordError(callErr)
				}
				span.SetStatus(codes.Error, err.Error())
				logger.Errorf("greet request contract violation: %v%s", err, tracefields.Fields(ctx))
			} else if callErr == nil {
				logger.Infof("greet response: %s%s", resp.GetGreeting(), tracefields.Fields(ctx))
			} else {
				logger.Infof("forced error observed as expected: %v%s", callErr, tracefields.Fields(ctx))
			}
		} else if callErr != nil {
			span.RecordError(callErr)
			span.SetStatus(codes.Error, callErr.Error())
			logger.Errorf("greet request failed: %v%s", callErr, tracefields.Fields(ctx))
		} else {
			logger.Infof("greet response: %s%s", resp.GetGreeting(), tracefields.Fields(ctx))
		}
		if cancelTimer != nil {
			cancelTimer.Stop()
		}
		span.End()
		cancel()

		if *requests != 0 && i == *requests {
			break
		}
		time.Sleep(*interval)
	}

	if provider, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); ok {
		_ = provider.Shutdown(context.Background())
	}
	if contractViolations > 0 {
		logger.Fatalf("observability integration observed %d contract violations", contractViolations)
	}
}
