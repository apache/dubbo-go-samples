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
	"os/signal"
	"syscall"
	"time"
)

import (
	"github.com/dubbogo/gost/log"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	"dubbo.apache.org/dubbo-go/v3/protocol/jsonrpc"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
)

var (
	survivalTimeout int = 10e9
)

var (
	appName         = "UserConsumer"
	referenceConfig = config.ReferenceConfig{
		InterfaceName: "org.apache.dubbo.UserProvider",
		Cluster:       "failover",
		Registry:      "demoZk",
		Protocol:      jsonrpc.JSONRPC,
		Generic:       true,
	}
)

func init() {
	config.Load()
	referenceConfig.GenericLoad(appName)
	time.Sleep(3 * time.Second)
}

func main() {
	callGetUser()
	initSignal()
}

func callGetUser() {
	ctx := context.WithValue(
		context.Background(),
		constant.AttachmentKey,
		map[string]interface{}{
			"generic": "jsonrpc",
		},
	)
	resp, err := referenceConfig.GetRPCService().(*config.GenericService).Invoke(
		ctx,
		[]interface{}{
			"GetUser",
			[]string{"java.lang.String"},
			[]interface{}{"A000"},
		},
	)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("res: %+v\n", resp)
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP,
		syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
		// reload()
		default:
			time.AfterFunc(time.Duration(survivalTimeout), func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})

			// The program exits normally or timeout forcibly exits.
			fmt.Println("app exit now...")
			return
		}
	}
}
