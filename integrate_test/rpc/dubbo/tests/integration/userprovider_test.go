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
	"github.com/apache/dubbo-go-samples/rpc/dubbo/go-server/pkg"
)

func TestGetUserA000(t *testing.T) {
	user, err := userProvider.GetUser(context.TODO(), "A000")
	assert.Nil(t, err)
	assert.Equal(t, "000", user.Id)
	assert.Equal(t, "Alex Stocks", user.Name)
	assert.Equal(t, int32(31), user.Age)
	assert.Equal(t, Gender(pkg.MAN), user.Sex)
	assert.NotNil(t, user.Time)
}

func TestGetUserA001(t *testing.T) {
	reqUser := &User{}
	reqUser.Id = "001"
	user, err := userProvider.GetUser(context.TODO(), "A001")
	assert.Nil(t, err)
	assert.Equal(t, "001", user.Id)
	assert.Equal(t, "ZhangSheng", user.Name)
	assert.Equal(t, int32(18), user.Age)
	assert.Equal(t, Gender(pkg.MAN), user.Sex)
	assert.NotNil(t, user.Time)
}

func TestGetUserA002(t *testing.T) {
	reqUser := &User{}
	reqUser.Id = "002"
	user, err := userProvider.GetUser(context.TODO(), "A002")
	assert.Nil(t, err)
	assert.Equal(t, "002", user.Id)
	assert.Equal(t, "Lily", user.Name)
	assert.Equal(t, int32(20), user.Age)
	assert.Equal(t, Gender(pkg.WOMAN), user.Sex)
	assert.NotNil(t, user.Time)
}

func TestGetUserA003(t *testing.T) {
	reqUser := &User{}
	reqUser.Id = "003"
	user, err := userProvider.GetUser(context.TODO(), "A003")
	assert.Nil(t, err)
	assert.Equal(t, "113", user.Id)
	assert.Equal(t, "Moorse", user.Name)
	assert.Equal(t, int32(30), user.Age)
	assert.Equal(t, Gender(pkg.WOMAN), user.Sex)
	assert.NotNil(t, user.Time)
}
