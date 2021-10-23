// +build integration

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
	"os"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	routeFilePath := os.Getenv("PROJECT_PATH")
	routeBytes, err := ioutil.ReadFile(routeFilePath + "/go-client/conf/mesh_route.yml")
	assert.Nil(t, err)
	dynamicConfiguration, err := config.GetRootConfig().ConfigCenter.GetDynamicConfiguration()
	assert.Nil(t, err)
	// 1. publish mesh route config
	err = dynamicConfiguration.PublishConfig("dubbo.io.MESHAPPRULE", "dubbo", string(routeBytes))
	assert.Nil(t, err)

	req := &api.HelloRequest{
		Name: "Dong",
	}
	reply, err := grpcGreeterImpl.SayHello(context.Background(), req)
	assert.Nil(t, err)

	assert.Equal(t, "Hello Dong", reply.Name)
	assert.Equal(t, "12345", reply.Id)
	assert.Equal(t, 21, reply.Age)
}
