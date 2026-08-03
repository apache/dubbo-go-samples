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

package tracefields

import (
	"context"
	"strings"
	"testing"
)

import (
	"go.opentelemetry.io/otel/trace"
)

func TestFieldsWithoutSpan(t *testing.T) {
	if got := Fields(context.Background()); got != " trace_id=- span_id=-" {
		t.Fatalf("Fields() without span = %q, want empty identifiers", got)
	}
}

func TestFieldsWithSpan(t *testing.T) {
	traceID := trace.TraceID{1, 2, 3}
	spanID := trace.SpanID{4, 5, 6}
	ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	}))

	got := Fields(ctx)
	if !strings.Contains(got, "trace_id="+traceID.String()) {
		t.Fatalf("Fields() = %q, missing trace ID", got)
	}
	if !strings.Contains(got, "span_id="+spanID.String()) {
		t.Fatalf("Fields() = %q, missing span ID", got)
	}
}
