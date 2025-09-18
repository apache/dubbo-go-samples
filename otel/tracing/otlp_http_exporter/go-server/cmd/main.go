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
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/common"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/otel/trace"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"

	"github.com/dubbogo/gost/log/logger"

	"github.com/golang/protobuf/proto"

	collecttracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

import (
	greet "github.com/apache/dubbo-go-samples/otel/tracing/stdout/proto"
)

type GreetMultiRPCServer struct {
}

func (srv *GreetMultiRPCServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	return resp, nil
}

type GreetProvider struct {
}

func (*GreetProvider) SayHello(req string, req1 string, req2 string) (string, error) {
	return req + " " + req1 + " " + req2, nil
}

var (
	// triple + dubbo + jsonrpc
	serverReceivesChan = make(chan bool, 3)
	clientReceivesChan = make(chan bool, 3)
	errChan            = make(chan error, 6)
)

func mockOtlpReceiver() {
	http.HandleFunc("/v1/traces", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var reader io.Reader = r.Body
		if strings.EqualFold(r.Header.Get("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				errChan <- err
				return
			}
			defer gr.Close()
			reader = gr
		}
		raw, err := io.ReadAll(reader)
		if err != nil {
			errChan <- err
			return
		}

		var req collecttracepb.ExportTraceServiceRequest
		if err = proto.Unmarshal(raw, &req); err != nil {
			errChan <- err
			return
		}

		reqStr := req.String()
		if strings.Contains(reqStr, "dubbo_otel_server") {
			serverReceivesChan <- true
		} else if strings.Contains(reqStr, "dubbo_otel_client") {
			clientReceivesChan <- true
		} else {
			errChan <- errors.New("unknown trace" + reqStr)
		}
	})

	err := http.ListenAndServe("127.0.0.1:4318", nil)

	if err != nil {
		panic(err)
	}
}

func main() {
	go mockOtlpReceiver()
	go func() {
		var (
			serverCount = 0
			clientCount = 0
		)
		for i := 0; i < 6; i++ {
			select {
			case <-serverReceivesChan:
				serverCount++
			case <-clientReceivesChan:
				clientCount++
			case err := <-errChan:
				panic(err)
			case <-time.After(7 * time.Second):
				panic("timeout")
			}
		}

		logger.Infof("server count: %d, client count: %d", serverCount, clientCount)
		if serverCount != 3 || clientCount != 3 {
			panic("trace received count not match")
		}
	}()

	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_otel_server"),
		dubbo.WithTracing(
			trace.WithEnabled(),
			trace.WithOtlpHttpExporter(),
			trace.WithW3cPropagator(),
			trace.WithAlwaysMode(),
			trace.WithEndpoint("127.0.0.1:4318"),
			trace.WithInsecure(),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20000)),
		dubbo.WithProtocol(
			protocol.WithDubbo(),
			protocol.WithPort(20001)),
		dubbo.WithProtocol(
			protocol.WithJSONRPC(),
			protocol.WithPort(20002)),
	)
	if err != nil {
		panic(err)
	}

	//Triple
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	if err = greet.RegisterGreetServiceHandler(srv, &GreetMultiRPCServer{}); err != nil {
		panic(err)
	}

	//Dubbo & JsonRPC
	if err = srv.Register(&GreetProvider{}, &common.ServiceInfo{}, server.WithInterface("GreetProvider")); err != nil {
		panic(err)
	}
	if err = srv.Serve(); err != nil {
		logger.Error(err)
	}
}
