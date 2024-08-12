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
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"

	"github.com/dubbogo/gost/log"

	perrors "github.com/pkg/errors"
)

func init() {
	config.SetProviderService(&UserProvider1{})
}

type UserProvider1 struct {
}

func (u *UserProvider1) getUser(userID string) (*User, error) {
	if user, ok := userMap[userID]; ok {
		return &user, nil
	}

	return nil, fmt.Errorf("invalid user id:%s", userID)
}

func (u *UserProvider1) GetUser(ctx context.Context, userID string) (*User, error) {
	var (
		err  error
		user *User
	)

	gxlog.CInfo("userID:%#v", userID)
	user, err = u.getUser(userID)
	if err == nil {
		gxlog.CInfo("rsp:%#v", user)
	}
	return user, err
}

func (u *UserProvider1) GetUser0(userID string, name string) (User, error) {
	var err error

	gxlog.CInfo("userID:%s, name:%s", userID, name)
	user, err := u.getUser(userID)
	if err != nil {
		return User{}, err
	}
	if user.Name != name {
		return User{}, perrors.New("name is not " + user.Name)
	}
	return *user, err
}

func (u *UserProvider1) GetUser2(ctx context.Context, userID string) (*User, error) {
	var err error

	gxlog.CInfo("userID:%#v", userID)
	rsp := &User{
		ID:  userID,
		Sex: MAN.String(),
	}
	return rsp, err
}

func (u *UserProvider1) GetUser3() error {
	return nil
}

func (u *UserProvider1) GetUsers(req []interface{}) ([]*User, error) {
	return []*User{}, nil
}

func (s *UserProvider1) MethodMapper() map[string]string {
	return map[string]string{
		"GetUser2": "getUser2",
	}
}
