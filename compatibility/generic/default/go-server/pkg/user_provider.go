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

package pkg

import (
	"context"
	"strconv"
	"time"
)

import (
	"github.com/dubbogo/gost/log/logger"
)

type UserProvider struct {
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

func (u *UserProvider) MethodMapper(_ context.Context) map[string]string {
	return map[string]string{
		"QueryUser":  "queryUser",
		"QueryUsers": "queryUsers",
		"QueryAll":   "queryAll",
	}
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}
