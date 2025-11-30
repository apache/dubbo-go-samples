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

package user

import (
	"time"
)

var (
	// UserMap stores user data
	UserMap = make(map[string]*User)
)

func init() {
	UserMap["000"] = &User{
		Id:   "000",
		Name: "Alex Stocks",
		Age:  31,
		Time: time.Now().Unix(),
		Sex:  Gender_MAN,
	}
	UserMap["001"] = &User{
		Id:   "001",
		Name: "ZhangSheng",
		Age:  18,
		Time: time.Now().Unix(),
		Sex:  Gender_MAN,
	}
	UserMap["002"] = &User{
		Id:   "002",
		Name: "Lily",
		Age:  20,
		Time: time.Now().Unix(),
		Sex:  Gender_WOMAN,
	}
	UserMap["003"] = &User{
		Id:   "003",
		Name: "Moorse",
		Age:  30,
		Time: time.Now().Unix(),
		Sex:  Gender_WOMAN,
	}
}

// GetUserByID returns user by id
func GetUserByID(id string) (*User, bool) {
	user, ok := UserMap[id]
	return user, ok
}
