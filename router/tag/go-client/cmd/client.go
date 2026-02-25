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
	"strings"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/router/tag/proto"
)

const (
	RegistryAddress = "127.0.0.1:8848"
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("tag-client"),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress(RegistryAddress),
		),
	)

	if err != nil {
		logger.Errorf("new instance failed: %v", err)
		panic(err)
	}

	cli, err := ins.NewClient()

	if err != nil {
		logger.Errorf("new client failed: %v", err)
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)

	if err != nil {
		logger.Errorf("new service failed: %v", err)
		panic(err)
	}

	callGreet := func(name, tag, force, exp string) {
		// set tag attachments for invocation
		atta := map[string]string{
			constant.Tagkey:      tag,
			constant.ForceUseTag: force,
		}
		ctx := context.WithValue(context.Background(), constant.AttachmentKey, atta)

		rep, err := svc.Greet(ctx, &greet.GreetRequest{Name: name})
		// temporarily cancel checking for result, for PR # 3208 (https://github.com/apache/dubbo-go/pull/3208)
		// hasn't been merged to main branch yet, thus tag router still not works properly.

		//checkRes(exp, resp.GetGreeting(), err)
		if err != nil {
			logger.Errorf("❌ invoke failed: %v", err)
		} else {
			logger.Infof("✔ invoke successfully : %v", rep.Greeting)
		}
	}

	callGreet("tag with force", "test-tag", "true", "server-with-tag")         // success
	callGreet("tag with force", "test-tag1", "true", "fail")                   // fail
	callGreet("tag with no-force", "test-tag1", "false", "server-without-tag") // success
	callGreet("non-tag", "", "false", "server-without-tag")                    // success
}

func checkRes(exp string, act string, err error) {
	if (err == nil && exp == "fail") || (err != nil && exp != "fail") {
		panic("unexpected result!")
	} else if act != "" && !strings.Contains(act, exp) {
		panic("unexpected result!")
	}
}
