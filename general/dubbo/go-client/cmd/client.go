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
	"time"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go-samples/general/dubbo/go-client/pkg"
	"github.com/apache/dubbo-go/config"
	"github.com/dubbogo/gost/log"
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

var (
	userProvider        = new(pkg.UserProvider)
	userProvider1       = new(pkg.UserProvider1)
	userProvider2       = new(pkg.UserProvider2)
)

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	hessian.RegisterJavaEnum(pkg.Gender(pkg.MAN))
	hessian.RegisterJavaEnum(pkg.Gender(pkg.WOMAN))
	hessian.RegisterPOJO(&pkg.User{})

	config.SetConsumerService(userProvider)

	// FIXME: cannot make multi-references work
	config.SetConsumerService(userProvider1)
	config.SetConsumerService(userProvider2)
	config.Load()

	time.Sleep(6 * time.Second)

	gxlog.CInfo("\n\ntest")
	test()
	gxlog.CInfo("\n\ntest1")
	test1()
	gxlog.CInfo("\n\ntest2")
	test2()
}

func test() {

	gxlog.CInfo("\n\n\necho")
	res, err := userProvider.Echo(context.TODO(), "OK")
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("res: %v\n", res)

	gxlog.CInfo("\n\n\nstart to test dubbo")
	user := &pkg.User{}
	err = userProvider.GetUser(context.TODO(), []interface{}{"A003"}, user)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", user)

	gxlog.CInfo("\n\n\nstart to test dubbo - enum")
	gender, err := userProvider.GetGender(1)
	if err != nil {
		gxlog.CInfo("error: %v", err)
	} else {
		gxlog.CInfo("response result: %v", gender)
	}

	gxlog.CInfo("\n\n\nstart to test dubbo - GetUser0")
	ret, err := userProvider.GetUser0("A003", "Moorse")
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", ret)

	gxlog.CInfo("\n\n\nstart to test dubbo - GetUsers")
	ret1, err := userProvider.GetUsers([]interface{}{[]interface{}{"A002", "A003"}})
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", ret1)

	gxlog.CInfo("\n\n\nstart to test dubbo - getUser")
	user = &pkg.User{}
	var i int32 = 1
	err = userProvider.GetUser2(context.TODO(), []interface{}{i}, user)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", user)

	gxlog.CInfo("\n\n\nstart to test dubbo - getUser - overload")
	user = &pkg.User{}
	err = userProvider.GetUser2(context.TODO(), []interface{}{i, "overload"}, user)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", user)

	gxlog.CInfo("\n\n\nstart to test dubbo - GetUser3")
	err = userProvider.GetUser3()
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("succ!")

	gxlog.CInfo("\n\n\nstart to test dubbo - getErr")
	user = &pkg.User{}
	err = userProvider.GetErr(context.TODO(), []interface{}{"A003"}, user)
	if err == nil {
		panic("err is nil")
	}
	gxlog.CInfo("getErr - error: %v", err)

	gxlog.CInfo("\n\n\nstart to test dubbo illegal method")
	err = userProvider.GetUser1(context.TODO(), []interface{}{"A003"}, user)
	if err == nil {
		panic("err is nil")
	}
	gxlog.CInfo("error: %v", err)
}

func test1() {
	gxlog.CInfo("\n\n\necho")
	res, err := userProvider1.Echo(context.TODO(), "OK")
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("res: %v\n", res)

	gxlog.CInfo("\n\n\nstart to test dubbo")
	user := &pkg.User{}
	err = userProvider1.GetUser(context.TODO(), []interface{}{"A003"}, user)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", user)

	gxlog.CInfo("\n\n\nstart to test dubbo - GetUser0")
	ret, err := userProvider1.GetUser0("A003", "Moorse")
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", ret)

	gxlog.CInfo("\n\n\nstart to test dubbo - GetUsers")
	ret1, err := userProvider1.GetUsers([]interface{}{[]interface{}{"A002", "A003"}})
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", ret1)

	gxlog.CInfo("\n\n\nstart to test dubbo - getUser")
	user = &pkg.User{}
	var i int32 = 1
	err = userProvider1.GetUser2(context.TODO(), []interface{}{i}, user)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", user)

	gxlog.CInfo("\n\n\nstart to test dubbo - getUser - overload")
	user = &pkg.User{}
	err = userProvider1.GetUser2(context.TODO(), []interface{}{i, "overload"}, user)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", user)

	gxlog.CInfo("\n\n\nstart to test dubbo - GetUser3")
	err = userProvider1.GetUser3()
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("succ!")

	gxlog.CInfo("\n\n\nstart to test dubbo - getErr")
	user = &pkg.User{}
	err = userProvider1.GetErr(context.TODO(), []interface{}{"A003"}, user)
	if err == nil {
		panic("err is nil")
	}
	gxlog.CInfo("getErr - error: %v", err)

	gxlog.CInfo("\n\n\nstart to test dubbo illegal method")
	err = userProvider1.GetUser1(context.TODO(), []interface{}{"A003"}, user)
	if err == nil {
		panic("err is nil")
	}
	gxlog.CInfo("error: %v", err)
}

func test2() {
	gxlog.CInfo("\n\n\necho")
	res, err := userProvider2.Echo(context.TODO(), "OK")
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("res: %v\n", res)

	time.Sleep(3e9)

	gxlog.CInfo("\n\n\nstart to test dubbo")
	user := &pkg.User{}
	err = userProvider2.GetUser(context.TODO(), []interface{}{"A003"}, user)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", user)

	gxlog.CInfo("\n\n\nstart to test dubbo - GetUser0")
	ret, err := userProvider2.GetUser0("A003", "Moorse")
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", ret)

	gxlog.CInfo("\n\n\nstart to test dubbo - GetUsers")
	ret1, err := userProvider2.GetUsers([]interface{}{[]interface{}{"A002", "A003"}})
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", ret1)

	gxlog.CInfo("\n\n\nstart to test dubbo - getUser")
	user = &pkg.User{}
	var i int32 = 1
	err = userProvider2.GetUser2(context.TODO(), []interface{}{i}, user)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", user)

	gxlog.CInfo("\n\n\nstart to test dubbo - getUser - overload")
	user = &pkg.User{}
	err = userProvider2.GetUser2(context.TODO(), []interface{}{i, "overload"}, user)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("response result: %v", user)

	gxlog.CInfo("\n\n\nstart to test dubbo - GetUser3")
	err = userProvider2.GetUser3()
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("succ!")

	gxlog.CInfo("\n\n\nstart to test dubbo - getErr")
	user = &pkg.User{}
	err = userProvider2.GetErr(context.TODO(), []interface{}{"A003"}, user)
	if err == nil {
		panic("err is nil")
	}
	gxlog.CInfo("getErr - error: %v", err)

	gxlog.CInfo("\n\n\nstart to test dubbo illegal method")
	err = userProvider2.GetUser1(context.TODO(), []interface{}{"A003"}, user)
	if err == nil {
		panic("err is nil")
	}
	gxlog.CInfo("error: %v", err)
}
