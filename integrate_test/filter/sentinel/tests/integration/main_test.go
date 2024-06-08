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
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"os"
	"testing"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/filter/proto/sentinel"
)

var greetService greet.SentinelGreetService

func TestMain(m *testing.M) {
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}

	greetService, err = greet.NewSentinelGreetService(cli, client.WithFilter(constant.SentinelConsumerFilterKey))

	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
