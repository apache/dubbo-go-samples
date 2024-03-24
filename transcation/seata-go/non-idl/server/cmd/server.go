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
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"

	"github.com/dubbogo/gost/log/logger"

	_ "github.com/seata/seata-go/pkg/imports"
	"github.com/seata/seata-go/pkg/integration"
	"github.com/seata/seata-go/pkg/rm/tcc"

	"github.com/apache/dubbo-go-samples/transcation/seata-go/non-idl/server/service"
)

func main() {
	integration.UseDubbo()
	userProviderProxy, err := tcc.NewTCCServiceProxy(&service.UserProvider{})
	if err != nil {
		logger.Errorf("get userProviderProxy tcc service proxy error, %v", err.Error())
		return
	}
	err = userProviderProxy.RegisterResource()
	if err != nil {
		logger.Errorf("userProviderProxy register resource error, %v", err.Error())
	}
	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_seata_server"),
	)
	if err != nil {
		panic(err)
	}
	srv, err := ins.NewServer(
		server.WithServerProtocol(
			protocol.WithDubbo(),
			protocol.WithPort(20000),
		),
	)
	if err != nil {
		panic(err)
	}
	if err := srv.Register(userProviderProxy, nil, server.WithInterface("UserProvider"), server.WithSerialization(constant.Hessian2Serialization)); err != nil {
		panic(err)
	}
	if err := srv.Serve(); err != nil {
		panic(err)
	}
}
