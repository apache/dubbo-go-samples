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
	"github.com/apache/dubbo-go-samples/compatibility/rpc/triple/pb2/models"
)

func TestStream(t *testing.T) {

	ctx := context.Background()
	ctx = context.WithValue(ctx, tripleConstant.TripleCtxKey("tri-req-id"), "triple-request-id-demo")

	req := models.HelloRequest{
		Name: "dubbo-go",
	}

	r, err := greeterProvider.SayHelloStream(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, r)

	for i := 0; i < 2; i++ {
		err = r.Send(&req)
		assert.Nil(t, err)
	}

	rspUser := &models.User{}
	err = r.RecvMsg(rspUser)
	assert.Nil(t, err)
	assert.NotNil(t, rspUser)

	err = r.Send(&req)
	assert.Nil(t, err)

	rspUser2 := &models.User{}
	err = r.RecvMsg(rspUser2)
	assert.Nil(t, err)
	assert.NotNil(t, rspUser2)

}

func TestUnary(t *testing.T) {

	req := models.HelloRequest{
		Name: "dubbo-go",
	}
	user, err := greeterProvider.SayHello(context.Background(), &req)
	assert.Nil(t, err)
	assert.NotNil(t, user)

}
