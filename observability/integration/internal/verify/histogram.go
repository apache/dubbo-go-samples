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
	"io"
	"math"
	"net/http"
	"reflect"
	"time"
)

import (
	dto "github.com/prometheus/client_model/go"

	"github.com/prometheus/common/expfmt"
)

type histogramTarget struct {
	role, application, url string
	before                 map[string]*dto.MetricFamily
}

// HistogramRun compares exports before and after a bounded batch of Greet RPCs.
// Baselines allow a provider to remain running between client invocations.
type HistogramRun struct {
	targets []histogramTarget
	bounds  []float64
}

// NewHistogramRun waits for both exporters and captures their initial counters.
func NewHistogramRun(ctx context.Context, consumerURL, providerURL string, bounds []float64) (*HistogramRun, error) {
	run := &HistogramRun{bounds: append([]float64(nil), bounds...), targets: []histogramTarget{
		{role: "consumer", application: "dubbo-observability-client", url: consumerURL},
		{role: "provider", application: "dubbo-observability-server", url: providerURL},
	}}
	for i := range run.targets {
		target := &run.targets[i]
		if err := pollHistogram(ctx, func() error {
			var err error
			target.before, err = readMetrics(ctx, target.url)
			return err
		}); err != nil {
			return nil, fmt.Errorf("%s baseline: %w", target.role, err)
		}
	}
	return run, nil
}

// Check waits for asynchronous metric events after the batch. No new Greet
// calls should be made until this returns, including in continuous mode.
func (r *HistogramRun) Check(ctx context.Context, requests, failures int) error {
	for _, target := range r.targets {
		if err := pollHistogram(ctx, func() error {
			after, err := readMetrics(ctx, target.url)
			if err != nil {
				return err
			}
			return target.check(after, requests, failures, r.bounds)
		}); err != nil {
			return fmt.Errorf("%s histogram verification: %w", target.role, err)
		}
	}
	return nil
}

func pollHistogram(ctx context.Context, check func() error) error {
	for {
		err := check()
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: %v", ctx.Err(), err)
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func readMetrics(ctx context.Context, url string) (map[string]*dto.MetricFamily, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := (&http.Client{Timeout: 2 * time.Second}).Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metrics endpoint returned HTTP %d", response.StatusCode)
	}
	var parser expfmt.TextParser
	return parser.TextToMetricFamilies(io.LimitReader(response.Body, 8<<20))
}

func greetMetric(families map[string]*dto.MetricFamily, name string) (*dto.Metric, error) {
	var found *dto.Metric
	for _, metric := range families[name].GetMetric() {
		labels := metricLabels(metric)
		if labels["interface"] != "observability.GreetService" || labels["method"] != "Greet" {
			continue
		}
		if found != nil {
			return nil, fmt.Errorf("multiple Greet series for %s", name)
		}
		found = metric
	}
	return found, nil
}

func metricLabels(metric *dto.Metric) map[string]string {
	labels := make(map[string]string)
	for _, label := range metric.GetLabel() {
		labels[label.GetName()] = label.GetValue()
	}
	return labels
}

func (t histogramTarget) check(after map[string]*dto.MetricFamily, requests, failures int, bounds []float64) error {
	prefix := "dubbo_" + t.role + "_"
	for suffix, expected := range map[string]int{
		"requests_total": requests, "requests_succeed_total": requests - failures, "requests_failed_total": failures,
	} {
		current, err := greetMetric(after, prefix+suffix)
		if err != nil {
			return err
		}
		previous, err := greetMetric(t.before, prefix+suffix)
		if err != nil {
			return err
		}
		if current == nil && expected != 0 {
			return fmt.Errorf("missing %s", prefix+suffix)
		}
		delta := current.GetCounter().GetValue() - previous.GetCounter().GetValue()
		if delta != float64(expected) {
			return fmt.Errorf("%s delta=%v, want %d", prefix+suffix, delta, expected)
		}
	}
	request, err := greetMetric(after, prefix+"requests_total")
	if err != nil {
		return err
	}
	name := prefix + "rt_milliseconds_histogram_seconds"
	current, err := greetMetric(after, name)
	if err != nil {
		return err
	}
	previous, err := greetMetric(t.before, name)
	if err != nil {
		return err
	}
	if current.GetHistogram() == nil || after[name].GetType() != dto.MetricType_HISTOGRAM {
		return fmt.Errorf("missing Histogram %s", name)
	}
	if current.Histogram.SampleCount == nil || current.Histogram.SampleSum == nil {
		return fmt.Errorf("Histogram count or sum is missing")
	}
	labels := metricLabels(current)
	if !reflect.DeepEqual(labels, metricLabels(request)) || labels["application_name"] != t.application ||
		labels["group"] != "observability" || labels["version"] != "1.0.0" || labels["ip"] == "" {
		return fmt.Errorf("unexpected Histogram labels: %v", labels)
	}
	if previous != nil && !reflect.DeepEqual(labels, metricLabels(previous)) {
		return fmt.Errorf("Histogram labels changed during the batch")
	}
	return checkHistogramBuckets(current.GetHistogram(), previous.GetHistogram(), uint64(requests), bounds)
}

func checkHistogramBuckets(current, before *dto.Histogram, expected uint64, bounds []float64) error {
	if current.GetSampleCount() < before.GetSampleCount() || current.GetSampleCount()-before.GetSampleCount() != expected {
		return fmt.Errorf("Histogram count must include every successful and failed RPC")
	}
	wantBounds := append(append([]float64(nil), bounds...), math.Inf(1))
	if len(current.GetBucket()) != len(wantBounds) || (before != nil && len(before.GetBucket()) != len(wantBounds)) {
		return fmt.Errorf("unexpected Histogram bucket count")
	}
	var previous uint64
	var lower, minimum, maximum float64
	for i, bucket := range current.GetBucket() {
		bound := wantBounds[i]
		if bucket.GetUpperBound() != bound {
			return fmt.Errorf("unexpected bucket boundary %v, want %v", bucket.GetUpperBound(), bound)
		}
		var baseline uint64
		if before != nil {
			if before.Bucket[i].GetUpperBound() != bound {
				return fmt.Errorf("baseline bucket boundary differs")
			}
			baseline = before.Bucket[i].GetCumulativeCount()
		}
		if bucket.GetCumulativeCount() < baseline {
			return fmt.Errorf("bucket counter decreased")
		}
		count := bucket.GetCumulativeCount() - baseline
		if count < previous || count > expected {
			return fmt.Errorf("bucket counts are not cumulative")
		}
		observations := float64(count - previous)
		minimum += observations * lower
		if observations > 0 {
			maximum += observations * bound
		}
		previous, lower = count, bound
	}
	if previous != expected {
		return fmt.Errorf("Histogram +Inf differs from count")
	}
	total := current.GetSampleSum() - before.GetSampleSum()
	if math.IsNaN(total) || math.IsInf(total, 0) || total < 0 || total < minimum-1e-9 || total > maximum+1e-9 {
		return fmt.Errorf("Histogram sum %v is outside bucket bounds in seconds", total)
	}
	return nil
}
