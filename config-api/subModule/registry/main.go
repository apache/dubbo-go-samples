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
	"strconv"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

func main() {
	registryConfig := config.NewRegistryConfigWithProtocolDefaultPort("zookeeper")
	reg, err := registryConfig.GetInstance(common.PROVIDER)
	if err != nil {
		panic(err)
	}

	ivkURL, err := common.NewURL("mock://localhost:8080",
		common.WithPath("com.alibaba.dubbogo.HelloService"),
		common.WithParamsValue(constant.RoleKey, strconv.Itoa(common.PROVIDER)),
		common.WithMethods([]string{"GetUser", "SayHello"}),
	)
	if err != nil {
		panic(err)
	}
	if err := reg.Register(ivkURL); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 30)
	if err := reg.UnRegister(ivkURL); err != nil {
		panic(err)
	}
}
