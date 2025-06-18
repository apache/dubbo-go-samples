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
	reflection "dubbo.apache.org/dubbo-go/v3/protocol/triple/reflection/triple_reflection"

	greet "github.com/apache/dubbo-go-samples/rpc/triple/reflection/proto"
)

import (
	"github.com/golang/protobuf/proto"

	"github.com/stretchr/testify/assert"

	"google.golang.org/protobuf/types/descriptorpb"
)

func TestSayHello(t *testing.T) {
	req := &greet.GreetRequest{Name: "hello world"}

	ctx := context.Background()

	reply, err := greetService.Greet(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, "hello world", reply.Greeting)
}

func TestFileByFilename(t *testing.T) {
	err := stream.Send(&reflection.ServerReflectionRequest{
		MessageRequest: &reflection.ServerReflectionRequest_FileByFilename{FileByFilename: "reflection.proto"},
	})
	assert.Nil(t, err)
	recv, err := stream.Recv()
	assert.Nil(t, err)
	assert.NotNil(t, recv)
	m := new(descriptorpb.FileDescriptorProto)
	err = proto.Unmarshal(recv.GetFileDescriptorResponse().GetFileDescriptorProto()[0], m)
	assert.Nil(t, err)
}

func TestFileContainingSymbol(t *testing.T) {
	err := stream.Send(&reflection.ServerReflectionRequest{
		MessageRequest: &reflection.ServerReflectionRequest_FileContainingSymbol{FileContainingSymbol: "dubbo.reflection.v1alpha.ServerReflection"},
	})
	assert.Nil(t, err)
	recv, err := stream.Recv()
	assert.Nil(t, err)
	assert.NotNil(t, recv)
	m := new(descriptorpb.FileDescriptorProto)
	err = proto.Unmarshal(recv.GetFileDescriptorResponse().GetFileDescriptorProto()[0], m)
	assert.Nil(t, err)
}

func TestListServices(t *testing.T) {
	err := stream.Send(&reflection.ServerReflectionRequest{
		MessageRequest: &reflection.ServerReflectionRequest_ListServices{},
	})
	assert.Nil(t, err)
	recv, err := stream.Recv()
	assert.Nil(t, err)
	assert.NotNil(t, recv)
}
