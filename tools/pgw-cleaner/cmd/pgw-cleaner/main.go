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
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

import (
	"github.com/dubbogo/gost/log/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/common/expfmt"
)

// Configuration flags
var (
	pushgw      = flag.String("pushgw", "http://127.0.0.1:9091", "Pushgateway address")
	ttl         = flag.Int("ttl", 3600, "TTL seconds for job/instance (default 300s)")
	daemon      = flag.Bool("daemon", false, "Run in daemon mode")
	interval    = flag.Int("interval", 60, "Interval seconds for daemon mode")
	jobPrefix   = flag.String("job-prefix", "dubbo_", "Only clean jobs with this prefix")
	metricsPort = flag.Int("metrics-port", 9105, "Cleaner metrics port (daemon mode)")
)

// Cleaner metrics
var (
	checkedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pgw_cleaner_checked_total",
		Help: "Total number of job/instance checked",
	})
	deletedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pgw_cleaner_deleted_total",
		Help: "Total number of job/instance deleted",
	})
	errorsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pgw_cleaner_errors_total",
		Help: "Total number of errors occurred",
	})
)

func init() {
	prometheus.MustRegister(checkedTotal, deletedTotal, errorsTotal)
}

// parsePushgatewayMetrics fetches and parses metrics from Pushgateway
func parsePushgatewayMetrics(pushgwURL string) (map[[2]string]float64, error) {
	resp, err := http.Get(pushgwURL + "/metrics")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	parser := expfmt.TextParser{}
	mf, err := parser.TextToMetricFamilies(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	result := make(map[[2]string]float64)
	for _, fam := range mf {
		// Check all job_pushed_at_seconds metrics regardless of name
		for _, m := range fam.Metric {
			var job, instance string
			for _, l := range m.Label {
				switch *l.Name {
				case "job":
					job = *l.Value
				case "instance":
					instance = *l.Value
				}
			}

			// Skip if job doesn't match prefix
			if *jobPrefix != "" && !strings.HasPrefix(job, *jobPrefix) {
				continue
			}

			if m.Gauge != nil && m.Gauge.Value != nil {
				result[[2]string{job, instance}] = *m.Gauge.Value
			}
		}
	}
	return result, nil
}

// deleteJobInstance removes metrics for a given job/instance
func deleteJobInstance(pushgwURL, job, instance string) error {
	u, err := url.Parse(pushgwURL)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/metrics/job/%s", url.PathEscape(job))
	if instance != "" {
		path += fmt.Sprintf("/instance/%s", url.PathEscape(instance))
	}
	u.Path = path

	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}

// runCleaner executes a single cleaning cycle
func runCleaner() {
	metrics, err := parsePushgatewayMetrics(*pushgw)
	if err != nil {
		logger.Errorf("Failed to fetch metrics: %v", err)
		errorsTotal.Inc()
		return
	}

	now := time.Now().Unix()
	for key, ts := range metrics {
		job, instance := key[0], key[1]
		checkedTotal.Inc()

		if age := now - int64(ts); age > int64(*ttl) {
			instDisplay := instance
			if instDisplay == "" {
				instDisplay = "<empty>"
			}

			logger.Infof("Deleting job=%s instance=%s (age=%ds)", job, instDisplay, age)
			if err := deleteJobInstance(*pushgw, job, instance); err != nil {
				logger.Errorf("Delete failed: %v", err)
				errorsTotal.Inc()
			} else {
				deletedTotal.Inc()
			}
		}
	}
}

func main() {
	flag.Parse()

	if *daemon {
		// Start metrics endpoint
		http.Handle("/metrics", promhttp.Handler())
		go func() {
			addr := fmt.Sprintf(":%d", *metricsPort)
			logger.Infof("Cleaner metrics listening on %s/metrics", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				logger.Fatalf("Metrics server failed: %v", err)
			}
		}()

		// Run periodically
		ticker := time.NewTicker(time.Duration(*interval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			runCleaner()
		}
	} else {
		runCleaner()
	}
}
