// +build integration

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
	. "github.com/apache/dubbo-go-samples/general/rest/go-server/pkg"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	// userProvider
	user := &User{}
	err := userProvider.GetUser(context.TODO(), []interface{}{"A001"}, user)
	assert.Nil(t, err)
	assert.Equal(t, "001", user.ID)
	assert.Equal(t, "ZhangSheng", user.Name)
	assert.Equal(t, int(18), user.Age)
	assert.NotNil(t, user.Time)

	// userProvider1
	user = &User{}
	err = userProvider1.GetUser(context.TODO(), []interface{}{"A003"}, user)
	assert.Nil(t, err)
	assert.Equal(t, "113", user.ID)
	assert.Equal(t, "Moorse中文", user.Name)
	assert.Equal(t, int(30), user.Age)
	assert.NotNil(t, user.Time)

	// userProvider2
	user = &User{}
	err = userProvider2.GetUser(context.TODO(), []interface{}{"A003"}, user)
	assert.Nil(t, err)
	assert.Equal(t, "113", user.ID)
	assert.Equal(t, "Moorse中文", user.Name)
	assert.Equal(t, int(30), user.Age)
	assert.NotNil(t, user.Time)
}

func TestGetUser0(t *testing.T) {
	// userProvider
	ret, err := userProvider.GetUser0("A003", "Moorse中文", 30)
	assert.Nil(t, err)
	assert.Equal(t, "113", ret.ID)
	assert.Equal(t, "Moorse中文", ret.Name)
	assert.Equal(t, int(30), ret.Age)
	assert.NotNil(t, ret.Time)

	// userProvider1
	ret1, err1 := userProvider1.GetUser0("A003", "Moorse中文", 30)
	assert.Nil(t, err1)
	assert.Equal(t, "113", ret1.ID)
	assert.Equal(t, "Moorse中文", ret1.Name)
	assert.Equal(t, int(30), ret1.Age)
	assert.NotNil(t, ret1.Time)

	// userProvider2
	ret2, err2 := userProvider2.GetUser0("A003", "Moorse中文", 30)
	assert.Nil(t, err2)
	assert.Equal(t, "113", ret2.ID)
	assert.Equal(t, "Moorse中文", ret2.Name)
	assert.Equal(t, int(30), ret2.Age)
	assert.NotNil(t, ret2.Time)
}

func TestGetUsers(t *testing.T) {
	// userProvider
	m := make(map[string]interface{})
	m["ID"]="A002"
	ret1, err := userProvider.GetUsers([]interface{}{m})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ret1))
	u := ret1[0]
	assert.Equal(t, "002", u.ID)
	assert.Equal(t, "Lily", u.Name)
	assert.Equal(t, int(20), u.Age)
	assert.NotNil(t, u.Time)

	// userProvider1
	m = make(map[string]interface{})
	m["ID"]="A002"
	ret1, err = userProvider1.GetUsers([]interface{}{m})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(ret1))

	// userProvider2
	m = make(map[string]interface{})
	m["ID"]="A002"
	ret1, err = userProvider1.GetUsers([]interface{}{m})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(ret1))
}


func TestGetUser3(t *testing.T) {
	// userProvider
	err := userProvider.GetUser3()
	assert.Nil(t, err)

	// userProvider1
	err = userProvider1.GetUser3()
	assert.Nil(t, err)

	// userProvider2
	err = userProvider2.GetUser3()
	assert.Nil(t, err)
}

func TestGetUser1(t *testing.T) {
	// userProvider
	_, err := userProvider.GetUser1([]interface{}{"A003"})
	assert.NotNil(t, err)

	// userProvider1
	_, err = userProvider1.GetUser1([]interface{}{"A003"})
	assert.NotNil(t, err)

	// userProvider2
	_, err = userProvider2.GetUser1([]interface{}{"A003"})
	assert.NotNil(t, err)
}
