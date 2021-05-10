// +build integration

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
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	"github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

import (
	"os"
	"testing"
	"time"
)

var appName = "UserConsumerTest"
var referenceConfig = config.ReferenceConfig{
	InterfaceName: "org.apache.dubbo.UserProvider",
	Cluster:       "failover",
	Registry:      "demoZk",
	Protocol:      dubbo.DUBBO,
	Generic:       true,
}

func TestMain(m *testing.M) {
	referenceConfig.GenericLoad(appName)
	time.Sleep(3 * time.Second)

	os.Exit(m.Run())
}

type User struct {
	ID   string
	Name string
	Age  int32
	Time time.Time
}

func (User) JavaClassName() string {
	return "org.apache.dubbo.User"
}
