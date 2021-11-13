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

import (
	"dubbo.apache.org/dubbo-go/v3/config"

	"github.com/apache/dubbo-go-hessian2/java_exception"
)

func init() {
	config.SetProviderService(&UserProvider{})
}

type UserProvider struct {
}

func (u *UserProvider) GetUser(ctx context.Context, userID string) (*User, error) {
	if user, ok := userMap[userID]; ok {
		return &user, nil
	}

	return nil, java_exception.NewThrowable(fmt.Sprintf("invalid user id:%s", userID))
}

func (s *UserProvider) MethodMapper() map[string]string {
	return map[string]string{
		"GetUser2": "getUser",
	}
}
