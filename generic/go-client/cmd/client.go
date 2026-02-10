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
	"os"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/filter/generic"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-samples/generic/go-client/pkg"
)

const (
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

	cli, err := ins.NewClient(
		client.WithClientProtocolTriple(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		panic(err)
	}

	genericService, err := cli.NewGenericService(
		UserProvider,
		client.WithURL(TripleServerURL),
		client.WithVersion(ServiceVersion),
		client.WithGroup("triple"),
	)
	if err != nil {
		panic(err)
	}

	failed := false
	failed = runGenericTests(genericService) || failed

	if failed {
		logger.Errorf("Some generic call tests failed")
		os.Exit(1)
	}
	logger.Info("All generic call tests passed")
}

func runGenericTests(svc *generic.GenericService) bool {
	failed := false
	ctx := context.Background()

	// GetUser1(String)
	result, err := svc.Invoke(ctx, "GetUser1", []string{"java.lang.String"}, []hessian.Object{"A003"})
	if err != nil {
		logger.Errorf("GetUser1 failed: %v", err)
		failed = true
	} else {
		logger.Infof("GetUser1(userId string) res: %+v", result)
	}

	// GetUser2(String, String)
	result, err = svc.Invoke(ctx, "GetUser2", []string{"java.lang.String", "java.lang.String"}, []hessian.Object{"A003", "lily"})
	if err != nil {
		logger.Errorf("GetUser2 failed: %v", err)
		failed = true
	} else {
		logger.Infof("GetUser2(userId string, name string) res: %+v", result)
	}

	// GetUser3(int)
	result, err = svc.Invoke(ctx, "GetUser3", []string{"int"}, []hessian.Object{int32(1)})
	if err != nil {
		logger.Errorf("GetUser3 failed: %v", err)
		failed = true
	} else {
		logger.Infof("GetUser3(userCode int) res: %+v", result)
	}

	// GetUser4(int, String)
	result, err = svc.Invoke(ctx, "GetUser4", []string{"int", "java.lang.String"}, []hessian.Object{int32(1), "zhangsan"})
	if err != nil {
		logger.Errorf("GetUser4 failed: %v", err)
		failed = true
	} else {
		logger.Infof("GetUser4(userCode int, name string) res: %+v", result)
	}

	// GetOneUser()
	result, err = svc.Invoke(ctx, "GetOneUser", []string{}, []hessian.Object{})
	if err != nil {
		logger.Errorf("GetOneUser failed: %v", err)
		failed = true
	} else {
		logger.Infof("GetOneUser() res: %+v", result)
	}

	// GetUsers(String[])
	result, err = svc.Invoke(ctx, "GetUsers", []string{"[Ljava.lang.String;"}, []hessian.Object{[]string{"001", "002", "003"}})
	if err != nil {
		logger.Errorf("GetUsers failed: %v", err)
		failed = true
	} else {
		logger.Infof("GetUsers(userIdList []string) res: %+v", result)
	}

	// GetUsersMap(String[])
	result, err = svc.Invoke(ctx, "GetUsersMap", []string{"[Ljava.lang.String;"}, []hessian.Object{[]string{"001", "002"}})
	if err != nil {
		logger.Errorf("GetUsersMap failed: %v", err)
		failed = true
	} else {
		logger.Infof("GetUsersMap(userIdList []string) res: %+v", result)
	}

	// QueryAll()
	result, err = svc.Invoke(ctx, "QueryAll", []string{}, []hessian.Object{})
	if err != nil {
		logger.Errorf("QueryAll failed: %v", err)
		failed = true
	} else {
		logger.Infof("QueryAll() res: %+v", result)
	}

	// QueryUser(User)
	testUser := &pkg.User{
		ID:   "3213",
		Name: "panty",
		Age:  25,
		Time: time.Now(),
	}
	result, err = svc.Invoke(ctx, "QueryUser", []string{"org.apache.dubbo.samples.User"}, []hessian.Object{testUser})
	if err != nil {
		logger.Errorf("QueryUser failed: %v", err)
		failed = true
	} else {
		logger.Infof("QueryUser(user *User) res: %+v", result)
	}

	// QueryUsers(User[])
	testUsers := []*pkg.User{
		{ID: "3212", Name: "XavierNiu", Age: 24, Time: time.Now()},
		{ID: "3213", Name: "zhangsan", Age: 21, Time: time.Now()},
	}
	result, err = svc.Invoke(ctx, "QueryUsers", []string{"[Lorg.apache.dubbo.samples.User;"}, []hessian.Object{testUsers})
	if err != nil {
		logger.Errorf("QueryUsers failed: %v", err)
		failed = true
	} else {
		logger.Infof("QueryUsers(users []*User) res: %+v", result)
	}

	return failed
}
