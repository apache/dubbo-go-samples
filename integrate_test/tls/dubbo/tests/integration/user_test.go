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

package integration

import (
	"context"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-samples/tls/dubbo/go-client/pkg"
)

func TestUser(t *testing.T) {
	reqUser := &pkg.User{}
	reqUser.ID = "003"
	user, err := userProvider.GetUser(context.TODO(), reqUser)
	assert.Nil(t, err)
	assert.NotNil(t, user)

	gender, err := userProvider.GetGender(context.TODO(), 1)
	assert.Nil(t, err)
	assert.NotNil(t, gender)

	ret, err := userProvider.GetUser0("003", "Moorse")
	assert.Nil(t, err)
	assert.NotNil(t, ret)

	ret1, err := userProvider.GetUsers([]string{"002", "003"})
	assert.Nil(t, err)
	assert.NotNil(t, ret1)

	var i int32 = 1
	user, err = userProvider.GetUser2(context.TODO(), i)
	assert.Nil(t, err)
	assert.NotNil(t, user)

	reqUser.ID = "003"
	_, err = userProvider.GetErr(context.TODO(), reqUser)
	assert.NotNil(t, err)

}
