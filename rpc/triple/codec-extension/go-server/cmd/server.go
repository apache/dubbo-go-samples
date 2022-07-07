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
	"github.com/dubbogo/gost/log/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	_ "github.com/apache/dubbo-go-samples/rpc/triple/codec-extension/codec"
)

type User struct {
	ID   string
	Name string
	Age  int32
}

type UserProvider struct {
}

func (u *UserProvider) GetUser(ctx context.Context, req *User, req2 *User, name string) (*User, error) {
	logger.Infof("req:%#v", req)
	logger.Infof("req2:%#v", req2)
	logger.Infof("name%#v", name)
	rsp := User{"12345", req.Name + req2.Name, 18}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

// export DUBBO_GO_CONFIG_PATH=PATH_TO_SAMPLES/rpc/triple/codec-extension/go-server/conf/dubbogo.yml
func main() {
	config.SetProviderService(&UserProvider{})
	if err := config.Load(); err != nil {
		panic(err)
	}
	select {}
}
