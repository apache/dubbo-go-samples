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
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/registry"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-samples/generic/go-client/pkg"
)

const (
	RegistryAddress    = "127.0.0.1:2181"
	UserProvider       = "org.apache.dubbo.samples.UserProvider"
	ServiceVersion     = "1.0.0"
	ServiceGroupTriple = "dubbo"
)

func createServiceConnection(cli *client.Client, serviceInterface string) (*client.Connection, error) {
	return cli.Dial(
		serviceInterface,
		client.WithGeneric(),
		client.WithVersion(ServiceVersion),
		client.WithGroup(ServiceGroupTriple),
	)
}

func main() {
	hessian.RegisterPOJO(&pkg.User{})

	ins, err := dubbo.NewInstance(
		dubbo.WithName("generic-dubbo-client"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress(RegistryAddress),
		),
	)
	if err != nil {
		logger.Fatalf("Failed to create Dubbo instance: %v", err)
	}

	tripleCli, err := ins.NewClient(
		client.WithClientProtocolDubbo(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		logger.Fatalf("Failed to create Dubbo client: %v", err)
	}

	tripleConn, err := createServiceConnection(tripleCli, UserProvider)
	if err != nil {
		logger.Fatalf("Failed to create Dubbo connection: %v", err)
	}

	logger.Info("\n=== Testing Dubbo Protocol ===")
	testUserService(tripleConn)
}

func testUserService(conn *client.Connection) {
	call := func(methodName string, params []interface{}) (interface{}, error) {
		var result interface{}
		err := conn.CallUnary(
			context.TODO(),
			params,
			&result,
			methodName,
		)
		return result, err
	}

	testUserID := "A003"
	testUserName := "lily"
	testUserCode := int32(1)
	testUserIDs := []string{"001", "002", "003", "004"}

	testUser := &pkg.User{
		ID:   "3213",
		Name: "panty",
		Age:  25,
		Time: time.Now(),
	}

	testUsers := []*pkg.User{
		{
			ID:   "3212",
			Name: "XavierNiu",
			Age:  24,
			Time: time.Now().Add(4),
		},
		{
			ID:   "3213",
			Name: "zhangsan",
			Age:  21,
			Time: time.Now().Add(4),
		},
	}

	result, err := call("GetUser1", []interface{}{testUserID})
	logResult("GetUser1", result, err)

	result, err = call("GetUser2", []interface{}{testUserID, testUserName})
	logResult("GetUser2", result, err)

	result, err = call("GetUser3", []interface{}{testUserCode})
	logResult("GetUser3", result, err)

	result, err = call("GetUser4", []interface{}{testUserCode, "zhangsan"})
	logResult("GetUser4", result, err)

	result, err = call("GetUsers", []interface{}{testUserIDs})
	logResult("GetUsers", result, err)

	result, err = call("GetUsersMap", []interface{}{testUserIDs})
	logResult("GetUsersMap", result, err)

	result, err = call("QueryUser", []interface{}{testUser})
	logResult("QueryUser", result, err)

	result, err = call("QueryUsers", []interface{}{testUsers})
	logResult("QueryUsers", result, err)
}

func logResult(methodName string, result interface{}, err error) {
	if err != nil {
		logger.Errorf("❌ %s failed: %v", methodName, err)
	} else {
		logger.Infof("✅ %s succeeded: %+v", methodName, result)
	}
}
