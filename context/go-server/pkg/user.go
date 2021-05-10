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

package pkg

import (
	"context"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config"
	"github.com/dubbogo/gost/log"
)

func init() {
	config.SetProviderService(new(UserProvider))
	// ------for hessian2------
	hessian.RegisterPOJO(&ContextContent{})
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

type UserProvider struct {
}

func (u *UserProvider) GetContext(ctx context.Context, req []interface{}) (*ContextContent, error) {
	gxlog.CInfo("req:%#v", req)
	ctxAtta := ctx.Value(constant.DubboCtxKey("attachment")).(map[string]interface{})
	userDefinedval := ctxAtta["user-defined-value"].(*ContextContent)
	gxlog.CInfo("get user defined struct:%#v", userDefinedval)
	rsp := ContextContent{
		Path:              ctxAtta["path"].(string),
		InterfaceName:     ctxAtta["interface"].(string),
		DubboVersion:      ctxAtta["dubbo"].(string),
		LocalAddr:         ctxAtta["local-addr"].(string),
		RemoteAddr:        ctxAtta["remote-addr"].(string),
		UserDefinedStrVal: userDefinedval.InterfaceName,
		CtxIntVal:         ctxAtta["int-value"].(int64),
		CtxStrVal:         ctxAtta["string-value"].(string),
	}
	gxlog.CInfo("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}

func (u ContextContent) JavaClassName() string {
	return "org.apache.dubbo.ContextContent"
}
