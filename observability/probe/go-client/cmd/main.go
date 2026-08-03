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
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/direct/proto"
)

const (
	tripleAddr   = "127.0.0.1:20000"
	probeBaseURL = "http://127.0.0.1:22222"

	tcpReadyTimeout   = 20 * time.Second
	probeReadyTimeout = 40 * time.Second
	requestTimeout    = 2 * time.Second
)

func main() {
	ctx := context.Background()

	if err := waitTCPReady(tripleAddr, tcpReadyTimeout); err != nil {
		logger.Fatalf("triple port is not ready: %v", err)
		panic("triple port is not ready")
	}
	logger.Infof("triple port is ready: %s", tripleAddr)

	if err := waitTCPReady("127.0.0.1:22222", tcpReadyTimeout); err != nil {
		logger.Fatalf("probe port is not ready: %v", err)
		panic("probe port is not ready")
	}
	logger.Info("probe port is ready: 127.0.0.1:22222")

	if err := expectHTTPStatus(probeBaseURL+"/live", http.StatusOK, requestTimeout); err != nil {
		logger.Fatalf("/live check failed: %v", err)
		panic("/live check failed")
	}
	logger.Info("/live is healthy")

	if err := waitHTTPStatus(probeBaseURL+"/ready", http.StatusOK, probeReadyTimeout, requestTimeout); err != nil {
		logger.Fatalf("/ready did not become healthy: %v", err)
		panic("/ready did not become healthy")
	}
	logger.Info("/ready is healthy")

	if err := waitHTTPStatus(probeBaseURL+"/startup", http.StatusOK, probeReadyTimeout, requestTimeout); err != nil {
		logger.Fatalf("/startup did not become healthy: %v", err)
		panic("/startup did not become healthy")
	}
	logger.Info("/startup is healthy")

	if err := callGreet(ctx); err != nil {
		logger.Fatalf("greet rpc check failed: %v", err)
		panic("greet rpc check failed")
	}
	logger.Info("probe sample integration checks passed")
}

func callGreet(ctx context.Context) error {
	cli, err := client.NewClient(
		client.WithClientURL("tri://" + tripleAddr),
	)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		return fmt.Errorf("create greet service: %w", err)
	}

	rpcCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := svc.Greet(rpcCtx, &greet.GreetRequest{Name: "probe-check"})
	if err != nil {
		return fmt.Errorf("invoke greet: %w", err)
	}
	if strings.TrimSpace(resp.Greeting) == "" {
		return fmt.Errorf("empty greet response")
	}
	logger.Infof("greet rpc succeeded: %s", resp.Greeting)
	return nil
}

func waitTCPReady(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("tcp %s not ready within %s: %w", addr, timeout, err)
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func waitHTTPStatus(url string, expected int, timeout, reqTimeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for {
		err := expectHTTPStatus(url, expected, reqTimeout)
		if err == nil {
			return nil
		}
		lastErr = err
		if time.Now().After(deadline) {
			return fmt.Errorf("url %s not ready within %s: %w", url, timeout, lastErr)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func expectHTTPStatus(url string, expected int, reqTimeout time.Duration) error {
	client := http.Client{Timeout: reqTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != expected {
		return fmt.Errorf("expect status %d but got %d, body=%s", expected, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}
