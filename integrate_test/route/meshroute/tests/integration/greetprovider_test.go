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
	"testing"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/stretchr/testify/assert"
)

import (
	api "github.com/apache/dubbo-go-samples/api"
)

const MeshRouteConf = "apiVersion: service.dubbo.apache.org/v1alpha1\n" +
	"kind: DestinationRule\n" +
	"metadata: { name: demo-route }\n" +
	"spec:\n" +
	"  host: demo\n" +
	"  subsets:\n" +
	"    - labels: { env-sign: xxx, tag1: hello }\n" +
	"      name: isolation\n" +
	"    - labels: { env-sign: yyy }\n" +
	"      name: testing-trunk\n" +
	"    - labels: { env-sign: zzz }\n" +
	"      name: testing\n" +
	"  trafficPolicy:\n" +
	"    loadBalancer: { simple: ROUND_ROBIN }\n\n" +
	"---\n\n" +
	"apiVersion: service.dubbo.apache.org/v1alpha1\n" +
	"kind: VirtualService\n" +
	"metadata: {name: demo-route}\n" +
	"spec:\n" +
	"  dubbo:\n" +
	"    - routedetail:\n" +
	"        - match:\n" +
	"            - sourceLabels: {trafficLabel: xxx}\n" +
	"          name: xxx-project\n" +
	"          route:\n" +
	"            - destination: {host: demo, subset: isolation}\n" +
	"        - match:\n" +
	"            - sourceLabels: {trafficLabel: testing-trunk}\n" +
	"          name: testing-trunk\n" +
	"          route:\n" +
	"            - destination: {host: demo, subset: testing-trunk}\n" +
	"        - name: testing\n" +
	"          route:\n" +
	"            - destination: {host: demo, subset: testing}\n" +
	"      services:\n" +
	"        - {exact: com.apache.dubbo.sample.basic.IGreeter}\n" +
	"  hosts: [demo]"

func TestGetUser(t *testing.T) {
	config.GetRootConfig().ConfigCenter = config.NewConfigCenterConfigBuilder().SetProtocol("zookeeper").SetAddress("127.0.0.1:2181").SetDataID("dubbo-go-samples-configcenter-zookeeper-client").Build()
	dynamicConfiguration, err := config.GetRootConfig().ConfigCenter.GetDynamicConfiguration()
	assert.Nil(t, err)
	// 1. publish mesh route config
	err = dynamicConfiguration.PublishConfig("dubbo.io.MESHAPPRULE", "dubbo", MeshRouteConf)
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
