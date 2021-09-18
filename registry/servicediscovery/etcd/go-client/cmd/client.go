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
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log"
)

import (
	"github.com/apache/dubbo-go-samples/registry/servicediscovery/etcd/go-client/pkg"
)

var userProvider = new(pkg.UserProvider)

func init() {
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&pkg.User{})
}

// export DUBBO_GO_CONFIG_PATH= $PATH_TO_SAMPLES/registry/servicediscovery/go-client/conf/dubbogo.yml
func main() {
	hessian.RegisterPOJO(&pkg.User{})
	config.Load()
	time.Sleep(8 * time.Second)

	gxlog.CInfo("\n\n\nstart to test dubbo")
	for i := 0; i < 123; i++ {
		user, err := userProvider.GetUser(context.TODO(), []interface{}{"A001"})
		if err != nil {
			gxlog.CError("error: %v\n", err)
			os.Exit(1)
			return
		}
		gxlog.CInfo("response result: %v\n", user)
		time.Sleep(1 * time.Second)
	}
}
