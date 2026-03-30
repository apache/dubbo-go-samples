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
	"fmt"
	"os"
	"strings"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/cluster/router"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/global"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/direct/proto"
)

const (
	clientApplication = "static-tag-client"
	tagName           = "gray"
	grayAddress       = "127.0.0.1:20002"
	directURL         = "tri://127.0.0.1:20000;tri://127.0.0.1:20002?dubbo.tag=gray"
	expectedServer    = "server-with-gray-tag"
	attemptCount      = 5
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName(clientApplication),
		dubbo.WithRouter(
			router.WithScope("application"),
			router.WithKey("static-tag-provider"),
			router.WithPriority(100),
			router.WithForce(false),
			router.WithTags([]global.Tag{
				{
					Name:      tagName,
					Addresses: []string{grayAddress},
				},
			}),
		),
	)
	if err != nil {
		logger.Errorf("new instance failed: %v", err)
		panic(err)
	}

	cli, err := ins.NewClient(client.WithClientURL(directURL))
	if err != nil {
		logger.Errorf("new client failed: %v", err)
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		logger.Errorf("new service failed: %v", err)
		panic(err)
	}

	ctx := context.WithValue(context.Background(), constant.AttachmentKey, map[string]string{
		constant.Tagkey: tagName,
	})

	for i := 1; i <= attemptCount; i++ {
		resp, err := svc.Greet(ctx, &greet.GreetRequest{
			Name: fmt.Sprintf("static tag router attempt %d", i),
		})
		if err != nil {
			logger.Errorf("invoke failed on attempt %d/%d: %v", i, attemptCount, err)
			os.Exit(1)
		}

		if !strings.Contains(resp.Greeting, expectedServer) {
			logger.Errorf("routed to unexpected server on attempt %d/%d, want %s, got %q",
				i, attemptCount, expectedServer, resp.Greeting)
			os.Exit(1)
		}

		logger.Infof("invoke successfully on attempt %d/%d: %v", i, attemptCount, resp.Greeting)
	}
}
