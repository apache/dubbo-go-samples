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
	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"math/rand"
	"os"
	"time"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
	"github.com/dubbogo/gost/log"
	// tracing zipkin & prometheus
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

import (
	"github.com/apache/dubbo-go-samples/group/go-client/pkg"
)

var userProviderA = new(pkg.UserProviderGroupA)
var userProviderB = new(pkg.UserProviderGroupB)

func init() {
	config.SetConsumerService(userProviderA)
	config.SetConsumerService(userProviderB)
	hessian.RegisterPOJO(&pkg.User{})
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {

	config.Load()

	// register zipkin & prometheus exporter
	registerZipkin()
	registerPrometheus()

	time.Sleep(3 * time.Second)

	gxlog.CInfo("\n\n\nstart to test dubbo")

	getUserAll(context.Background())
}

func getUserAll(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getUserAll")

	getUserA(ctx)
	span.Finish()

	getUserB(ctx)
	span.Finish()

	getUserA(ctx)
	span.Finish()

	getUserB(ctx)
	span.Finish()
}

func getUserA(ctx context.Context) {

	time.Sleep(time.Duration(rand.Intn(977)+300) * time.Millisecond)
	user := &pkg.User{}
	err := userProviderA.GetUser(ctx, []interface{}{"A001"}, user)
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("response result: %v\n", user)

}
func getUserB(ctx context.Context) {

	time.Sleep(time.Duration(rand.Intn(977)+300) * time.Millisecond)
	user := &pkg.User{}
	err := userProviderB.GetUser(ctx, []interface{}{"A001"}, user)
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("response result: %v\n", user)
}

// zipkin / opentracing specific stuff
func registerZipkin() {
	// set up a span reporter
	reporter := zipkinhttp.NewReporter("http://localhost:9411/api/v2/spans")

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint("go-client", "localhost:80")
	if err != nil {
		gxlog.CError("unable to create local endpoint: %+v\n", err)
	}

	// set sampler , default AlwaysSample
	//sampler := zipkin.NewModuloSampler(1)

	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))

	if err != nil {
		gxlog.CError("unable to create tracer: %+v\n", err)
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)

	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)
}

// register prometheus exporter for zipkin
func registerPrometheus() *prometheus.Exporter {
	pe, err := prometheus.NewExporter(prometheus.Options{Namespace: "go-client"})
	if err != nil {
		gxlog.CError("Failed to create Prometheus exporter: %v", err)
	}
	view.RegisterExporter(pe)
	return pe
}
