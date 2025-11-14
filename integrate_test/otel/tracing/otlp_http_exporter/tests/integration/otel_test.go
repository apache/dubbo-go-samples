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

package integration

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/otel/trace"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greetpb "github.com/apache/dubbo-go-samples/otel/tracing/stdout/proto"
)

// ---- Test configuration ----
const (
	defaultOtelEndpoint = "127.0.0.1:4318"

	tripleAddr  = "127.0.0.1:20000"
	dubboAddr   = "127.0.0.1:20001"
	jsonRPCAddr = "127.0.0.1:20002"

	testTimeout   = 8 * time.Second
	exporterGrace = 5 * time.Second // best-effort to allow OTLP export
)

var ins *dubbo.Instance

// TestMain handles global setup/teardown once for all integration tests.
func TestMain(m *testing.M) {
	var err error

	otel := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otel == "" {
		otel = defaultOtelEndpoint
	}

	ins, err = dubbo.NewInstance(
		dubbo.WithName("dubbo_otel_client_it"),
		dubbo.WithTracing(
			trace.WithEnabled(),
			trace.WithOtlpHttpExporter(), // OTLP/HTTP -> /v1/traces
			trace.WithW3cPropagator(),
			trace.WithAlwaysMode(),
			trace.WithEndpoint(otel),
			trace.WithInsecure(),
		),
	)
	if err != nil {
		fmt.Printf("init dubbo instance failed: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	// Give exporter a short window to flush spans (best-effort).
	time.Sleep(exporterGrace)
	os.Exit(code)
}

func waitPort(addr string, deadline time.Duration) error {
	end := time.Now().Add(deadline)
	for {
		c, err := net.DialTimeout("tcp", addr, 1*time.Second)
		if err == nil {
			c.Close()
			return nil
		}
		if time.Now().After(end) {
			return fmt.Errorf("port %s not ready: %w", addr, err)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func TestMultiProtocols(t *testing.T) {

	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	must(waitPort("127.0.0.1:20000", 10*time.Second)) // triple
	must(waitPort("127.0.0.1:20001", 10*time.Second)) // dubbo(getty)
	must(waitPort("127.0.0.1:20002", 10*time.Second)) // jsonrpc

	t.Run("triple", testTriple)
	t.Run("dubbo", testDubbo)
	t.Run("jsonrpc", testJSONRPC)
}

// ----- Subtests -----

func testTriple(t *testing.T) {
	t.Helper()

	cli, err := ins.NewClient(
		client.WithClientProtocolTriple(),
		client.WithClientURL(tripleAddr),
	)
	if err != nil {
		t.Fatalf("create triple client: %v", err)
	}

	svc, err := greetpb.NewGreetService(cli)
	if err != nil {
		t.Fatalf("create greet service: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	resp, err := svc.Greet(ctx, &greetpb.GreetRequest{Name: "hello world"})
	if err != nil {
		t.Fatalf("triple call failed: %v", err)
	}

	if resp == nil || resp.Greeting == "" {
		t.Fatalf("unexpected triple response: %#v", resp)
	}
	logger.Infof("[triple] Greet response: %s", resp.Greeting)

	time.Sleep(exporterGrace) // allow spans to export
}

func testDubbo(t *testing.T) {
	t.Helper()

	cli, err := ins.NewClient(
		client.WithClientProtocolDubbo(),
		client.WithClientSerialization(constant.Hessian2Serialization),
		client.WithClientURL(dubboAddr),
	)
	if err != nil {
		t.Fatalf("create dubbo client: %v", err)
	}

	conn, err := cli.Dial("GreetProvider")
	if err != nil {
		t.Fatalf("dial dubbo service: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var resp string
	if err = conn.CallUnary(ctx, []any{"hello", "new", "dubbo"}, &resp, "SayHello"); err != nil {
		t.Fatalf("dubbo call failed: %v", err)
	}
	if resp == "" {
		t.Fatalf("empty dubbo response")
	}
	logger.Infof("[dubbo] SayHello response: %s", resp)

	time.Sleep(exporterGrace)
}

func testJSONRPC(t *testing.T) {
	t.Helper()

	cli, err := ins.NewClient(
		client.WithClientProtocolJsonRPC(),
		client.WithClientSerialization(constant.JSONSerialization),
		client.WithClientURL(jsonRPCAddr),
	)
	if err != nil {
		t.Fatalf("create jsonrpc client: %v", err)
	}

	conn, err := cli.Dial("GreetProvider")
	if err != nil {
		t.Fatalf("dial jsonrpc service: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var resp string
	if err = conn.CallUnary(ctx, []any{"hello", "new", "jsonrpc"}, &resp, "SayHello"); err != nil {
		t.Fatalf("jsonrpc call failed: %v", err)
	}
	if resp == "" {
		t.Fatalf("empty jsonrpc response")
	}
	logger.Infof("[jsonrpc] SayHello response: %s", resp)

	time.Sleep(exporterGrace)
}
