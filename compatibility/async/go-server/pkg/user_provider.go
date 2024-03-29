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
	gxlog "github.com/dubbogo/gost/log"
)

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
