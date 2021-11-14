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

func (u *UserProvider) GetUser(ctx context.Context, userID string) (*User, error) {
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

func (u *UserProvider) GetUser2(ctx context.Context, userID string) (*User, error) {
	var err error

	gxlog.CInfo("userID:%#v", userID)
	rsp := &User{
		ID:  userID,
		Sex: Gender(MAN).String(),
	}
	return rsp, err
}

func (s *UserProvider) MethodMapper() map[string]string {
	return map[string]string{
		"GetUser2": "getUser",
	}
}
