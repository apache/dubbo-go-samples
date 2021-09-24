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

	"math/rand"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/metrics/prometheus"
)

import (
	"github.com/apache/dubbo-go-samples/api"
)

type GreeterProvider struct {
	api.GreeterProviderBase
}

func (s *GreeterProvider) SayHello(ctx context.Context, in *api.HelloRequest) (*api.User, error) {
	logger.Infof("Dubbo3 GreeterProvider get user name = %s\n", in.Name)
	return &api.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, nil
}

// export DUBBO_GO_CONFIG_PATH= PATH_TO_SAMPLES/helloworld/go-server/conf/dubbogo.yml
func main() {
	config.SetProviderService(&GreeterProvider{})
	if err := config.Load(); err != nil {
		panic(err)
	}
	for {
		time.Sleep(time.Second)
		prometheus.IncSummary("test_summary", rand.Float64())
		prometheus.IncSummaryWithLabel("test_summary_with_label", rand.Float64(), map[string]string{
			"summarylabel1":"value1",
			"summarylabel2":"value2",
		})
		prometheus.IncCounter("test_counter")
		prometheus.IncCounterWithLabel("test_counter_with_label",map[string]string{
			"counterlabel1":"value1",
			"counterlabel2":"value2",
		})

		prometheus.SetGauge("test_gauge", rand.Float64())
		prometheus.SetGaugeWithLabel("test_gauge_with_label", rand.Float64(), map[string]string{
			"gaugelabel1":"value1",
			"gaugelabel2":"value2",
		})
	}
}