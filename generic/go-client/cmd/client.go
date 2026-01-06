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
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/apache/dubbo-go-samples/generic/go-client/pkg"

	"github.com/dubbogo/gost/log/logger"
)

const (
	DubboServerURL  = "dubbo://127.0.0.1:20000"
	TripleServerURL = "tri://127.0.0.1:50052"
	UserProvider    = "org.apache.dubbo.samples.UserProvider"
	ServiceVersion  = "1.0.0"
)

func main() {
	hessian.RegisterPOJO(&pkg.User{})

	ins, err := dubbo.NewInstance(
		dubbo.WithName("generic-go-client"),
	)
	if err != nil {
		panic(err)
	}

	// Test Dubbo protocol generic call
	logger.Info("=== Testing Dubbo Protocol Generic Call ===")
	testDubboProtocol(ins)

	// Test Triple protocol generic call
	logger.Info("=== Testing Triple Protocol Generic Call ===")
	testTripleProtocol(ins)

	logger.Info("All generic call tests completed")
}

func testDubboProtocol(ins *dubbo.Instance) {
	cli, err := ins.NewClient(
		client.WithClientProtocolDubbo(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		logger.Errorf("Failed to create Dubbo client: %v", err)
		return
	}

	genericService, err := cli.NewGenericService(
		UserProvider,
		client.WithURL(DubboServerURL),
		client.WithVersion(ServiceVersion),
		client.WithGroup("dubbo"),
	)
	if err != nil {
		logger.Errorf("Failed to create Dubbo generic service: %v", err)
		return
	}

	runGenericTests(genericService, "Dubbo")
}

func testTripleProtocol(ins *dubbo.Instance) {
	cli, err := ins.NewClient(
		client.WithClientProtocolTriple(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		logger.Errorf("Failed to create Triple client: %v", err)
		return
	}

	genericService, err := cli.NewGenericService(
		UserProvider,
		client.WithURL(TripleServerURL),
		client.WithVersion(ServiceVersion),
		client.WithGroup("triple"),
	)
	if err != nil {
		logger.Errorf("Failed to create Triple generic service: %v", err)
		return
	}

	runGenericTests(genericService, "Triple")
}

func runGenericTests(svc *generic.GenericService, protocol string) {
	// 1. Basic type parameters
	testBasicTypes(svc, protocol)

	// 2. Array/Collection types
	testCollectionTypes(svc, protocol)

	// 3. Custom POJO types
	testPOJOTypes(svc, protocol)
}

func testBasicTypes(svc *generic.GenericService, protocol string) {
	ctx := context.Background()

	// GetUser1(String)
	result, err := svc.Invoke(ctx, "GetUser1", []string{"java.lang.String"}, []hessian.Object{"A003"})
	if err != nil {
		logger.Errorf("[%s] GetUser1 failed: %v", protocol, err)
	} else {
		logger.Infof("[%s] GetUser1(userId string) res: %+v", protocol, result)
	}

	// GetUser2(String, String)
	result, err = svc.Invoke(ctx, "GetUser2", []string{"java.lang.String", "java.lang.String"}, []hessian.Object{"A003", "lily"})
	if err != nil {
		logger.Errorf("[%s] GetUser2 failed: %v", protocol, err)
	} else {
		logger.Infof("[%s] GetUser2(userId string, name string) res: %+v", protocol, result)
	}

	// GetUser3(int)
	result, err = svc.Invoke(ctx, "GetUser3", []string{"int"}, []hessian.Object{int32(1)})
	if err != nil {
		logger.Errorf("[%s] GetUser3 failed: %v", protocol, err)
	} else {
		logger.Infof("[%s] GetUser3(userCode int) res: %+v", protocol, result)
	}

	// GetUser4(int, String)
	result, err = svc.Invoke(ctx, "GetUser4", []string{"int", "java.lang.String"}, []hessian.Object{int32(1), "zhangsan"})
	if err != nil {
		logger.Errorf("[%s] GetUser4 failed: %v", protocol, err)
	} else {
		logger.Infof("[%s] GetUser4(userCode int, name string) res: %+v", protocol, result)
	}

	// GetOneUser()
	result, err = svc.Invoke(ctx, "GetOneUser", []string{}, []hessian.Object{})
	if err != nil {
		logger.Errorf("[%s] GetOneUser failed: %v", protocol, err)
	} else {
		logger.Infof("[%s] GetOneUser() res: %+v", protocol, result)
	}
}

func testCollectionTypes(svc *generic.GenericService, protocol string) {
	ctx := context.Background()

	// GetUsers(String[])
	result, err := svc.Invoke(ctx, "GetUsers", []string{"[Ljava.lang.String;"}, []hessian.Object{[]string{"001", "002", "003"}})
	if err != nil {
		logger.Errorf("[%s] GetUsers failed: %v", protocol, err)
	} else {
		logger.Infof("[%s] GetUsers(userIdList []string) res: %+v", protocol, result)
	}

	// GetUsersMap(String[])
	result, err = svc.Invoke(ctx, "GetUsersMap", []string{"[Ljava.lang.String;"}, []hessian.Object{[]string{"001", "002"}})
	if err != nil {
		logger.Errorf("[%s] GetUsersMap failed: %v", protocol, err)
	} else {
		logger.Infof("[%s] GetUsersMap(userIdList []string) res: %+v", protocol, result)
	}

	// QueryAll()
	result, err = svc.Invoke(ctx, "QueryAll", []string{}, []hessian.Object{})
	if err != nil {
		logger.Errorf("[%s] QueryAll failed: %v", protocol, err)
	} else {
		logger.Infof("[%s] QueryAll() res: %+v", protocol, result)
	}
}

func testPOJOTypes(svc *generic.GenericService, protocol string) {
	ctx := context.Background()

	// QueryUser(User)
	testUser := &pkg.User{
		ID:   "3213",
		Name: "panty",
		Age:  25,
		Time: time.Now(),
	}
	result, err := svc.Invoke(ctx, "QueryUser", []string{"org.apache.dubbo.samples.User"}, []hessian.Object{testUser})
	if err != nil {
		logger.Errorf("[%s] QueryUser failed: %v", protocol, err)
	} else {
		logger.Infof("[%s] QueryUser(user *User) res: %+v", protocol, result)
	}

	// QueryUsers(User[])
	testUsers := []*pkg.User{
		{ID: "3212", Name: "XavierNiu", Age: 24, Time: time.Now()},
		{ID: "3213", Name: "zhangsan", Age: 21, Time: time.Now()},
	}
	result, err = svc.Invoke(ctx, "QueryUsers", []string{"[Lorg.apache.dubbo.samples.User;"}, []hessian.Object{testUsers})
	if err != nil {
		logger.Errorf("[%s] QueryUsers failed: %v", protocol, err)
	} else {
		logger.Infof("[%s] QueryUsers(users []*User) res: %+v", protocol, result)
	}
}
