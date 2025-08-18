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
	"os"
	"strings"
	"time"

	"github.com/dubbogo/gost/log/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
)

// Flags
var (
	pushgw      = flag.String("pushgw", "http://127.0.0.1:9091", "Pushgateway address")
	ttl         = flag.Int("ttl", 3600, "TTL seconds for job/instance (default 300s)")
	daemon      = flag.Bool("daemon", false, "Run in daemon mode")
	interval    = flag.Int("interval", 60, "Interval seconds for daemon mode")
	jobPrefix   = flag.String("job-prefix", "job_", "Only clean jobs with this prefix")
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

// parsePushgatewayMetrics fetch metrics from Pushgateway and extract job/instance
func parsePushgatewayMetrics(pushgw string) (map[[2]string]float64, error) {
	resp, err := http.Get(pushgw + "/metrics")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	parser := expfmt.TextParser{}
	mf, err := parser.TextToMetricFamilies(io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return nil, err
	}

	result := make(map[[2]string]float64)
	if fam, ok := mf["job_pushed_at_seconds"]; ok {
		for _, m := range fam.Metric {
			var job, instance string
			for _, l := range m.Label {
				if *l.Name == "job" {
					job = *l.Value
				}
				if *l.Name == "instance" {
					instance = *l.Value
				}
			}
			if m.Gauge != nil {
				result[[2]string{job, instance}] = *m.Gauge.Value
			}
		}
	}
	return result, nil
}

// deleteJobInstance remove metrics for given job/instance
func deleteJobInstance(pushgw, job, instance string) error {
	u, err := url.Parse(pushgw)
	if err != nil {
		return err
	}
	if instance == "" {
		u.Path = fmt.Sprintf("/metrics/job/%s", url.PathEscape(job))
	} else {
		u.Path = fmt.Sprintf("/metrics/job/%s/instance/%s",
			url.PathEscape(job), url.PathEscape(instance))
	}
	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("unexpected status: %s", resp.Status)
}

// runCleaner executes a single cleaning cycle against Pushgateway.
func runCleaner() {
	metrics, err := parsePushgatewayMetrics(*pushgw)
	if err != nil {
		logger.Errorf("failed to fetch metrics: %v", err)
		errorsTotal.Inc()
		return
	}

	now := time.Now().Unix()
	for key, ts := range metrics {
		job, instance := key[0], key[1]
		checkedTotal.Inc()
		age := now - int64(ts)

		if *jobPrefix != "" && !strings.HasPrefix(job, *jobPrefix) {
			continue
		}
		if age > int64(*ttl) {
			inst := instance
			if inst == "" {
				inst = "<empty>"
			}
			logger.Infof("Deleting job=%s instance=%s (age=%ds)", job, inst, age)
			if err := deleteJobInstance(*pushgw, job, instance); err != nil {
				logger.Errorf("delete failed: %v", err)
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
		// start metrics endpoint
		http.Handle("/metrics", promhttp.Handler())
		go func() {
			addr := fmt.Sprintf(":%d", *metricsPort)
			logger.Infof("Cleaner metrics listening on %s/metrics", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				logger.Fatalf("metrics server failed: %v", err)
				os.Exit(1)
			}
		}()

		// run periodically
		ticker := time.NewTicker(time.Duration(*interval) * time.Second)
		for range ticker.C {
			runCleaner()
		}
	} else {
		runCleaner()
	}
}
