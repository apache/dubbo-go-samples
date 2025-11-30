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
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/require"
)

import (
	greet "github.com/apache/dubbo-go-samples/filter/proto"
)

func TestCustomFilter_AttachmentsFlow(t *testing.T) {
	drainAttachmentChannel(requestAttachmentCh)
	drainAttachmentChannel(responseAttachmentCh)

	req := &greet.GreetRequest{Name: "integration-custom-filter"}
	resp, err := greetService.Greet(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, req.Name, resp.Greeting)

	reqAttachment := waitAttachment(t, requestAttachmentCh)
	requireAttachmentValues(t, reqAttachment, "request-key1", []string{"request-value1"})
	requireAttachmentValues(t, reqAttachment, "request-key2", []string{"request-value2.1", "request-value2.2"})

	respAttachment := waitAttachment(t, responseAttachmentCh)
	requireAttachmentValues(t, respAttachment, "key1", []string{"value1"})
	requireAttachmentValues(t, respAttachment, "key2", []string{"value1", "value2"})
}

func waitAttachment(t *testing.T, ch chan map[string]interface{}) map[string]interface{} {
	t.Helper()
	select {
	case attachment := <-ch:
		return attachment
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for attachment")
		return nil
	}
}

func drainAttachmentChannel(ch chan map[string]interface{}) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}

func requireAttachmentValues(t *testing.T, attachments map[string]interface{}, key string, expected []string) {
	t.Helper()
	raw, ok := attachments[key]
	require.Truef(t, ok, "attachment key %s not found", key)

	var actual []string
	switch v := raw.(type) {
	case string:
		actual = []string{v}
	case []string:
		actual = v
	default:
		t.Fatalf("unexpected attachment type for key %s: %T", key, raw)
	}

	require.ElementsMatch(t, expected, actual)
}
