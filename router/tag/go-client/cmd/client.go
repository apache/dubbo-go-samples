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

	callGreet := func(name, tag, force string) {
		// set tag attachments for invocation
		atta := map[string]string{
			constant.Tagkey:      tag,
			constant.ForceUseTag: force,
		}
		ctx := context.WithValue(context.Background(), constant.AttachmentKey, atta)

		resp, err := svc.Greet(ctx, &greet.GreetRequest{Name: name})
		printRes(resp, err)
	}

	callGreet("tag with force", "test-tag", "true")      // success
	callGreet("tag with force", "test-tag1", "true")     // fail
	callGreet("tag with no-force", "test-tag1", "false") // success
	callGreet("non-tag", "", "false")                    // success
}

func printRes(rep *greet.GreetResponse, err error) {
	if err != nil {
		logger.Errorf("❌ invoke failed: %v", err)
	} else {
		logger.Infof("✔ invoke successfully : %v", rep.Greeting)
	}
}
