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
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log"
)

type UserProvider struct {
	GetContext func(ctx context.Context, req *ContextContent) (rsp *ContextContent, err error)
}

func (u *UserProvider) Reference() string {
	return "userProvider"
}

type ContextContent struct {
	Path              string
	InterfaceName     string
	DubboVersion      string
	LocalAddr         string
	RemoteAddr        string
	UserDefinedStrVal string
	CtxStrVal         string
	CtxIntVal         int64
}

func (*ContextContent) JavaClassName() string {
	return "org.apache.dubbo.User"
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	var userProvider = &UserProvider{}
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&ContextContent{})
	config.Load()
	gxlog.CInfo("\n\n\nstart to test dubbo")

	atta := make(map[string]interface{})
	atta["string-value"] = "string-demo"
	atta["int-value"] = 1231242
	atta["user-defined-value"] = ContextContent{InterfaceName: "test.interface.name"}
	reqContext := context.WithValue(context.Background(), constant.DubboCtxKey("attachment"), atta)
	rspContent, err := userProvider.GetContext(reqContext, &ContextContent{CtxStrVal: "A001"})
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("response result: %+v\n", rspContent)
}
