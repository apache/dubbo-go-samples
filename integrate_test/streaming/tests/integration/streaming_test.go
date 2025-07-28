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
	"github.com/stretchr/testify/assert"
)

import (
	greet "github.com/apache/dubbo-go-samples/streaming/proto"
)

func TestStreaming(t *testing.T) {
	t.Run("test TRIPLE unary call", testUnary)
	t.Run("test TRIPLE bidi stream", testBidiStream)
	t.Run("test TRIPLE client stream", testClientStream)
	t.Run("test TRIPLE server stream", testServerStream)
}

func testUnary(t *testing.T) {
	resp, err := greeterProvider.Greet(context.Background(), &greet.GreetRequest{Name: "triple"})
	assert.Nil(t, err)
	assert.Equal(t, "triple", resp.Greeting)
}

func testBidiStream(t *testing.T) {
	stream, err := greeterProvider.GreetStream(context.Background())
	assert.Nil(t, err)
	assert.Nil(t, stream.Send(&greet.GreetStreamRequest{Name: "triple"}))
	resp, err := stream.Recv()
	assert.Nil(t, err)
	assert.Equal(t, "triple", resp.Greeting)
	assert.Nil(t, stream.CloseRequest())
	assert.Nil(t, stream.CloseResponse())
}

func testClientStream(t *testing.T) {
	stream, err := greeterProvider.GreetClientStream(context.Background())
	assert.Nil(t, err)
	for i := 0; i < 5; i++ {
		assert.Nil(t, stream.Send(&greet.GreetClientStreamRequest{Name: "triple"}))
	}
	resp, err := stream.CloseAndRecv()
	assert.Nil(t, err)
	assert.Equal(t, "triple,triple,triple,triple,triple", resp.Greeting)
}

func testServerStream(t *testing.T) {
	stream, err := greeterProvider.GreetServerStream(context.Background(), &greet.GreetServerStreamRequest{Name: "triple"})
	assert.Nil(t, err)
	for stream.Recv() {
		assert.Equal(t, "triple", stream.Msg().Greeting)
	}
	assert.Nil(t, stream.Err())
	assert.Nil(t, stream.Close())
}
