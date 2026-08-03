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
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/metrics"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"

	"github.com/prometheus/client_golang/prometheus"
)

import (
	greet "github.com/apache/dubbo-go-samples/helloworld/proto"
)

// Config structure for application settings
type Config struct {
	PushGatewayURL  string // PushGateway URL for metrics
	PushGatewayUser string // Username for PushGateway authentication
	PushGatewayPass string // Password for PushGateway authentication
	JobName         string // Job name for PushGateway
	ZkAddress       string // ZooKeeper address for service registry
	UsePush         bool   // Flag to enable/disable push mode
}

var config Config

// Initialize configuration from environment variables and flags
func init() {
	flag.BoolVar(&config.UsePush, "push", true, "use push mode")
	config.PushGatewayURL = getEnv("PUSHGATEWAY_URL", "127.0.0.1:9091")
	config.PushGatewayUser = getEnv("PUSHGATEWAY_USER", "username")
	config.PushGatewayPass = getEnv("PUSHGATEWAY_PASS", "1234")
	config.JobName = getEnv("JOB_NAME", "dubbo_client")
	config.ZkAddress = getEnv("ZK_ADDRESS", "127.0.0.1:2181")
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func main() {
	flag.Parse()

	// Prometheus metric
	pushedAt := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "job_pushed_at_seconds",
		Help: "Unix seconds of last push to Pushgateway",
	})
	prometheus.MustRegister(pushedAt)

	// Dubbo instance
	ins, err := dubbo.NewInstance(
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress(config.ZkAddress),
		),
		dubbo.WithMetrics(
			metrics.WithEnabled(),
			metrics.WithPrometheus(),
			metrics.WithPrometheusExporterEnabled(),
			metrics.WithPort(9097),
			metrics.WithPath("/prometheus"),
			metrics.WithMetadataEnabled(),
			metrics.WithRegistryEnabled(),
			metrics.WithConfigCenterEnabled(),

			metrics.WithPrometheusPushgatewayEnabled(),
			metrics.WithPrometheusGatewayUsername(config.PushGatewayUser),
			metrics.WithPrometheusGatewayPassword(config.PushGatewayPass),
			metrics.WithPrometheusGatewayUrl(config.PushGatewayURL),
			metrics.WithPrometheusGatewayInterval(10*time.Second),
			metrics.WithPrometheusGatewayJob(config.JobName),
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

	// Handle graceful shutdown
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	// Dead loop, can be stopped by signal
	go func() {
		for {
			resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
			if err != nil {
				logger.Error(err)
			} else {
				logger.Infof("Greet response: %s", resp.Greeting)
			}

			pushedAt.Set(float64(time.Now().Unix()))
			time.Sleep(1 * time.Second)
		}
	}()

	// Wait for signal
	<-stopCh
	logger.Info("Received shutdown signal, cleaning up...")

	// Cleanup: delete job from Pushgateway if push mode enabled
	if config.UsePush {
		path := fmt.Sprintf("http://%s", config.PushGatewayURL)
		if err := deletePushgatewayJob(path, config.JobName); err != nil {
			logger.Errorf("Delete pushgateway job failed: %v", err)
		} else {
			logger.Infof("Deleted job from Pushgateway: job=%s", config.JobName)
		}
	}
}

// delete Pushgateway job
func deletePushgatewayJob(pushgw, job string) error {
	u, err := url.Parse(pushgw)
	if err != nil {
		return err
	}
	u.Path = fmt.Sprintf("/metrics/job/%s", url.PathEscape(job))
	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(config.PushGatewayUser, config.PushGatewayPass)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("unexpected status: %s", resp.Status)
}
