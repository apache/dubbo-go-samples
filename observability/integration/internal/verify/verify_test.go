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

package verify

import (
	"errors"
	"strings"
	"testing"
)

import (
	"dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
)

import (
	observability "github.com/apache/dubbo-go-samples/observability/integration/proto"
)

func TestGreetRequestExpectedSuccess(t *testing.T) {
	resp := &observability.GreetResponse{Greeting: "hello alice"}
	if err := GreetRequestExpected("alice", resp, nil); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestGreetRequestExpectedNilResponse(t *testing.T) {
	if err := GreetRequestExpected("alice", nil, nil); err == nil {
		t.Fatal("nil response must be rejected")
	}
}

// Mutation: the provider returns an empty or wrong greeting.
func TestGreetRequestExpectedWrongGreeting(t *testing.T) {
	for _, greeting := range []string{"", "hello bob", "HELLO alice"} {
		resp := &observability.GreetResponse{Greeting: greeting}
		if err := GreetRequestExpected("alice", resp, nil); err == nil {
			t.Fatalf("greeting %q must be rejected for request %q", greeting, "alice")
		}
	}
}

// Mutation: the provider fails an otherwise successful request.
func TestGreetRequestExpectedUnexpectedError(t *testing.T) {
	err := errors.New("boom")
	got := GreetRequestExpected("alice", nil, err)
	if got == nil {
		t.Fatal("unexpected RPC error must be rejected")
	}
	if !strings.Contains(got.Error(), "boom") {
		t.Fatalf("rejection must preserve the original cause, got %v", got)
	}
}

func TestGreetRequestExpectedForcedError(t *testing.T) {
	err := triple_protocol.NewError(triple_protocol.CodeBizError, errors.New("forced error"))
	if got := GreetRequestExpected("error", nil, err); got != nil {
		t.Fatalf("expected typed forced error to be accepted, got %v", got)
	}
}

// Mutation: the provider error branch is removed and the forced request
// succeeds like a normal request.
func TestGreetRequestExpectedForcedErrorSucceeded(t *testing.T) {
	resp := &observability.GreetResponse{Greeting: "hello error"}
	if err := GreetRequestExpected("error", resp, nil); err == nil {
		t.Fatal("a succeeded forced-error request must be rejected")
	}
}

// Mutation: the provider returns an untyped or wrong-code error.
func TestGreetRequestExpectedForcedErrorWrongCode(t *testing.T) {
	err := triple_protocol.NewError(triple_protocol.CodeUnknown, errors.New("forced error"))
	if got := GreetRequestExpected("error", nil, err); got == nil {
		t.Fatal("forced error with a wrong code must be rejected")
	}
}
