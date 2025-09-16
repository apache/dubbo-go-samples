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
	triplepb "github.com/apache/dubbo-go-samples/compatibility/api"
)

func TestSayHello(t *testing.T) {
	ctx := context.Background()
	// set user defined context attachment, map value can be string or []string, otherwise it is not to be transferred
	userDefinedValueMap := make(map[string]any)
	userDefinedValueMap["key1"] = "user defined value 1"
	userDefinedValueMap["key2"] = "user defined value 2"
	userDefinedValueMap["key3"] = []string{"user defined value 3.1", "user defined value 3.2"}
	userDefinedValueMap["key4"] = []string{"user defined value 4.1", "user defined value 4.2"}
	ctx = context.WithValue(ctx, constant.AttachmentKey, userDefinedValueMap)

	req := triplepb.HelloRequest{
		Name: "laurence",
	}

	user, err := greeterProvider.SayHello(ctx, &req)

	assert.Nil(t, err)
	assert.Equal(t, "map[key1:[user defined value 1] key2:[user defined value 2] key3:[user defined value 3.1 user defined value 3.2] key4:[user defined value 4.1 user defined value 4.2]]", user.Name)
	assert.Equal(t, "12345", user.Id)
	assert.Equal(t, int32(21), user.Age)
}
