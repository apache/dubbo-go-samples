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
		ID   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
		sex  Gender
		Sex  string `json:"sex"`
	}
)

var (
	DefaultUser = User{
		ID: "0", Name: "Alex Stocks", Age: 31,
		sex: Gender(MAN),
	}

	userMap = make(map[string]User)
)

func init() {
	DefaultUser.Sex = DefaultUser.sex.String()
	userMap["A000"] = DefaultUser
	userMap["A001"] = User{ID: "001", Name: "ZhangSheng", Age: 18, sex: MAN}
	userMap["A002"] = User{ID: "002", Name: "Lily", Age: 20, sex: WOMAN}
	userMap["A003"] = User{ID: "113", Name: "Moorse", Age: 30, sex: MAN}
	for k, v := range userMap {
		v.Sex = userMap[k].sex.String()
		userMap[k] = v
	}
}
