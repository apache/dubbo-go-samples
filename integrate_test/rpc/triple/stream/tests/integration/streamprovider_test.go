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
	greet "github.com/apache/dubbo-go-samples/rpc/triple/stream/proto"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestSayHello(t *testing.T) {
	req := &greet.GreetRequest{Name: "hello world"}

	ctx := context.Background()

	reply, err := greetService.Greet(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, "hello world", reply.Greeting)
}

func TestUnary(t *testing.T) {
	resp, err := greetService.Greet(context.Background(), &greet.GreetRequest{Name: "triple"})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestBidiStream(t *testing.T) {
	stream, err := greetService.GreetStream(context.Background())
	assert.Nil(t, err)

	sendErr := stream.Send(&greet.GreetStreamRequest{Name: "triple"})
	assert.Nil(t, sendErr)

	resp, err := stream.Recv()
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	err = stream.CloseRequest()
	assert.Nil(t, err)

	err = stream.CloseResponse()
	assert.Nil(t, err)
}

func TestClientStream(t *testing.T) {
	stream, err := greetService.GreetClientStream(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, stream)

	for i := 0; i < 5; i++ {
		sendErr := stream.Send(&greet.GreetClientStreamRequest{Name: "triple"})
		assert.Nil(t, sendErr)
	}

	resp, err := stream.CloseAndRecv()
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestServerStream(t *testing.T) {
	stream, err := greetService.GreetServerStream(context.Background(), &greet.GreetServerStreamRequest{Name: "triple"})
	assert.Nil(t, err)
	assert.NotNil(t, stream)

	for stream.Recv() {
		assert.NotNil(t, stream.Msg())
	}
	assert.Nil(t, stream.Err())

	err = stream.Close()
	assert.Nil(t, err)
}
