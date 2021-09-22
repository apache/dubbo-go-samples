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
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	"github.com/apache/dubbo-go-samples/rpc/jsonrpc/go-client/pkg"
)

var (
	survivalTimeout int = 10e9
	userProvider        = &pkg.UserProvider{}
	userProvider1       = &pkg.UserProvider1{}
	userProvider2       = &pkg.UserProvider2{}
)

func init() {
	config.SetConsumerService(userProvider)
	config.SetConsumerService(userProvider1)
	config.SetConsumerService(userProvider2)
}

// Do some checking before the system starts up:
// 1. env config
// 		`export DUBBO_GO_CONFIG_PATH= ROOT_PATH/conf/dubbogo.yml` or `dubbogo.yaml`
func main() {

	config.Load()

	logger.Info("\n\ntest")
	test()
	logger.Info("\n\ntest1")
	test1()
	logger.Info("\n\ntest2")
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
	logger.Info("\n\n\necho")
	res, err := userProvider.Echo(context.TODO(), "OK")
	if err != nil {
		logger.Info("echo - error: %v", err)
	} else {
		logger.Info("res: %v", res)
	}

	time.Sleep(3e9)

	logger.Info("\n\n\nstart to test jsonrpc")

	user, err := userProvider.GetUser(context.TODO(), "A003")

	if err != nil {
		panic(err)
	}
	logger.Info("response result: %v", user)

	logger.Info("\n\n\nstart to test jsonrpc - GetUser0")
	ret, err := userProvider.GetUser0("A003", "Moorse")
	if err != nil {
		panic(err)
	}
	logger.Info("response result: %v", ret)

	logger.Info("\n\n\nstart to test jsonrpc - GetUsers")

	ret1, err := userProvider.GetUsers([]interface{}{[]interface{}{"A002", "A003"}})
	if err != nil {
		panic(err)
	}
	logger.Info("response result: %v", ret1)

	logger.Info("\n\n\nstart to test jsonrpc - getUser")
	rep2, err := userProvider.GetUser2(context.TODO(), "1")
	if err != nil {
		panic(err)
	}
	logger.Info("response result: %v", rep2)

	logger.Info("\n\n\nstart to test jsonrpc - GetUser3")
	err = userProvider.GetUser3()
	if err != nil {
		panic(err)
	}
	logger.Info("succ!")

	logger.Info("\n\n\nstart to test jsonrpc illegal method")
	rep3, err := userProvider.GetUser1(context.TODO(), "A003")
	if err == nil {
		panic("err is nil")
	}
	logger.Info("response result: %v", rep3)
}

func test1() {
	logger.Info("\n\n\necho")
	res, err := userProvider1.Echo(context.TODO(), "OK")
	if err != nil {
		logger.Info("echo - error: %v", err)
	} else {
		logger.Info("res: %v", res)
	}

	time.Sleep(3e9)

	logger.Info("\n\n\nstart to test jsonrpc")
	user, err := userProvider1.GetUser(context.TODO(), "A003")
	if err != nil {
		panic(err)
	}
	logger.Info("response result: %v", user)

	logger.Info("\n\n\nstart to test jsonrpc - GetUser0")
	ret, err := userProvider1.GetUser0("A003", "Moorse")
	if err != nil {
		panic(err)
	}
	logger.Info("response result: %v", ret)

	logger.Info("\n\n\nstart to test jsonrpc - getUser")
	user, err = userProvider1.GetUser2(context.TODO(), "1")
	if err != nil {
		panic(err)
	}
	logger.Info("response result: %v", user)

	logger.Info("\n\n\nstart to test jsonrpc - GetUser3")
	err = userProvider1.GetUser3()
	if err != nil {
		panic(err)
	}
	logger.Info("succ!")

	logger.Info("\n\n\nstart to test jsonrpc illegal method")
	user, err = userProvider1.GetUser1(context.TODO(), "A003")
	if err == nil {
		panic("err is nil")
	}
	logger.Info("error: %v", err)
}

func test2() {
	logger.Info("\n\n\necho")
	res, err := userProvider2.Echo(context.TODO(), "OK")
	if err != nil {
		logger.Info("echo - error: %v", err)
	} else {
		logger.Info("res: %v", res)
	}

	time.Sleep(3e9)

	logger.Info("\n\n\nstart to test jsonrpc")
	user, err := userProvider2.GetUser(context.TODO(), "A003")
	if err != nil {
		panic(err)
	}
	logger.Info("response result: %v", user)

	logger.Info("\n\n\nstart to test jsonrpc - GetUser0")
	ret, err := userProvider2.GetUser0("A003", "Moorse")
	if err != nil {
		panic(err)
	}
	logger.Info("response result: %v", ret)

	logger.Info("\n\n\nstart to test jsonrpc - getUser")
	user, err = userProvider2.GetUser2(context.TODO(), "1")
	if err != nil {
		panic(err)
	}
	logger.Info("response result: %v", user)

	logger.Info("\n\n\nstart to test jsonrpc - GetUser3")
	err = userProvider2.GetUser3()
	if err != nil {
		panic(err)
	}
	logger.Info("succ!")

	logger.Info("\n\n\nstart to test jsonrpc illegal method")
	user, err = userProvider2.GetUser1(context.TODO(), "A003")
	if err == nil {
		panic("err is nil")
	}
	logger.Info("error: %v", err)
}
