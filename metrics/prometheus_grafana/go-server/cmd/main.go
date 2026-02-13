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
	"math/rand"
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
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"dubbo.apache.org/dubbo-go/v3/server"

	"github.com/dubbogo/gost/log/logger"

	"github.com/pkg/errors"

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
	config.JobName = getEnv("JOB_NAME", "dubbo_server")
	config.ZkAddress = getEnv("ZK_ADDRESS", "127.0.0.1:2181")
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

type GreetTripleServer struct{}

func (srv *GreetTripleServer) Greet(_ context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	if r.Intn(101) > 99 { // mock error here
		return nil, errors.New("random error")
	}
	time.Sleep(10 * time.Millisecond) // mock business delay
	return resp, nil
}

func main() {
	flag.Parse()

	pushedAt := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "job_pushed_at_seconds",
		Help: "Unix seconds of last push to Pushgateway",
	})
	prometheus.MustRegister(pushedAt)

	ins, err := dubbo.NewInstance(
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress(config.ZkAddress),
		),
		dubbo.WithMetrics(
			metrics.WithEnabled(),
			metrics.WithPrometheus(),
			metrics.WithPrometheusExporterEnabled(),
			metrics.WithPort(9099),
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

			metrics.WithAggregationEnabled(),
			metrics.WithAggregationTimeWindowSeconds(30),
			metrics.WithAggregationBucketNum(10),
		),
	)
	if err != nil {
		panic(err)
	}

	srv, err := ins.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000),
			protocol.WithTriple(),
		),
	)
	if err != nil {
		panic(err)
	}

	if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
		panic(err)
	}

	// Custom pushedAt metric goroutine
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				pushedAt.Set(float64(time.Now().Unix()))
			}
		}
	}()

	// Handle graceful shutdown
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	// Run server
	go func() {
		logger.Info("Starting server...")
		if err := srv.Serve(); err != nil {
			logger.Error("Server error:", err)
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
