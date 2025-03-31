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
	"dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
	"net/http"
	"testing"

	greet "github.com/apache/dubbo-go-samples/context/proto"
	"github.com/stretchr/testify/assert"
)

func TestSayHello(t *testing.T) {
	req := &greet.GreetRequest{Name: "hello world"}

	header := http.Header{"testKey1": []string{"testVal1"}, "testKey2": []string{"testVal2"}}
	// to store outgoing data ,and reserve the location for the receive field.
	// header will be copy , and header's key will change to be lowwer.
	ctx := triple_protocol.NewOutgoingContext(context.Background(), header)
	ctx = triple_protocol.AppendToOutgoingContext(ctx, "testKey3", "testVal3")

	reply, err := greeterProvider.Greet(ctx, req)

	extractedHeader, _ := triple_protocol.FromIncomingContext(ctx)
	var serverValue1, serverValue2 string
	if values, ok := extractedHeader["outgoingcontextkey1"]; ok && len(values) > 0 {
		serverValue1 = values[0]
	}
	if values, ok := extractedHeader["outgoingcontextkey2"]; ok && len(values) > 0 {
		serverValue2 = values[0]
	}

	assert.Nil(t, err)
	assert.Equal(t, "name: hello world, testKey1: testVal1, testKey2: testVal2", reply.Greeting)

	assert.Equal(t, "OutgoingDataVal1", serverValue1)
	assert.Equal(t, "OutgoingDataVal2", serverValue2)
}
