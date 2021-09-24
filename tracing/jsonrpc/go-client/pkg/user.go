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

type JsonRPCUser struct {
	Id    string    `json:"id"`
	Name  string    `json:"name"`
	Age   int       `json:"age"`
	Birth time.Time `json:"time"`
}

func (u JsonRPCUser) String() string {
	return fmt.Sprintf(
		"User{ID:%s, Name:%s, Age:%d, Time:%s}",
		u.Id, u.Name, u.Age, u.Birth,
	)
}

type UserProvider struct {
	GetUser func(ctx context.Context, req []interface{}) (*JsonRPCUser, error)
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}
