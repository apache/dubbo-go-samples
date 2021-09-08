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
	"os"
	"testing"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"

	hessian "github.com/apache/dubbo-go-hessian2"
)

import (
	pkg "github.com/apache/dubbo-go-samples/generic/default/go-server/pkg"
)

var (
	appName         = "dubbo.io"
	referenceConfig config.ReferenceConfig
)

func init() {
	registryConfig := &config.RegistryConfig{
		Protocol: "zookeeper",
		Address:  "127.0.0.1:2181",
	}

	referenceConfig = config.ReferenceConfig{
		InterfaceName: "org.apache.dubbo.UserProvider",
		Cluster:       "failover",
		Registry:      []string{"zk"},
		Protocol:      dubbo.DUBBO,
		Generic:       "true",
	}

	rootConfig := config.NewRootConfig(config.WithRootRegistryConfig("zk", registryConfig))
	_ = rootConfig.Init()
	_ = referenceConfig.Init(rootConfig)
	referenceConfig.GenericLoad(appName)
}

func TestMain(m *testing.M) {
	hessian.RegisterPOJO(&pkg.User{})
	config.SetProviderService(&pkg.User{})
	config.Load()

	os.Exit(m.Run())
}
