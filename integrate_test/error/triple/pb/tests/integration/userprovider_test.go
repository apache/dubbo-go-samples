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
	"fmt"
	"testing"
)

import (
	"github.com/dubbogo/triple/pkg/common"

	"github.com/stretchr/testify/assert"
)

import (
	triplepb "github.com/apache/dubbo-go-samples/api"
)

func TestSayHello(t *testing.T) {
	ctx := context.Background()
	req := triplepb.HelloRequest{
		Name: "laurence",
	}
	_, err := greeterProvider.SayHello(ctx, &req)
	stacks, ok := err.(common.TripleError)
	assert.True(t, ok)
	assert.Equal(t, "user defined error", fmt.Sprintf("%s", err))
	assert.NotEqual(t, "", stacks.Stacks())
	assert.Equal(t, uint32(1234), uint32(stacks.Code()))
	assert.Equal(t, "user defined error", stacks.Message())
}
