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
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	user "github.com/apache/dubbo-go-samples/router/polaris/proto"
)

func main() {
	polarisAddr := "127.0.0.1:8091"

	ins, err := dubbo.NewInstance(
		dubbo.WithName("myApp"),
		dubbo.WithRegistry(
			registry.WithPolaris(),
			registry.WithAddress(polarisAddr),
			registry.WithNamespace("dubbogo"),
			registry.WithRegisterInterface(),
		),
	)
	if err != nil {
		logger.Errorf("new dubbo instance failed: %v", err)
		panic(err)
	}

	cli, err := ins.NewClient()
	if err != nil {
		logger.Errorf("new client failed: %v", err)
		panic(err)
	}

	svc, err := user.NewUserService(cli)
	if err != nil {
		logger.Errorf("create user service failed: %v", err)
		panic(err)
	}

	logger.Infof("\n\n\nstart to test dubbo")

	uid := os.Getenv("uid")

	for i := 0; i < 5; i++ {
		time.Sleep(200 * time.Millisecond)
		req := &user.User{Name: "Alex001"}
		ctx := context.WithValue(context.Background(), constant.AttachmentKey, map[string]interface{}{
			"uid": uid,
		})
		resp, err := svc.GetUser(ctx, req)
		if err != nil {
			panic(err)
		}
		logger.Infof("uid=%s, response: %v\n", uid, resp)
	}
}
