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
	_ "dubbo.apache.org/dubbo-go/v3/protocol/rest"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"

	"github.com/dubbogo/gost/log"
)

import (
	"github.com/apache/dubbo-go-samples/general/rest/go-client/pkg"
)

var (
	survivalTimeout int = 10e9
)

// they are necessary:
// 		export CONF_CONSUMER_FILE_PATH="xxx"
// 		export APP_LOG_CONF_FILE="xxx"
func main() {

	config.Load()

	gxlog.CInfo("\n\ntest")
	test()
	gxlog.CInfo("\n\ntest1")
	test1()
	gxlog.CInfo("\n\ntest2")
	test2()

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

func test() {

	gxlog.CInfo("\n\n\nstart to test rest")
	user := &pkg.User{}
	err := pkg.UserProviderVar.GetUser(context.TODO(), []interface{}{"A003"}, user)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", user)

	gxlog.CInfo("\n\n\nstart to test rest - GetUser0")
	ret, err := pkg.UserProviderVar.GetUser0("A003", "Moorse中文", 30)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", ret)

	gxlog.CInfo("\n\n\nstart to test rest - GetUsers")
	ret1, err := pkg.UserProviderVar.GetUsers([]interface{}{&pkg.User{ID: "A002"}})
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", ret1)

	gxlog.CInfo("\n\n\nstart to test rest - GetUser3")
	err = pkg.UserProviderVar.GetUser3()
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("succ!")

	gxlog.CInfo("\n\n\nstart to test rest illegal method")
	err = pkg.UserProviderVar.GetUser1(context.TODO(), []interface{}{"A003"}, user)
	if err == nil {
		panic("err is nil")
	}
	gxlog.CInfo("error: %v", err)
}

func test1() {

	time.Sleep(3e9)

	gxlog.CInfo("\n\n\nstart to test rest")
	user := &pkg.User{}
	err := pkg.UserProvider1Var.GetUser(context.TODO(), []interface{}{"A003"}, user)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", user)

	gxlog.CInfo("\n\n\nstart to test rest - GetUser0")
	ret, err := pkg.UserProvider1Var.GetUser0("A003", "Moorse中文", 30)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", ret)

	gxlog.CInfo("\n\n\nstart to test rest - GetUsers")
	ret1, err := pkg.UserProvider1Var.GetUsers([]interface{}{&pkg.User{ID: "A002"}})
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", ret1)

	gxlog.CInfo("\n\n\nstart to test rest - GetUser3")
	err = pkg.UserProvider1Var.GetUser3()
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("succ!")

	gxlog.CInfo("\n\n\nstart to test rest illegal method")
	err = pkg.UserProvider1Var.GetUser1(context.TODO(), []interface{}{"A003"}, user)
	if err == nil {
		panic("err is nil")
	}
	gxlog.CInfo("error: %v", err)
}

func test2() {

	gxlog.CInfo("\n\n\nstart to test rest")
	user := &pkg.User{}
	err := pkg.UserProvider2Var.GetUser(context.TODO(), []interface{}{"A003"}, user)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", user)

	gxlog.CInfo("\n\n\nstart to test rest - GetUser0")
	ret, err := pkg.UserProvider2Var.GetUser0("A003", "Moorse中文", 30)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", ret)

	gxlog.CInfo("\n\n\nstart to test rest - GetUsers")
	ret1, err := pkg.UserProvider2Var.GetUsers([]interface{}{&pkg.User{ID: "A002"}})
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", ret1)

	gxlog.CInfo("\n\n\nstart to test rest - GetUser3")
	err = pkg.UserProvider2Var.GetUser3()
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("succ!")

	gxlog.CInfo("\n\n\nstart to test rest illegal method")
	err = pkg.UserProvider2Var.GetUser1(context.TODO(), []interface{}{"A003"}, user)
	if err == nil {
		panic("err is nil")
	}
	gxlog.CInfo("error: %v", err)
}
