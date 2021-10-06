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
	"fmt"
	"strconv"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"

	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go-hessian2/java_exception"

	"github.com/dubbogo/gost/log"

	perrors "github.com/pkg/errors"
)

func init() {
	config.SetProviderService(&UserProvider{})
}

type UserProvider struct {
}

func (u *UserProvider) getUser(userID string) (*User, error) {
	if user, ok := userMap[userID]; ok {
		return &user, nil
	}

	return nil, fmt.Errorf("invalid user id:%s", userID)
}

func (u *UserProvider) GetUser(ctx context.Context, req *User) (*User, error) {
	var (
		err  error
		user *User
	)

	gxlog.CInfo("req:%#v", req)
	user, err = u.getUser(req.ID)
	if err == nil {
		gxlog.CInfo("rsp:%#v", user)
	}
	return user, err
}

func (u *UserProvider) GetUser0(id string, name string) (User, error) {
	var err error

	gxlog.CInfo("id:%s, name:%s", id, name)
	user, err := u.getUser(id)
	if err != nil {
		return User{}, err
	}
	if user.Name != name {
		return User{}, perrors.New("name is not " + user.Name)
	}
	return *user, err
}

func (u *UserProvider) GetUser2(ctx context.Context, req int32) (*User, error) {
	var err error

	gxlog.CInfo("req:%#v", req)
	user := &User{}
	user.ID = strconv.Itoa(int(req))
	return user, err
}

func (u *UserProvider) GetUser3() error {
	return nil
}

func (u *UserProvider) GetErr(ctx context.Context, req *User) (*User, error) {
	return nil, java_exception.NewThrowable("exception")
}

func (u *UserProvider) GetUsers(req []string) ([]*User, error) {
	var err error

	gxlog.CInfo("req:%s", req)
	user, err := u.getUser(req[0])
	if err != nil {
		return nil, err
	}
	gxlog.CInfo("user:%v", user)
	user1, err := u.getUser(req[1])
	if err != nil {
		return nil, err
	}
	gxlog.CInfo("user1:%v", user1)

	return []*User{user, user1}, err
}

func (s *UserProvider) GetGender(i int32) (hessian.JavaEnum, error) {
	return hessian.JavaEnum(i), nil
}

func (s *UserProvider) MethodMapper() map[string]string {
	return map[string]string{
		"GetUser2": "getUser",
	}
}
