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
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	log "dubbo.apache.org/dubbo-go/v3/logger"
	"dubbo.apache.org/dubbo-go/v3/protocol"

	"github.com/dubbogo/gost/log/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20000),
		),
		dubbo.WithLogger(
			log.WithZap(),
			log.WithLevel("debug"),
			log.WithTraceIntegration(true),
			log.WithRecordErrorToSpan(true),
		),
	)
	if err != nil {
		panic(err)
	}

	server, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	go server.Serve()

	// setup OpenTelemetry tracer
	tp := trace.NewTracerProvider()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer("demo")

	// run for 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// create trace span from the timeout context
	traceCtx, span := tracer.Start(ctx, "demo-operation")
	defer span.End()

	// get CtxLogger
	rawLogger := logger.GetLogger()
	ctxLog := rawLogger.(log.CtxLogger)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// log with trace context (automatically includes trace_id, span_id)
			ctxLog.CtxInfo(traceCtx, "hello dubbogo this is info log")
			ctxLog.CtxDebug(traceCtx, "hello dubbogo this is debug log")
			ctxLog.CtxWarn(traceCtx, "hello dubbogo this is warn log")
			ctxLog.CtxError(traceCtx, "hello dubbogo this is error log")
			time.Sleep(time.Second * 1)
		}
	}
}
