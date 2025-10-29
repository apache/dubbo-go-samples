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
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/constant"

	"github.com/stretchr/testify/assert"
)

import (
	greet "github.com/apache/dubbo-go-samples/context/proto"
)

func TestSayHello(t *testing.T) {
	req := &greet.GreetRequest{Name: "hello world"}

	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.AttachmentKey, map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	})

	reply, err := greeterProvider.Greet(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, "name: hello world, key1: value1, key2: value2", reply.Greeting)
}
