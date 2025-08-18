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

// Pushgateway Config (flags)
var (
	pushGatewayURL = getEnv("PUSHGATEWAY_URL", "127.0.0.1:9091")
	jobName        = getEnv("JOB_NAME", "push")
	usePush        = flag.Bool("push", true, "use push mode")
)

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
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
			registry.WithAddress("127.0.0.1:2181"),
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
			metrics.WithPrometheusGatewayUsername("username"),
			metrics.WithPrometheusGatewayPassword("1234"),
			metrics.WithPrometheusGatewayUrl(pushGatewayURL),
			metrics.WithPrometheusGatewayInterval(10*time.Second),
			metrics.WithPrometheusGatewayJob(jobName),
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

	// business loop
	for i := 0; i < 10; i++ {
		resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
		if err != nil {
			logger.Error(err)
		} else {
			logger.Infof("Greet response: %s", resp.Greeting)
		}

		pushedAt.Set(float64(time.Now().Unix()))

		time.Sleep(1 * time.Second)
	}

	// cleanup: delete job from Pushgateway if push mode enabled
	if *usePush {
		path := fmt.Sprintf("http://%s", pushGatewayURL)
		if err := deletePushgatewayJob(path, jobName); err != nil {
			logger.Errorf("Delete pushgateway job failed: %v", err)
		} else {
			logger.Infof("Deleted job from Pushgateway: job=%s", jobName)
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
	req.SetBasicAuth("username", "1234")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("unexpected status: %s", resp.Status)
}
