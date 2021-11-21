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

	"google.golang.org/grpc"
)

import (
	grpcpb "github.com/apache/dubbo-go-samples/rpc/triple/pb/dubbogo-grpc/protobuf/grpc"
)

func TestGRPCClientHello(t *testing.T) {
	// Set up a connection to the client.
	conn, err := grpc.Dial("127.0.0.1:20000", grpc.WithInsecure())
	assert.Nil(t, err)
	defer conn.Close()
	c := grpcpb.NewGreeterClient(conn)

	req := &grpcpb.HelloRequest{
		Name: "laurence",
	}
	ctx := context.Background()
	rsp, err := c.SayHello(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "Hello laurence", rsp.Name)
	assert.Equal(t, "12345", rsp.Id)
	assert.Equal(t, int32(21), rsp.Age)
}

func TestGRPCClientStreamSayHello(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:20000", grpc.WithInsecure())
	assert.Nil(t, err)
	defer conn.Close()
	c := grpcpb.NewGreeterClient(conn)

	req := &grpcpb.HelloRequest{
		Name: "grpc laurence",
	}
	clientStream, err := c.SayHelloStream(context.Background())
	assert.Nil(t, err)
	for i := 0; i < 2; i++ {
		err = clientStream.Send(req)
		assert.Nil(t, err)
	}

	rspUser := &grpcpb.User{}
	err = clientStream.RecvMsg(rspUser)
	assert.Nil(t, err)
	assert.Equal(t, "hello grpc laurence", rspUser.Name)
	assert.Equal(t, "123456789", rspUser.Id)
	assert.Equal(t, int32(18), rspUser.Age)

	err = clientStream.Send(req)
	assert.Nil(t, err)

	err = clientStream.RecvMsg(rspUser)
	assert.Nil(t, err)
	assert.Equal(t, "hello grpc laurence", rspUser.Name)
	assert.Equal(t, "123456789", rspUser.Id)
	assert.Equal(t, int32(19), rspUser.Age)
}
