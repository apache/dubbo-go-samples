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
	"time"
)

type Gender int

const (
	MAN = iota
	WOMAN
)

var genderStrings = [...]string{
	"MAN",
	"WOMAN",
}

func (g Gender) String() string {
	return genderStrings[g]
}

type (
	User struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Age   int    `json:"age"`
		sex   Gender
		Birth int    `json:"time"`
		Sex   string `json:"sex"`
		Time  int64
	}
)

func (User) JavaClassName() string {
	return "org.apache.dubbo.User"
}

var (
	DefaultUser = User{
		ID: "0", Name: "Alex Stocks", Age: 31,
		// Birth: int(time.Date(1985, time.November, 10, 23, 0, 0, 0, time.UTC).Unix()),
		Birth: int(time.Date(1985, 11, 24, 15, 15, 0, 0, time.Local).Unix()),
		sex:   Gender(MAN),
	}

	UserMap = make(map[string]User)
)

func init() {
	DefaultUser.Sex = DefaultUser.sex.String()
	UserMap["A000"] = DefaultUser
	UserMap["A001"] = User{ID: "001", Name: "ZhangSheng", Age: 18, sex: MAN}
	UserMap["A002"] = User{ID: "002", Name: "Lily", Age: 20, sex: WOMAN}
	UserMap["A003"] = User{ID: "113", Name: "Moorse中文", Age: 30, sex: MAN}
	for k, v := range UserMap {
		v.Birth = int(time.Now().AddDate(-1*v.Age, 0, 0).Unix())
		v.Sex = UserMap[k].sex.String()
		UserMap[k] = v
	}
}
