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
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"

	hessian "github.com/apache/dubbo-go-hessian2"

	tpconst "github.com/dubbogo/triple/pkg/common/constant"
)

import (
	"github.com/apache/dubbo-go-samples/generic/default/go-client/pkg"
)

const appName = "dubbo.io"

// export DUBBO_GO_CONFIG_PATH= PATH_TO_SAMPLES/generic/default/go-client/conf/dubbogo.yml
func main() {
	// register POJOs
	hessian.RegisterPOJO(&pkg.User{})

	// generic invocation samples using hessian serialization on Dubbo protocol
	dubboRefConf := newRefConf("org.apache.dubbo.UserProvider", dubbo.DUBBO)
	callGetUser(dubboRefConf)

	// generic invocation samples using hessian serialization on Triple protocol
	tripleRefConf := newRefConf("org.apache.dubbo.UserProviderTriple", tpconst.TRIPLE)
	callGetUser(tripleRefConf)
}

func callGetUser(refConf config.ReferenceConfig) {
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser",
		[]string{"java.lang.String", "java.lang.String"},
		[]hessian.Object{"A003", "lily"},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("GetUser(userId string, name string) res: %+v", resp)
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
