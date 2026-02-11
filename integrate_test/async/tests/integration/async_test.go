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
	"github.com/stretchr/testify/assert"
)

import (
	user "github.com/apache/dubbo-go-samples/async/proto"
)

func TestAsync(t *testing.T) {
	req := &user.GetUserRequest{Id: "003"}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		_, err := userProvider.GetUser(ctx, req)
		done <- err
	}()

	select {
	case err := <-done:
		assert.Nil(t, err)
	case <-ctx.Done():
		assert.Fail(t, "async call timeout", ctx.Err().Error())
	}
}

func TestAsyncOneWay(t *testing.T) {
	req := &user.SayHelloRequest{UserId: "003"}
	_, err := userProviderV2.SayHello(context.Background(), req)
	assert.Nil(t, err)
}
