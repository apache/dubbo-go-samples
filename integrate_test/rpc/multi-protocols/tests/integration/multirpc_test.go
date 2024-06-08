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
	greet "github.com/apache/dubbo-go-samples/rpc/multi-protocols/proto"
	"github.com/dubbogo/gost/log/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSayHello(t *testing.T) {
	//Triple
	req := &greet.GreetRequest{Name: "hello world"}

	ctx := context.Background()

	reply, err := greeterProvider.Greet(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, "hello world", reply.Greeting)
	logger.Infof("GreetProvider.Greet reply: %s", reply.Greeting)

	//Dubbo
	var respDubbo string
	if err = connDubbo.CallUnary(context.Background(), []interface{}{"hello", "new", "dubbo"}, &respDubbo, "SayHello"); err != nil {
		logger.Errorf("GreetProvider.Greet err: %s", err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, "hellonewdubbo", respDubbo)

	//JsonRpc
	var respJsonRpc string
	if err = connJsonRpc.CallUnary(context.Background(), []interface{}{"hello", "new", "jsonrpc"}, &respJsonRpc, "SayHello"); err != nil {
		logger.Errorf("GreetProvider.Greet err: %s", err)
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, "hellonewjsonrpc", respJsonRpc)
}
