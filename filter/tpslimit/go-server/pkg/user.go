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
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
)

type User struct {
	ID   string
	Name string
	Age  int32
	Time time.Time
}

type UserProvider struct {
}

func (u *UserProvider) GetUser(ctx context.Context, req []string) (*User, error) {
	logger.Infof("req:%#v", req)
	rsp := User{"A001", "Alex Stocks", 18, time.Now()}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *User) JavaClassName() string {
	return "org.apache.dubbo.User"
}
