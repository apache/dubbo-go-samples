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
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log"
)

var (
	appName         = "UserConsumer"
	referenceConfig = config.ReferenceConfig{
		InterfaceName: "org.apache.dubbo.UserProvider",
		Cluster:       "failover",
		Registry:      "demoZk",
		Protocol:      dubbo.DUBBO,
		Generic:       "true",
	}
)

func init() {
	config.Load()
	referenceConfig.GenericLoad(appName) //appName is the unique identification of RPCService
	time.Sleep(3 * time.Second)
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	gxlog.CInfo("\n\ncall getUser")
	callGetUser()
	gxlog.CInfo("\n\ncall queryUser")
	callQueryUser()
	gxlog.CInfo("\n\ncall queryUsers")
	callQueryUsers()
	gxlog.CInfo("\n\ncall callGetOneUser")
	callGetOneUser()

	initSignal()
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
			time.AfterFunc(10*time.Second, func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})

			// The program exits normally or timeout forcibly exits.
			fmt.Println("app exit now...")
			return
		}
	}
}

func callGetUser() {
	gxlog.CInfo("\n\n\nstart to generic invoke")
	resp, err := referenceConfig.GetRPCService().(*config.GenericService).Invoke(
		context.TODO(),
		[]interface{}{
			"GetUser",
			[]string{"java.lang.String"},
			[]hessian.Object{"A003"},
		},
	)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("res: %+v\n", resp)
	gxlog.CInfo("success!")

}
func callQueryUser() {
	gxlog.CInfo("\n\n\nstart to generic invoke")
	resp, err := referenceConfig.GetRPCService().(*config.GenericService).Invoke(
		context.TODO(),
		[]interface{}{
			"queryUser",
			[]string{"org.apache.dubbo.User"},
			// the map represents a User object:
			// &User {
			// 		ID: "3213",
			// 		Name: "panty",
			// 		Age: 25,
			// 		Time: time.Now(),
			// }
			[]hessian.Object{map[string]hessian.Object{
				"iD":   "3213",
				"name": "panty",
				"age":  25,
				"time": time.Now(),
			}},
		},
	)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("res: %+v\n", resp)
	gxlog.CInfo("success!")
}

func callQueryUsers() {
	resp, err := referenceConfig.GetRPCService().(*config.GenericService).Invoke(
		context.TODO(),
		[]interface{}{
			"QueryUsers",
			[]string{"java.lang.Array"},
			[]hessian.Object{
				[]hessian.Object{
					map[string]hessian.Object{
						"iD":   "3213",
						"name": "panty",
						"age":  25,
						"time": time.Now(),
					},
					map[string]hessian.Object{
						"iD":   "3212",
						"name": "XavierNiu",
						"age":  24,
						"time": time.Now().Add(4),
					},
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("res: %+v\n", resp)
	gxlog.CInfo("success!")
}

func callGetOneUser() {
	gxlog.CInfo("\n\n\nstart to generic invoke")
	resp, err := referenceConfig.GetRPCService().(*config.GenericService).Invoke(
		context.TODO(),
		[]interface{}{
			"GetOneUser",
			[]string{"org.apache.dubbo.User"},
			[]hessian.Object{},
		},
	)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("res: %+v\n", resp)
	gxlog.CInfo("success!")
}
