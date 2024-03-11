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
	"github.com/pkg/errors"
	"math/rand"
	"time"

	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/metrics"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	greet "github.com/apache/dubbo-go-samples/helloworld/proto"
	"github.com/dubbogo/gost/log/logger"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(_ context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(101) > 99 { // mock error here
		return nil, errors.New("random error")
	}
	time.Sleep(10 * time.Millisecond) // mock business delay
	return resp, nil
}

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithMetrics(
			metrics.WithEnabled(),
			metrics.WithPrometheus(),                // set prometheus metric
			metrics.WithPrometheusExporterEnabled(), // enable prometheus exporter
			metrics.WithPort(9099),                  // prometheus http exporter listen at 9099
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
	srv, err := ins.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000), // triple protocol port
			protocol.WithTriple(),
		),
	)
	if err != nil {
		panic(err)
	}

	if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
