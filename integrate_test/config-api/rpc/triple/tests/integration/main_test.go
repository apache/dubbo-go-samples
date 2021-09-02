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
)

import (
	"github.com/apache/dubbo-go-samples/api"
)

var greeterProvider = new(api.GreeterClientImpl)

func TestMain(m *testing.M) {
	config.SetConsumerService(greeterProvider)

	referenceConfig := config.NewReferenceConfig(
		config.WithReferenceInterface("com.apache.dubbo.sample.basic.IGreeter"),
		config.WithReferenceProtocolName("tri"),
		config.WithReferenceRegistry("zkRegistryKey"),
	)

	consumerConfig := config.NewConsumerConfig(
		config.WithConsumerReferenceConfig("greeterImpl", referenceConfig),
	)

	registryConfig := config.NewRegistryConfigWithProtocolDefaultPort("zookeeper")

	rootConfig := config.NewRootConfig(
		config.WithRootRegistryConfig("zkRegistryKey", registryConfig),
		config.WithRootConsumerConfig(consumerConfig),
	)

	if err := rootConfig.Init(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())

}
