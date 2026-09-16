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
	"strings"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/tls"

	"github.com/dubbogo/gost/log/logger"
)

import (
	quic "github.com/apache/dubbo-go-samples/http3/internal/quic"
	greet "github.com/apache/dubbo-go-samples/http3/proto"
)

func main() {
	logger.SetLoggerLevel("debug")

	cli, err := client.NewClient(
		client.WithClientURL("127.0.0.1:20000"),
		client.WithClientTLSOption(
			tls.WithCACertFile("x509/server_ca_cert.pem"),
			tls.WithCertFile("x509/server2_cert.pem"),
			tls.WithKeyFile("x509/server2_key_pkcs8.pem"),
			tls.WithServerName("dubbogo.test.example.com"),
		),
		// Enable HTTP/3 support on client side
		// This configures the client to use dualTransport which supports
		// both HTTP/2 and HTTP/3 with Alt-Svc negotiation
		client.WithClientProtocol(
			protocol.WithTriple(
				quic.Options()...,
			),
		),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	// Make multiple RPC calls to demonstrate HTTP/3 upgrade process:
	// 1. First request: Uses HTTP/2 (HTTPS over TCP)
	// 2. Server responds with Alt-Svc header advertising HTTP/3 support
	// 3. Subsequent requests: Can upgrade to HTTP/3 (QUIC over UDP)
	// The client-side dualTransport automatically handles this negotiation
	for i := 0; i < 3; i += 1 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)

		resp, err := svc.Greet(ctx, &greet.GreetRequest{Name: "Go Client"})
		if err != nil {
			panic(err)
		}

		if resp == nil {
			panic("expected greeting response, got empty response")
		}

		// The greeting must echo the request name and come from an HTTP/3 backend,
		// e.g. "Hello Go Client, from Go HTTP/3!" or "..., from Java HTTP/3!".
		const prefix = "Hello Go Client, from "
		const suffix = "HTTP/3!"
		if !strings.HasPrefix(resp.Greeting, prefix) || !strings.HasSuffix(resp.Greeting, suffix) {
			panic(fmt.Sprintf("unexpected greeting %q, want prefix %q and suffix %q",
				resp.Greeting, prefix, suffix))
		}

		logger.Infof("Greet response: %s", resp.Greeting)

		cancel()
		time.Sleep(time.Second)
	}

	logger.Info("All requests completed successfully!")
}
