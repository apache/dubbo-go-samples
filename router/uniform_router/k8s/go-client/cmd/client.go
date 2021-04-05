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
	"github.com/apache/dubbo-go-samples/router/condition/go-client/pkg"
	"github.com/dubbogo/gost/log"
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/cluster/router/uniform"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

var userProvider = new(pkg.UserProvider)

func init() {
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&pkg.User{})
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
// need to setup environment variable "CONF_ROUTER_FILE_PATH" to "conf/router_config.yml" before run
func main() {
	hessian.RegisterPOJO(&pkg.User{})
	config.Load()

	gxlog.CInfo("\n\n\nstart to test dubbo")
	user := &pkg.User{}
	for {
		err := userProvider.GetUser(context.TODO(), []interface{}{"A001"}, user)
		if err != nil {
			gxlog.CError("error: %v\n", err)
			os.Exit(1)
			return
		}
		gxlog.CInfo("response result: %v\n", user)
		time.Sleep(time.Second)
	}

}
