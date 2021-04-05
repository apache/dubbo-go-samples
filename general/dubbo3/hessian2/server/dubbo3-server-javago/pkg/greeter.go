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

type User struct {
	Id   string
	Name string
	Age  int32
	//Time time.Time
}

type IGreeter2Impl struct {
}

func (u *IGreeter2Impl) SayHello0(ctx context.Context, request []interface{}) (*User, error) {
	req1 := request[0].(*User)
	req2 := request[1].(*User)
	fmt.Printf("get  %+v, %+v\n", req1, req2)
	rsp := &User{
		Name: req1.Name + req2.Name,
	}

	return rsp, nil
}

func (u *IGreeter2Impl) Reference() string {
	return "IGreeter2"
}

func (u User) JavaClassName() string {
	return "com.apache.dubbo.sample.basic.User"
}
