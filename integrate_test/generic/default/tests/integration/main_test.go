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
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"

	hessian "github.com/apache/dubbo-go-hessian2"

	tpconst "github.com/dubbogo/triple/pkg/common/constant"
)

import (
	"github.com/apache/dubbo-go-samples/generic/default/go-server/pkg"
)

const appName = "dubbo.io"

var (
	dubboRefConf  config.ReferenceConfig
	tripleRefConf config.ReferenceConfig
)

func TestMain(m *testing.M) {
	hessian.RegisterPOJO(&pkg.User{})

	dubboRefConf = newRefConf("org.apache.dubbo.UserProvider", dubbo.DUBBO)
	tripleRefConf = newRefConf("org.apache.dubbo.UserProviderTriple", tpconst.TRIPLE)

	os.Exit(m.Run())
}

func newRefConf(iface, protocol string) config.ReferenceConfig {
	registryConfig := &config.RegistryConfig{
		Protocol: "zookeeper",
		Address:  "127.0.0.1:2181",
	}

	refConf := config.ReferenceConfig{
		InterfaceName: iface,
		Cluster:       "failover",
		RegistryIDs:   []string{"zk"},
		Protocol:      protocol,
		Generic:       "true",
	}

	rootConfig := config.NewRootConfigBuilder().
		AddRegistry("zk", registryConfig).
		Build()

	if err := config.Load(config.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}
	_ = refConf.Init(rootConfig)
	refConf.GenericLoad(appName)

	return refConf
}
