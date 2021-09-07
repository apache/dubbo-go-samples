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

package integration

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/service/local"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"

	hessian "github.com/apache/dubbo-go-hessian2"
)

var (
	appName         = "dubbo.io"
	referenceConfig config.ReferenceConfig
)

func init() {
	registryConfig := &config.RegistryConfig{
		Protocol: "zookeeper",
		Address:  "127.0.0.1:2181",
	}

	referenceConfig = config.ReferenceConfig{
		InterfaceName: "org.apache.dubbo.UserProvider",
		Cluster:       "failover",
		Registry:      []string{"zk"},
		Protocol:      dubbo.DUBBO,
		Generic:       "true",
	}

	rootConfig := config.NewRootConfig(config.WithRootRegistryConfig("zk", registryConfig))
	_ = rootConfig.Init()
	_ = referenceConfig.Init(rootConfig)
	referenceConfig.GenericLoad(appName)
}

var userProvider = new(UserProvider)

func TestMain(m *testing.M) {
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&User{})
	config.Load()
	initUserMap()
	time.Sleep(3 * time.Second)

	os.Exit(m.Run())
}

type User struct {
	// !!! Cannot define lowercase names of variable
	ID   string
	Name string
	Age  int32
	Time time.Time
}

func (u User) String() string {
	return fmt.Sprintf(
		"User{ID:%s, Name:%s, Age:%d, Time:%s}",
		u.ID, u.Name, u.Age, u.Time,
	)
}

func (User) JavaClassName() string {
	return "org.apache.dubbo.User"
}

type UserResponse struct {
	Users []*User
}

type UserProvider struct {
}

var userMap = make(map[string]*User)

func initUserMap() {
	userMap["001"] = &User{"001", "other-zhangsan", 23, time.Date(1998, 1, 2, 3, 4, 5, 0, time.Local)}
	userMap["002"] = &User{"002", "other-lisi", 25, time.Date(1996, 1, 2, 3, 4, 5, 0, time.Local)}
	userMap["003"] = &User{"003", "other-lily", 28, time.Date(1993, 1, 2, 3, 4, 5, 0, time.Local)}
	userMap["004"] = &User{"004", "other-lisa", 36, time.Date(1985, 1, 2, 3, 4, 5, 0, time.Local)}
}

func (u *UserProvider) GetUser1(_ context.Context, userID string) (*User, error) {
	logger.Infof("req:%#v", userID)
	rsp := User{userID, "Joe", 48, time.Now()}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) GetUser2(_ context.Context, userID string, name string) (*User, error) {
	logger.Infof("req:%#v, %#v", userID, name)
	rsp := User{userID, name, 48, time.Now()}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) GetUser3(_ context.Context, userCode int) (*User, error) {
	logger.Infof("req:%#v", userCode)
	rsp := User{strconv.Itoa(userCode), "Alex Stocks", 18, time.Now()}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) GetUser4(_ context.Context, userCode int, name string) (*User, error) {
	logger.Infof("req:%#v, %#v", userCode, name)
	rsp := User{strconv.Itoa(userCode), name, 18, time.Now()}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) GetOneUser(_ context.Context) (*User, error) {
	return &User{
		ID:   "1000",
		Name: "xavierniu",
		Age:  24,
		Time: time.Now(),
	}, nil
}

func (u *UserProvider) GetUsers(_ context.Context, userIdList []string) (*UserResponse, error) {
	logger.Infof("req:%#v", userIdList)
	var users []*User
	for _, i := range userIdList {
		users = append(users, userMap[i])
	}
	return &UserResponse{
		Users: users,
	}, nil
}

func (u *UserProvider) GetUsersMap(_ context.Context, userIdList []string) (map[string]*User, error) {
	logger.Infof("req:%#v", userIdList)
	var users = make(map[string]*User)
	for _, i := range userIdList {
		users[i] = userMap[i]
	}
	return users, nil
}

func (u *UserProvider) QueryUser(_ context.Context, user *User) (*User, error) {
	logger.Infof("req1:%#v", user)
	rsp := User{user.ID, user.Name, user.Age, time.Now()}
	logger.Infof("rsp1:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) QueryUsers(_ context.Context, users []*User) (*UserResponse, error) {
	return &UserResponse{
		Users: users,
	}, nil
}

func (u *UserProvider) QueryAll(_ context.Context) (*UserResponse, error) {
	users := []*User{
		{
			ID:   "001",
			Name: "Joe",
			Age:  18,
			Time: time.Now(),
		},
		{
			ID:   "002",
			Name: "Wen",
			Age:  20,
			Time: time.Now(),
		},
	}

	return &UserResponse{
		Users: users,
	}, nil
}

func (u *UserProvider) MethodMapper() map[string]string {
	return map[string]string{
		"QueryUser":  "queryUser",
		"QueryUsers": "queryUsers",
		"QueryAll":   "queryAll",
	}
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}
