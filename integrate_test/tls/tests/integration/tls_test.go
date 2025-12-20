// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
)

import (
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/tls"
)

import (
	greet "github.com/apache/dubbo-go-samples/tls/proto"
)

// buildCertPath returns an absolute path to a file inside tls/x509.
func buildCertPath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	base := filepath.Dir(file)                 // .../integrate_test/tls/tests/integration
	root := filepath.Join(base, "../../../..") // repo root
	return filepath.Join(root, "tls", "x509", name)
}

func TestTLSGreet(t *testing.T) {
	cli, err := client.NewClient(
		client.WithClientURL("127.0.0.1:20000"),
		client.WithClientProtocolTriple(),
		client.WithClientTLSOption(
			tls.WithCACertFile(buildCertPath("server_ca_cert.pem")),
			tls.WithServerName("dubbogo.test.example.com"),
		),
	)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello tls"})
	if err != nil {
		t.Fatalf("greet call failed: %v", err)
	}

	if resp == nil || resp.Greeting != "hello tls" {
		t.Fatalf("unexpected response: %#v", resp)
	}
}
