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
	"fmt"
	"time"
)

import (
	"github.com/apache/dubbo-go-samples/chain/frontend/pkg"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/service/inmemory"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
)

func main() {
	var chinese = new(pkg.ChineseService)
	var american = new(pkg.AmericanService)
	config.SetConsumerService(chinese)
	config.SetConsumerService(american)
	config.Load()
	time.Sleep(3 * time.Second)
	have, _ := chinese.Have()
	fmt.Printf("chinese.Have(): %s\n", have)
	hear, _ := chinese.Hear()
	fmt.Printf("chinese.Hear(): %s\n", hear)
	have, _ = american.Have()
	fmt.Printf("american.Have(): %s\n", have)
	hear, _ = american.Hear()
	fmt.Printf("american.Hear(): %s\n", hear)
}
