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

type UserProvider struct{}

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

func (u *UserProvider) GetUser3(_ context.Context, userCode int32) (*User, error) {
	logger.Infof("req:%#v", userCode)
	rsp := User{strconv.Itoa(int(userCode)), "Alex Stocks", 18, time.Now()}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) GetUser4(_ context.Context, userCode int32, name string) (*User, error) {
	logger.Infof("req:%#v, %#v", userCode, name)
	rsp := User{strconv.Itoa(int(userCode)), name, 18, time.Now()}
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

func (u *UserProvider) GetUsers(_ context.Context, userIdList []string) ([]*User, error) {
	logger.Infof("req:%#v", userIdList)
	var users []*User
	for _, i := range userIdList {
		if v, ok := userMap[i]; ok {
			users = append(users, v)
		} else {
			users = append(users, &User{ID: i, Name: "Unknown"})
		}
	}
	return users, nil
}

func (u *UserProvider) GetUsersMap(_ context.Context, userIdList []string) (map[string]*User, error) {
	logger.Infof("req:%#v", userIdList)
	var users = make(map[string]*User)
	for _, i := range userIdList {
		if v, ok := userMap[i]; ok {
			users[i] = v
		}
	}
	return users, nil
}

func (u *UserProvider) QueryUser(_ context.Context, user *User) (*User, error) {
	logger.Infof("req1:%#v", user)
	rsp := User{user.ID, user.Name, user.Age, time.Now()}
	logger.Infof("rsp1:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) QueryUsers(_ context.Context, users []*User) ([]*User, error) {
	return users, nil
}

func (u *UserProvider) QueryAll(_ context.Context) (map[string]*User, error) {
	users := map[string]*User{
		"001": {ID: "001", Name: "Joe", Age: 18, Time: time.Now()},
		"002": {ID: "002", Name: "Wen", Age: 20, Time: time.Now()},
	}
	return users, nil
}

func (u *UserProvider) Reference() string {
	return "org.apache.dubbo.samples.UserProvider"
}

func (u *UserProvider) MethodMapper(_ context.Context) map[string]string {
	return map[string]string{}
}

// Invoke handles generic call via $invoke method
// Parameters: methodName (string), types ([]string), args ([]interface{})
func (u *UserProvider) Invoke(ctx context.Context, methodName string, types []string, args []interface{}) (interface{}, error) {
	logger.Infof("Generic invoke: method=%s, types=%v, args=%v", methodName, types, args)

	switch methodName {
	case "GetUser1":
		return u.GetUser1(ctx, args[0].(string))
	case "GetUser2":
		return u.GetUser2(ctx, args[0].(string), args[1].(string))
	case "GetUser3":
		return u.GetUser3(ctx, args[0].(int32))
	case "GetUser4":
		return u.GetUser4(ctx, args[0].(int32), args[1].(string))
	case "GetOneUser":
		return u.GetOneUser(ctx)
	case "GetUsers":
		return u.GetUsers(ctx, args[0].([]string))
	case "GetUsersMap":
		return u.GetUsersMap(ctx, args[0].([]string))
	case "QueryUser":
		user := args[0].(*User)
		return u.QueryUser(ctx, user)
	case "QueryUsers":
		users := args[0].([]*User)
		return u.QueryUsers(ctx, users)
	case "QueryAll":
		return u.QueryAll(ctx)
	default:
		logger.Errorf("Unknown method: %s", methodName)
		return nil, nil
	}
}
