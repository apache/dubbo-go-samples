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
	"dubbo.apache.org/dubbo-go/v3/metrics"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/helloworld/proto"
)

func main() {
	zookeeper := os.Getenv("ZOOKEEPER_ADDRESS")
	if zookeeper == "" {
		zookeeper = "localhost"
	}
	ins, err := dubbo.NewInstance(
		dubbo.WithRegistry(
			registry.WithAddress("zookeeper://"+zookeeper+":2181"),
		),
		dubbo.WithMetrics(
			metrics.WithEnabled(),
			metrics.WithPrometheus(),                // set prometheus metric
			metrics.WithPrometheusExporterEnabled(), // enable prometheus exporter
			metrics.WithPort(9097),                  // prometheus http exporter listen at 9097
			metrics.WithPath("/prometheus"),         // prometheus http exporter url path
			metrics.WithMetadataEnabled(),           // enable metadata center metrics
			metrics.WithRegistryEnabled(),           // enable registry metrics
			metrics.WithConfigCenterEnabled(),       // enable config center metrics

			metrics.WithPrometheusPushgatewayEnabled(), // enable prometheus pushgateway
			metrics.WithPrometheusGatewayUsername("username"),
			metrics.WithPrometheusGatewayPassword("1234"),
			metrics.WithPrometheusGatewayUrl("127.0.0.1:9091"), // host:port or ip:port,“http://” is added automatically,do not include the “/metrics/jobs/…” part
			metrics.WithPrometheusGatewayInterval(time.Second*10),
			metrics.WithPrometheusGatewayJob("push"), // set a metric job label, job=push to metric

			metrics.WithAggregationEnabled(), // enable rpc metrics aggregations，Most of the time there is no need to turn it on
			metrics.WithAggregationTimeWindowSeconds(30),
			metrics.WithAggregationBucketNum(10), // agg bucket num
		),
	)
	if err != nil {
		panic(err)
	}
	cli, err := ins.NewClient()
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	for true {
		resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
		if err != nil {
			logger.Error(err)
		} else {
			logger.Infof("Greet response: %s", resp.Greeting)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
