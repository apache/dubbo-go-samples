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
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"

	"github.com/stretchr/testify/assert"
)

import (
	dubbo3pb "github.com/apache/dubbo-go-samples/api"
)

func TestStreamSayHello(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, tripleConstant.TripleCtxKey(tripleConstant.TripleRequestID), "triple-request-id-demo")
	req := dubbo3pb.HelloRequest{
		Name: "laurence",
	}

	r, err := greeterProvider.SayHelloStream(ctx)
	assert.Nil(t, err)

	for i := 0; i < 2; i++ {
		err = r.Send(&req)
		assert.Nil(t, err)
	}

	rspUser := &dubbo3pb.User{}
	err = r.RecvMsg(rspUser)
	assert.Nil(t, err)
	assert.Equal(t, "hello laurence", rspUser.Name)
	assert.Equal(t, "123456789", rspUser.Id)
	assert.Equal(t, int32(18), rspUser.Age)

	err = r.Send(&req)
	assert.Nil(t, err)

	err = r.RecvMsg(rspUser)
	assert.Nil(t, err)
	assert.Equal(t, "hello laurence", rspUser.Name)
	assert.Equal(t, "123456789", rspUser.Id)
	assert.Equal(t, int32(19), rspUser.Age)
}

func TestSayHello(t *testing.T) {
	req := &dubbo3pb.HelloRequest{
		Name: "laurence",
	}

	ctx := context.Background()

	reply, err := greeterProvider.SayHello(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, "Hello laurence", reply.Name)
	assert.Equal(t, "12345", reply.Id)
	assert.Equal(t, int32(21), reply.Age)
}
