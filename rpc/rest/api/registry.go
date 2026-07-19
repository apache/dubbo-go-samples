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

package api

import (
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

const (
	ServerAppName = "dubbo_rest_basic_server"
	ClientAppName = "dubbo_rest_basic_client"

	RegistryDirect    = "direct"
	RegistryZookeeper = "zookeeper"
	RegistryNacos     = "nacos"

	RegistryTypeInterface = "interface"
	RegistryTypeService   = "service"
	RegistryTypeAll       = "all"

	DefaultRegistry     = RegistryNacos
	DefaultRegistryType = RegistryTypeService
)

func RegistryOptions(name string, registryType string) ([]registry.Option, error) {
	if name == RegistryDirect || name == "" {
		return nil, nil
	}

	var opts []registry.Option
	switch name {
	case RegistryZookeeper:
		opts = append(opts, registry.WithZookeeper(), registry.WithAddress("127.0.0.1:2181"))
	case RegistryNacos:
		opts = append(opts, registry.WithNacos(), registry.WithAddress("127.0.0.1:8848"))
	default:
		return nil, logger.Errorf("unsupported registry %q, use direct, zookeeper, or nacos", name)
	}

	opts = append(opts, registry.WithoutUseAsConfigCenter())
	switch registryType {
	case "", RegistryTypeInterface:
		opts = append(opts, registry.WithRegisterInterface())
	case RegistryTypeService:
		opts = append(opts, registry.WithRegisterService())
	case RegistryTypeAll:
		opts = append(opts, registry.WithRegisterServiceAndInterface())
	default:
		return nil, logger.Errorf("unsupported registry type %q, use interface, service, or all", registryType)
	}
	return opts, nil
}
