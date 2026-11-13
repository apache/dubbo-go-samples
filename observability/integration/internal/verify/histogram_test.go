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

package verify

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

import (
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

var sampleBounds = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.3, 1}

// Synthetic data is used only for checker regression tests, not as RPC proof.
func histogramFixture(role string, batches int) string {
	app := "server"
	if role == "consumer" {
		app = "client"
	}
	labels := fmt.Sprintf(`application_name="dubbo-observability-%s",interface="observability.GreetService",method="Greet",group="observability",version="1.0.0",ip="127.0.0.1"`, app)
	prefix := "dubbo_" + role + "_"
	var out strings.Builder
	for suffix, count := range map[string]int{"requests_total": 5, "requests_succeed_total": 4, "requests_failed_total": 1} {
		fmt.Fprintf(&out, "# TYPE %s%s counter\n%s%s{%s} %d\n", prefix, suffix, prefix, suffix, labels, count*batches)
	}
	name := prefix + "rt_milliseconds_histogram_seconds"
	fmt.Fprintf(&out, "# TYPE %s histogram\n", name)
	for i, bound := range []string{"0.005", "0.01", "0.025", "0.05", "0.1", "0.3", "1", "+Inf"} {
		count := 5
		if i < 2 {
			count = 1
		}
		fmt.Fprintf(&out, "%s_bucket{%s,le=%q} %d\n", name, labels, bound, count*batches)
	}
	fmt.Fprintf(&out, "%s_count{%s} %d\n%s_sum{%s} %g\n", name, labels, 5*batches, name, labels, 0.081*float64(batches))
	return out.String()
}

func parseFixture(t *testing.T, text string) map[string]*dto.MetricFamily {
	t.Helper()
	var parser expfmt.TextParser
	families, err := parser.TextToMetricFamilies(strings.NewReader(text))
	if err != nil {
		t.Fatal(err)
	}
	return families
}

func TestHistogramContract(t *testing.T) {
	for _, role := range []string{"provider", "consumer"} {
		t.Run(role, func(t *testing.T) {
			app := "dubbo-observability-server"
			if role == "consumer" {
				app = "dubbo-observability-client"
			}
			target := histogramTarget{role: role, application: app}
			if err := target.check(parseFixture(t, histogramFixture(role, 1)), 5, 1, sampleBounds); err != nil {
				t.Fatal(err)
			}
			target.before = parseFixture(t, histogramFixture(role, 1))
			if err := target.check(parseFixture(t, histogramFixture(role, 2)), 5, 1, sampleBounds); err != nil {
				t.Fatalf("provider baseline must not count old requests: %v", err)
			}
		})
	}
}

func TestHistogramRejectsInvalidExports(t *testing.T) {
	for _, name := range []string{"missing histogram", "missing count", "missing sum", "success only", "wrong bounds", "noncumulative", "wrong infinity", "milliseconds", "negative", "NaN", "extra label", "missing IP"} {
		t.Run(name, func(t *testing.T) {
			families := parseFixture(t, histogramFixture("provider", 1))
			family := families["dubbo_provider_rt_milliseconds_histogram_seconds"]
			metric := family.Metric[0]
			h := metric.Histogram
			switch name {
			case "missing histogram":
				delete(families, family.GetName())
			case "missing count":
				h.SampleCount = nil
			case "missing sum":
				h.SampleSum = nil
			case "success only":
				*h.SampleCount = 4
			case "wrong bounds":
				*h.Bucket[0].UpperBound = 0.006
			case "noncumulative":
				*h.Bucket[1].CumulativeCount = 0
			case "wrong infinity":
				*h.Bucket[len(h.Bucket)-1].CumulativeCount = 4
			case "milliseconds":
				*h.SampleSum = 81
			case "negative":
				*h.SampleSum = -1
			case "NaN":
				families = parseFixture(t, strings.Replace(histogramFixture("provider", 1), "} 0.081", "} NaN", 1))
			case "extra label":
				label, value := "side", "provider"
				metric.Label = append(metric.Label, &dto.LabelPair{Name: &label, Value: &value})
			case "missing IP":
				for _, m := range []*dto.Metric{metric, families["dubbo_provider_requests_total"].Metric[0]} {
					for _, label := range m.Label {
						if label.GetName() == "ip" {
							*label.Value = ""
						}
					}
				}
			}
			target := histogramTarget{role: "provider", application: "dubbo-observability-server"}
			if err := target.check(families, 5, 1, sampleBounds); err == nil {
				t.Fatal("invalid export was accepted")
			}
		})
	}
}

func TestHistogramRunWaitsForAsyncExports(t *testing.T) {
	var complete atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if complete.Load() {
			fmt.Fprint(w, histogramFixture(r.URL.Path[1:], 1))
		}
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	run, err := NewHistogramRun(ctx, server.URL+"/consumer", server.URL+"/provider", sampleBounds)
	if err != nil {
		t.Fatal(err)
	}
	timer := time.AfterFunc(20*time.Millisecond, func() { complete.Store(true) })
	defer timer.Stop()
	if err := run.Check(ctx, 5, 1); err != nil {
		t.Fatal(err)
	}
	ctx, cancelMissing := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancelMissing()
	complete.Store(false)
	if err := run.Check(ctx, 5, 1); err == nil {
		t.Fatal("missing metrics must fail after the deadline")
	}
}
