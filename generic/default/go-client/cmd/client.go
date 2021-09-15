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
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"

	hessian "github.com/apache/dubbo-go-hessian2"
)

import (
	"github.com/apache/dubbo-go-samples/generic/default/go-client/pkg"
)

var (
	appName = "dubbo.io"
)

// export DUBBO_GO_CONFIG_PATH= PATH_TO_SAMPLES/generic/default/go-client/conf/dubbogo.yml
func main() {
	// register POJOs
	hessian.RegisterPOJO(&pkg.User{})

	// generic invocation samples using hessian serialization on Dubbo protocol
	dubboRefConf := newRefConf("org.apache.dubbo.samples.UserProvider", dubbo.DUBBO)

	callGetUser(dubboRefConf)
	//callGetOneUser(dubboRefConf)
	callGetUsers(dubboRefConf)
	callGetUsersMap(dubboRefConf)
	callQueryUser(dubboRefConf)
	callQueryUsers(dubboRefConf)
	//callQueryAll(dubboRefConf)

	// generic invocation samples using hessian serialization on Triple protocol
	//tripleRefConf := newRefConf("org.apache.dubbo.samples.UserProviderTriple", tpconst.TRIPLE)

	initSignal()
}

func callGetUser(refConf config.ReferenceConfig) {
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser1",
		[]string{"java.lang.String"},
		[]hessian.Object{"A003"},
	)

	if err != nil {
		panic(err)
	}
	logger.Infof("GetUser1(userId string) res: %+v", resp)

	resp, err = refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser2",
		[]string{"java.lang.String", "java.lang.String"},
		[]hessian.Object{"A003", "lily"},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("GetUser2(userId string, name string) res: %+v", resp)

	resp, err = refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser3",
		[]string{"int"},
		[]hessian.Object{1},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("GetUser3(userCode int) res: %+v", resp)

	resp, err = refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser4",
		[]string{"int", "java.lang.String"},
		[]hessian.Object{1, "zhangsan"},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("GetUser4(userCode int, name string) res: %+v", resp)
}

func callGetOneUser(refConf config.ReferenceConfig) {
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetOneUser",
		[]string{},
		// TODO go-go []hessian.Object{}, go-java []string{}
		[]hessian.Object{},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("GetOneUser() res: %+v", resp)
}

func callGetUsers(refConf config.ReferenceConfig) {
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUsers",
		[]string{"java.util.List"},
		[]hessian.Object{
			[]hessian.Object{
				"001", "002", "003", "004",
			},
		},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("GetUsers1(userIdList []*string) res: %+v", resp)
}

func callGetUsersMap(refConf config.ReferenceConfig) {
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUsersMap",
		[]string{"java.util.List"},
		[]hessian.Object{
			[]hessian.Object{
				"001", "002", "003", "004",
			},
		},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("GetUserMap(userIdList []*string) res: %+v", resp)
}

func callQueryUser(refConf config.ReferenceConfig) {
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"queryUser",
		[]string{"org.apache.dubbo.User"},
		// the map represents a User object:
		// &User {
		// 		ID: "3213",
		// 		Name: "panty",
		// 		Age: 25,
		// 		Time: time.Now(),
		// }
		[]hessian.Object{
			map[string]hessian.Object{
				"iD":   "3213",
				"name": "panty",
				"age":  25,
				"time": time.Now(),
			}},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("queryUser(user *User) res: %+v", resp)
}

func callQueryUsers(refConf config.ReferenceConfig) {
	var resp, err = refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"queryUsers",
		[]string{"java.util.ArrayList"},
		[]hessian.Object{
			[]hessian.Object{
				map[string]hessian.Object{
					"id":    "3212",
					"name":  "XavierNiu",
					"age":   24,
					"time":  time.Now().Add(4),
					"class": "org.apache.dubbo.User",
				},
				map[string]hessian.Object{
					"iD":    "3213",
					"name":  "zhangsan",
					"age":   21,
					"time":  time.Now().Add(4),
					"class": "org.apache.dubbo.User",
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("queryUsers(users []*User) res: %+v", resp)
}

func callQueryAll(refConf config.ReferenceConfig) {
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"queryAll",
		[]string{},
		// TODO go-go []hessian.Object{}, go-java []string{}
		//[]hessian.Object{},
		[]hessian.Object{},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("queryAll() res: %+v", resp)
}

func newRefConf(iface, protocol string) config.ReferenceConfig {
	registryConfig := &config.RegistryConfig{
		Protocol: "zookeeper",
		Address:  "127.0.0.1:2181",
	}

	refConf := config.ReferenceConfig{
		InterfaceName: iface,
		Cluster:       "failover",
		Registry:      []string{"zk"},
		Protocol:      protocol,
		Generic:       "true",
	}

	rootConfig := config.NewRootConfig(config.WithRootRegistryConfig("zk", registryConfig))
	_ = rootConfig.Init()
	_ = refConf.Init(rootConfig)
	refConf.GenericLoad(appName)

	return refConf
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
