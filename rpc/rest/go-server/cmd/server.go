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
	"strconv"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/filter/echo"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/mapping/metadata"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/report/nacos"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/report/zookeeper"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/rest"
	_ "dubbo.apache.org/dubbo-go/v3/proxy/proxy_factory"
	_ "dubbo.apache.org/dubbo-go/v3/registry/nacos"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/servicediscovery"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
	dubboserver "dubbo.apache.org/dubbo-go/v3/server"

	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-samples/rpc/rest/api"
)

type GreetingProvider struct {
	api.GreetingService
}

func (p *GreetingProvider) GetGreeting(_ context.Context, args []any) (*api.GreetingResponse, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf("expected 4 REST arguments, got %d", len(args))
	}

	userID, err := toInt(args[0])
	if err != nil {
		return nil, fmt.Errorf("parse userID: %w", err)
	}

	name := fmt.Sprint(args[1])
	traceID := fmt.Sprint(args[2])
	message := bodyMessage(args[3])

	return &api.GreetingResponse{
		UserID:   userID,
		Name:     name,
		TraceID:  traceID,
		Message:  message,
		Greeting: fmt.Sprintf("hello %s, userID=%d, traceID=%s, message=%s", name, userID, traceID, message),
	}, nil
}

func toInt(v any) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case int32:
		return int(val), nil
	case int64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	default:
		return strconv.Atoi(fmt.Sprint(val))
	}
}

func bodyMessage(v any) string {
	switch body := v.(type) {
	case map[string]any:
		return fmt.Sprint(body["message"])
	case map[string]string:
		return body["message"]
	case *api.GreetingRequest:
		return body.Message
	case api.GreetingRequest:
		return body.Message
	default:
		return fmt.Sprint(v)
	}
}

func main() {
	registryName := flag.String("registry", api.DefaultRegistry, "registry mode: direct, zookeeper, or nacos")
	registryType := flag.String("registry-type", api.DefaultRegistryType, "registry type: interface, service, or all")
	flag.Parse()

	api.InstallProviderRestConfig()

	instanceOpts := []dubbo.InstanceOption{
		dubbo.WithName(api.ServerAppName),
		dubbo.WithProtocol(
			protocol.WithREST(),
			protocol.WithIp("127.0.0.1"),
			protocol.WithPort(20080),
		),
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

	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}

	if err := srv.RegisterService(
		&GreetingProvider{},
		dubboserver.WithInterface(api.InterfaceName),
		dubboserver.WithFilter("echo"),
	); err != nil {
		panic(err)
	}

	logger.Infof("REST provider is listening on http://127.0.0.1:20080%s, registry=%s, registry-type=%s",
		api.RestPath, *registryName, *registryType)
	if err := srv.Serve(); err != nil {
		panic(err)
	}
}
