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
	"os"
	"time"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"

	"github.com/dubbogo/gost/log"
)

import (
	"github.com/apache/dubbo-go-samples/config-api/go-client/pkg"
)

var userProvider = new(pkg.UserProvider)

func setConfigByAPI() {
	consumerConfig := config.NewConsumerConfig(
		config.WithConsumerAppConfig(config.NewDefaultApplicationConfig()),
		config.WithConsumerConnTimeout(time.Second*3),
		config.WithConsumerRequestTimeout(time.Second*3),
		config.WithConsumerRegistryConfig("demoZk", config.NewDefaultRegistryConfig("zookeeper")),
		config.WithConsumerReferenceConfig("UserProvider", config.NewReferenceConfigByAPI(
			config.WithReferenceRegistry("demoZk"),
			config.WithReferenceProtocol("dubbo"),
			config.WithReferenceInterface("org.apache.dubbo.UserProvider"),
			config.WithReferenceMethod("GetUser", "3", "random"),
			config.WithReferenceCluster("failover"),
		)),
	)
	config.SetConsumerConfig(*consumerConfig)
}

func init() {
	setConfigByAPI()
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&pkg.User{})
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	hessian.RegisterPOJO(&pkg.User{})
	config.Load()
	time.Sleep(3 * time.Second)

	gxlog.CInfo("\n\n\nstart to test dubbo")
	user := &pkg.User{}
	err := userProvider.GetUser(context.TODO(), []interface{}{"A001"}, user)
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("response result: %v\n", user)
}
