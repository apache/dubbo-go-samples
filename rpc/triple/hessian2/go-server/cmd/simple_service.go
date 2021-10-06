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
)

import (
	gxlog "github.com/dubbogo/gost/log"
)

type User struct {
	Id   string
	Name string
	Age  int32
}

func (u *User) JavaClassName() string {
	return "com.apache.dubbo.sample.basic.User"
}

type UserProvider struct {
}

func (u *UserProvider) GetUser(ctx context.Context, usr *User) (*User, error) {
	gxlog.CInfo("req:%#v", usr)
	rsp := User{"12345", "Hello " + usr.Name, 18}
	gxlog.CInfo("rsp:%#v", rsp)
	return &rsp, nil
}
