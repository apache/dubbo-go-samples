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

type UserProviderTriple struct{}

func (u *UserProviderTriple) GetUser(_ context.Context, userID string, name string) (*User, error) {
	logger.Infof("req:%#v, %#v", userID, name)
	rsp := User{userID, name, 48, time.Now()}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *UserProviderTriple) Reference() string {
	return "UserProviderTriple"
}
