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
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log"
)

import (
	pkg2 "github.com/apache/dubbo-go-samples/context/dubbo/go-client/pkg"
)

var userProvider = new(pkg2.UserProvider)

func init() {
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&pkg2.ContextContent{})
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	hessian.RegisterPOJO(&pkg2.ContextContent{})
	config.Load()
	time.Sleep(3 * time.Second)

	gxlog.CInfo("\n\n\nstart to test dubbo")
	rspContent := &pkg2.ContextContent{}
	atta := make(map[string]interface{})
	atta["string-value"] = "string-demo"
	atta["int-value"] = 1231242
	atta["user-defined-value"] = pkg2.ContextContent{InterfaceName: "test.interface.name"}
	reqContext := context.WithValue(context.Background(), constant.DubboCtxKey("attachment"), atta)
	err := userProvider.GetContext(reqContext, []interface{}{"A001"}, rspContent)
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("response result: %+v\n", rspContent)
}
