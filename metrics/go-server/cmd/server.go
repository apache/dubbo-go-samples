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
	"github.com/dubbogo/gost/log/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/metrics/prometheus"
)

import (
	"github.com/apache/dubbo-go-samples/api"
)

type GreeterProvider struct {
	api.UnimplementedGreeterServer
}

func (s *GreeterProvider) SayHello(ctx context.Context, in *api.HelloRequest) (*api.User, error) {
	logger.Infof("Dubbo3 GreeterProvider get user name = %s\n", in.Name)
	return &api.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, nil
}

/* metrics in localhost:9090:
userCounterVec# HELP dubbo_test_counter
# TYPE dubbo_test_counter counter
dubbo_test_counter 36
# HELP dubbo_test_counter_with_label
# TYPE dubbo_test_counter_with_label counter
dubbo_test_counter_with_label{counterlabel1="value1",counterlabel2="value2"} 36
# HELP dubbo_test_gauge
# TYPE dubbo_test_gauge gauge
dubbo_test_gauge 0.4503573689185021
# HELP dubbo_test_gauge_with_label
# TYPE dubbo_test_gauge_with_label gauge
dubbo_test_gauge_with_label{gaugelabel1="value1",gaugelabel2="value2"} 0.3152981056786652
# HELP dubbo_test_summary
# TYPE dubbo_test_summary summary
dubbo_test_summary_sum 19.88640045117142
dubbo_test_summary_count 36
# HELP dubbo_test_summary_with_label
# TYPE dubbo_test_summary_with_label summary
dubbo_test_summary_with_label{summarylabel1="value1",summarylabel2="value2",quantile="0.5"} 0.6103291015556667
dubbo_test_summary_with_label{summarylabel1="value1",summarylabel2="value2",quantile="0.75"} 0.7991992117286294
dubbo_test_summary_with_label{summarylabel1="value1",summarylabel2="value2",quantile="0.9"} 0.937731048284269
dubbo_test_summary_with_label{summarylabel1="value1",summarylabel2="value2",quantile="0.98"} 0.9806989584755814
dubbo_test_summary_with_label{summarylabel1="value1",summarylabel2="value2",quantile="0.99"} 0.9806989584755814
dubbo_test_summary_with_label{summarylabel1="value1",summarylabel2="value2",quantile="0.999"} 0.9806989584755814
dubbo_test_summary_with_label_sum{summarylabel1="value1",summarylabel2="value2"} 18.896320346046032
dubbo_test_summary_with_label_count{summarylabel1="value1",summarylabel2="value2"} 36
*/
func main() {
	config.SetProviderService(&GreeterProvider{})
	if err := config.Load(); err != nil {
		panic(err)
	}
	for {
		// metrics refresh per second
		time.Sleep(time.Second)
		prometheus.IncSummary("test_summary", rand.Float64())
		prometheus.IncSummaryWithLabel("test_summary_with_label", rand.Float64(), map[string]string{
			"summarylabel1": "value1", // label and value for this summary
			"summarylabel2": "value2",
		})

		prometheus.IncCounter("test_counter")
		prometheus.IncCounterWithLabel("test_counter_with_label", map[string]string{
			"counterlabel1": "value1", // label and value for this counter
			"counterlabel2": "value2",
		})

		prometheus.SetGauge("test_gauge", rand.Float64())
		prometheus.SetGaugeWithLabel("test_gauge_with_label", rand.Float64(), map[string]string{
			"gaugelabel1": "value1",
			"gaugelabel2": "value2",
		})
	}
}
