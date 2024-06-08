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

	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"

	_ "github.com/seata/seata-go/pkg/imports"
	"github.com/seata/seata-go/pkg/integration"
	"github.com/seata/seata-go/pkg/rm/tcc"
	"github.com/seata/seata-go/pkg/tm"

	SeataClient "github.com/apache/dubbo-go-samples/transcation/seata-go/non-idl/client/seata-client"
	"github.com/apache/dubbo-go-samples/transcation/seata-go/non-idl/client/service"
)

// need to setup environment variable "DUBBO_GO_CONFIG_PATH" to "seata-go/tcc/client/conf/dubbogo.yml"
// and run "seata-go/tcc/server/cmd/server.go" before run
func main() {
	integration.UseDubbo()
	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_seata_client"),
	)
	if err != nil {
		panic(err)
	}
	cli, err := ins.NewClient(
		client.WithClientURL("127.0.0.1:20000"),
		client.WithClientProtocolDubbo(),
	)
	if err != nil {
		panic(err)
	}
	conn, err := cli.Dial("UserProvider", client.WithSerialization(constant.Hessian2Serialization))
	if err != nil {
		panic(err)
	}
	test(conn)
}

func test(conn *client.Connection) {
	ctx := tm.Begin(context.Background(), "TestTCCServiceBusiness")
	proxy, err := tcc.NewTCCServiceProxy(&service.UserProvider{})
	if err != nil {
		logger.Error(err)
	}
	defer func() {
		err = SeataClient.CommitOrRollback(ctx, err == nil)
		if err != nil {
			logger.Errorf("response commit or rollback: %v", err)
			return
		}
		logger.Info("complete commit or rollback")
	}()
	if resp, err := SeataClient.Prepare(proxy, ctx, conn, 1); err != nil {
		logger.Errorf("response prepare: %v", resp)
		return
	} else {
		logger.Infof("get resp %#v", resp)
	}
	// business
}
