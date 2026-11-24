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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type jaegerTrace struct {
	Spans []struct {
		ProcessID string `json:"processID"`
		Tags      []struct {
			Key   string `json:"key"`
			Value any    `json:"value"`
		} `json:"tags"`
	} `json:"spans"`
	Processes map[string]struct {
		ServiceName string `json:"serviceName"`
	} `json:"processes"`
}

// jaegerHasCrossServiceTrace looks for one trace in which both applications
// recorded a span for rpcService, i.e. context propagated across the call.
func jaegerHasCrossServiceTrace(base, consumer, provider, rpcService string) error {
	params := url.Values{"service": {consumer}, "limit": {"20"}, "lookback": {"1h"}}
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Get(base + "/api/traces?" + params.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result struct {
		Data []jaegerTrace `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	for _, trace := range result.Data {
		seen := map[string]bool{}
		for _, span := range trace.Spans {
			for _, tag := range span.Tags {
				if tag.Key == "rpc.service" && tag.Value == rpcService {
					seen[trace.Processes[span.ProcessID].ServiceName] = true
				}
			}
		}
		if seen[consumer] && seen[provider] {
			return nil
		}
	}
	return fmt.Errorf("no trace in Jaeger has %s spans from both %s and %s (%d traces checked)",
		rpcService, consumer, provider, len(result.Data))
}
