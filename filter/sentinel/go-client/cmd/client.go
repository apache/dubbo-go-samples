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
	"time"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
	"github.com/alibaba/sentinel-golang/api"
	sentinelConf "github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/dubbogo/gost/log"
)

import (
	"github.com/apache/dubbo-go-samples/filter/sentinel/go-client/pkg"
)

var userProvider = new(pkg.UserProvider)

func init() {
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&pkg.User{})
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	hessian.RegisterPOJO(&pkg.User{})
	config.Load()
	time.Sleep(3 * time.Second)

	if err := initSentinel(); err != nil {
		panic(err)
	}
	_, err := flow.LoadRules([]*flow.Rule{
		{
			// protocol:consumer:interfaceName:group:version:method
			Resource:               "dubbo:consumer:org.apache.dubbo.UserProvider:::GetUser()",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              1,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		panic(err)
	}

	gxlog.CInfo("\n\n\nstart to test sentinel filter")
	for i := 0; i <= 5; i++ {
		user := &pkg.User{}
		err = userProvider.GetUser(context.TODO(), []interface{}{"A001"}, user)
		if err != nil {
			gxlog.CError("error: %v\n", err)
			continue
		}
		gxlog.CInfo("response result: %v\n", user)
	}

}

func initSentinel() error {
	// custom changes configs
	conf := sentinelConf.NewDefaultConfig()
	return api.InitWithConfig(conf)
}
