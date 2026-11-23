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

// The client sends one untagged and one "gray" tagged call per second, so the
// effect of rules created in Admin shows up directly in its log.
package main

import (
	"context"
	"os"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	dubbolog "dubbo.apache.org/dubbo-go/v3/logger"
	"dubbo.apache.org/dubbo-go/v3/metrics"
	"dubbo.apache.org/dubbo-go/v3/otel/trace"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	adminsample "github.com/apache/dubbo-go-samples/admin/proto"
)

const (
	appName      = "admin-sample-consumer"
	metricsPort  = 9099
	serviceVer   = "1.0.0"
	serviceGroup = "admin-sample"
	grayTag      = "gray"
	ruleGroup    = "dubbo"
	interval     = time.Second
	callTimeout  = 3 * time.Second
)

func main() {
	nacos := getEnv("NACOS_ADDRESS", "127.0.0.1:8848")
	ins, err := dubbo.NewInstance(
		dubbo.WithName(appName),
		dubbo.WithRegistry(registry.WithNacos(), registry.WithAddress(nacos)),
		// Routers watch Admin's rules only when a config center is configured;
		// Admin writes them to the "dubbo" group.
		dubbo.WithConfigCenter(config_center.WithNacos(), config_center.WithAddress(nacos), config_center.WithGroup(ruleGroup)),
		dubbo.WithMetrics(
			metrics.WithEnabled(),
			metrics.WithPrometheus(),
			metrics.WithPrometheusExporterEnabled(),
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
	cli, err := ins.NewClient(client.WithClientVersion(serviceVer), client.WithClientGroup(serviceGroup))
	if err != nil {
		panic(err)
	}
	svc, err := adminsample.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		call(svc, "")
		call(svc, grayTag)
	}
}

func call(svc adminsample.GreetService, tag string) {
	ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
	defer cancel()
	label := "untagged"
	if tag != "" {
		label = "tag=" + tag
		ctx = context.WithValue(ctx, constant.AttachmentKey, map[string]string{constant.Tagkey: tag})
	}
	resp, err := svc.Greet(ctx, &adminsample.GreetRequest{Name: "admin"})
	if err != nil {
		logger.Warnf("%s -> error: %v", label, err)
		return
	}
	logger.Infof("%s -> %s", label, resp.GetInstance())
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
