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
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	dubbo3pb "github.com/apache/dubbo-go-samples/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasProvider(t *testing.T) {
	req := &dubbo3pb.HelloRequest{
		Name: "laurence",
	}

	ctx := context.Background()
	atm := map[string]string{
		"dubbo.tag":       "tag1",
		"dubbo.force.tag": "true",
	}
	ctx = context.WithValue(ctx, constant.AttachmentKey, atm)

	reply, err := greeterProvider.SayHello(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, "Hello laurence", reply.Name)
	assert.Equal(t, "2000", reply.Id)
	assert.Equal(t, int32(21), reply.Age)
}

func TestNoProvider(t *testing.T) {
	req := &dubbo3pb.HelloRequest{
		Name: "laurence",
	}

	ctx := context.Background()
	atm := map[string]string{
		"dubbo.tag":       "tag2",
		"dubbo.force.tag": "true",
	}
	ctx = context.WithValue(ctx, constant.AttachmentKey, atm)

	reply, err := greeterProvider.SayHello(ctx, req)

	assert.Nil(t, reply)
	assert.NotNil(t, err)
}
