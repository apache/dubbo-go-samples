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
	"time"
)

import (
	gxlog "github.com/dubbogo/gost/log"

	"github.com/opentracing/opentracing-go"
)

type (
	User struct {
		Id   string
		Name string
		Age  int32
		Time time.Time
	}
)

func (u User) String() string {
	return fmt.Sprintf(
		"User{ID:%s, Name:%s, Age:%d, Time:%s}",
		u.Id, u.Name, u.Age, u.Time,
	)
}

func (u *User) JavaClassName() string {
	return "org.apache.dubbo.User"
}

type UserProvider struct {
}

func (u *UserProvider) GetUser(ctx context.Context, req *User) (*User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "User-Server-Span")
	gxlog.CInfo("req:%#v", req)
	user := &User{Id: "001", Name: "zhangsan-dubbogo", Age: 18, Time: time.Now()}
	span.Finish()
	return user, nil
}
