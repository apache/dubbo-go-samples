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

// Package verify validates that a Greet call outcome matches the business
// contract of the observability integration sample. The validation is
// deterministic so that a provider regression - a wrong or empty greeting,
// a removed error branch, or a forced-error request that silently succeeds -
// fails the client run instead of being accepted.
package verify

import (
	"errors"
	"fmt"
)

import (
	"dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
)

import (
	observability "github.com/apache/dubbo-go-samples/observability/integration/proto"
)

// GreetRequestExpected reports how the observed Greet outcome deviates from
// the sample contract. A request with name "error" must fail with a typed
// CodeBizError business error; any other name must succeed and return the
// exact greeting "hello <name>". It returns nil when the outcome matches.
func GreetRequestExpected(name string, resp *observability.GreetResponse, callErr error) error {
	if name == "error" {
		return forcedErrorExpected(callErr)
	}
	return successExpected(name, resp, callErr)
}

func successExpected(name string, resp *observability.GreetResponse, callErr error) error {
	if callErr != nil {
		return fmt.Errorf("unexpected error for request %q: %w", name, callErr)
	}
	if resp == nil {
		return fmt.Errorf("nil response for request %q", name)
	}
	want := "hello " + name
	if got := resp.GetGreeting(); got != want {
		return fmt.Errorf("unexpected greeting for request %q: got %q, want %q", name, got, want)
	}
	return nil
}

func forcedErrorExpected(callErr error) error {
	if callErr == nil {
		return errors.New(`request with name "error" succeeded; expected the forced business error`)
	}
	if code := triple_protocol.CodeOf(callErr); code != triple_protocol.CodeBizError {
		return fmt.Errorf("forced error request reported code %s (%v), want %s", code, callErr, triple_protocol.CodeBizError)
	}
	return nil
}
