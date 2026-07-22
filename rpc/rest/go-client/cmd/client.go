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

package main

import (
	"context"
	"flag"
	"fmt"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster/failover"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster/zoneaware"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance/random"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/router/condition"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/filter/graceful_shutdown"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/mapping/metadata"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/report/nacos"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/report/zookeeper"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/rest"
	_ "dubbo.apache.org/dubbo-go/v3/registry/directory"
	_ "dubbo.apache.org/dubbo-go/v3/registry/nacos"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/servicediscovery"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"

	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-samples/rpc/rest/api"
)

func main() {
	registryName := flag.String("registry", api.DefaultRegistry, "registry mode: direct, zookeeper, or nacos")
	registryType := flag.String("registry-type", api.DefaultRegistryType, "registry type: interface, service, or all")
	flag.Parse()

	api.InstallConsumerRestConfig()

	instanceOpts := []dubbo.InstanceOption{
		dubbo.WithName(api.ClientAppName),
	}
	registryOpts, err := api.RegistryOptions(*registryName, *registryType)
	if err != nil {
		panic(err)
	}
	if len(registryOpts) > 0 {
		instanceOpts = append(instanceOpts, dubbo.WithRegistry(registryOpts...))
	}

	ins, err := dubbo.NewInstance(instanceOpts...)
	if err != nil {
		panic(err)
	}

	clientOpts := []client.ClientOption{
		client.WithClientNoCheck(),
		client.WithClientRequestTimeout(5 * time.Second),
	}
	if *registryName == api.RegistryDirect || *registryName == "" {
		clientOpts = append(clientOpts, client.WithClientURL("rest://127.0.0.1:20080"))
	}

	cli, err := ins.NewClient(clientOpts...)
	if err != nil {
		panic(err)
	}

	conn, err := cli.Dial(
		api.InterfaceName,
		client.WithProtocol(constant.RESTProtocol),
	)
	if err != nil {
		panic(err)
	}

	resp := new(api.GreetingResponse)
	err = conn.CallUnary(
		context.Background(),
		[]any{
			101,
			"dubbo-go",
			"trace-rest-basic",
			&api.GreetingRequest{Message: "body-from-dubbo-rest-client"},
		},
		resp,
		api.MethodGetGreeting,
	)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	if resp.UserID != 101 || resp.Name != "dubbo-go" || resp.TraceID != "trace-rest-basic" {
		panic(fmt.Errorf("gets an unexpected result: userID=%d name=%s traceID=%s message=%s greeting=%q",
			resp.UserID, resp.Name, resp.TraceID, resp.Message, resp.Greeting))
	}

	logger.Infof("REST response: userID=%d name=%s traceID=%s message=%s greeting=%q\n",
		resp.UserID, resp.Name, resp.TraceID, resp.Message, resp.Greeting)
}
