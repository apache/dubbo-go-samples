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
	"strings"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-samples/compatibility/api"
)

type GreeterProvider struct {
	api.UnimplementedGreeterServer
}

func (s *GreeterProvider) SayHello(ctx context.Context, in *api.HelloRequest) (*api.User, error) {
	// map must be assert to map[string]interface, because of dubbo limitation
	attachments := ctx.Value(constant.AttachmentKey).(map[string]interface{})

	// value must be assert to []string[0], because of http2 header limitation
	logger.Infof("get triple attachment key1 = %s", attachments["key1"].([]string)[0])
	logger.Infof("get triple attachment key2 = %s", attachments["key2"].([]string)[0])
	logger.Infof("get triple attachment key3 = %s and %s", attachments["key3"].([]string)[0],
		attachments["key3"].([]string)[1])
	logger.Infof("get triple attachment key4 = %s and %s", attachments["key4"].([]string)[0],
		attachments["key4"].([]string)[1])
	logger.Infof("Dubbo3 GreeterProvider get user name = %s\n", in.Name)
	rspAttachment := make(map[string]interface{})
	for k, v := range attachments {
		if strings.HasPrefix(k, "key") {
			rspAttachment[k] = v
		}
	}
	return &api.User{Name: fmt.Sprintf("%s", rspAttachment), Id: "12345", Age: 21}, nil
}

func (s *GreeterProvider) SayHelloStream(svr api.Greeter_SayHelloStreamServer) error {
	// map must be assert to map[string]interface, because of dubbo limitation
	attachments := svr.Context().Value(constant.AttachmentKey).(map[string]interface{})

	// value must be assert to []string[0], because of http2 header limitation
	logger.Infof("get triple attachment key1 = %s", attachments["key1"].([]string)[0])
	logger.Infof("get triple attachment key2 = %s", attachments["key2"].([]string)[0])
	logger.Infof("get triple attachment key3 = %s and %s", attachments["key3"].([]string)[0],
		attachments["key3"].([]string)[1])
	logger.Infof("get triple attachment key4 = %s and %s", attachments["key4"].([]string)[0],
		attachments["key4"].([]string)[1])
	c, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go3 GreeterProvider recv 1 user, name = %s\n", c.Name)
	err = svr.Send(&api.User{
		Name: "hello " + c.Name,
		Age:  18,
		Id:   "123456789",
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	config.SetProviderService(&GreeterProvider{})
	if err := config.Load(); err != nil {
		panic(err)
	}
	select {}
}
