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
	"time"
)

import (
	"github.com/dubbogo/gost/log"
	"github.com/stretchr/testify/assert"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/rest"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
)

var survivalTimeout int = 10e9

func TestUserProvider(t *testing.T) {

	time.Sleep(3e9)
	gxlog.CInfo("test Userprovider")
	test(t)
	gxlog.CInfo("test Userprovider1")
	test1(t)
	gxlog.CInfo("test Userprovider2")
	test2(t)

}

func checkGetUserNil(user *User, err error, t *testing.T) {
	assert.NotNil(t, err)
	assert.Equal(t, "", user.ID)
	assert.Equal(t, "", user.Name)
	assert.Equal(t, int64(0), user.Age)
	assert.NotNil(t, user.Time)
}
func checkGetUser(user *User, err error, t *testing.T) {
	assert.Nil(t, err)
	assert.Equal(t, "113", user.ID)
	assert.Equal(t, "Moorse中文", user.Name)
	assert.Equal(t, int64(30), user.Age)
	assert.NotNil(t, user.Time)
}

func checkGetUsers(user *[]User, err error, t *testing.T) {
	if len(*user) != 0 {
		assert.Nil(t, err)
		assert.Equal(t, "002", (*user)[0].ID)
		assert.Equal(t, "Lily", (*user)[0].Name)
		assert.Equal(t, int64(20), (*user)[0].Age)
		assert.NotNil(t, (*user)[0].Time)
		if len(*user) == 2 {
			assert.Nil(t, err)
			assert.Equal(t, "113", (*user)[1].ID)
			assert.Equal(t, "Moorse", (*user)[1].Name)
			assert.Equal(t, int64(30), (*user)[1].Age)
			assert.NotNil(t, (*user)[1].Time)
		}
	}

}

func checkGetUser3(err error, t *testing.T) {
	assert.Nil(t, err)

}

func checkGetUser1(user *User, err error, t *testing.T) {
	assert.NotNil(t, err)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "", user.Name)
	assert.Equal(t, int64(0), user.Age)
	assert.NotNil(t, user.Time)
}

func checkGetUser2(user *User, err error, t *testing.T) {
	assert.NotNil(t, err)
	assert.Equal(t, "", user.ID)
	assert.Equal(t, "", user.Name)
	assert.Equal(t, int64(0), user.Age)
	assert.Equal(t, int64(0), user.Time)
}

func test(t *testing.T) {

	gxlog.CInfo("\n\n\nstart to test rest")
	user := &User{}
	err := userProvider.GetUser(context.TODO(), []interface{}{"A003"}, user)
	checkGetUser(user, err, t)
	gxlog.CInfo("\n\n\nstart to test rest - GetUser0")
	ret, err := userProvider.GetUser0("A003", "Moorse")
	checkGetUserNil(&ret, err, t)
	gxlog.CInfo("\n\n\nstart to test rest - GetUsers")
	ret1, err := userProvider.GetUsers([]interface{}{&User{ID: "A002"}})
	checkGetUsers(&ret1, err, t)
	user = &User{}
	err = userProvider.GetUser2(context.TODO(), []interface{}{1}, user)
	checkGetUser2(user, err, t)
	gxlog.CInfo("\n\n\nstart to test rest - GetUser3")
	err = userProvider.GetUser3()
	checkGetUser3(err, t)
}

func test1(t *testing.T) {

	gxlog.CInfo("\n\n\nstart to test1 rest")
	user := &User{}
	err := userProvider1.GetUser(context.TODO(), []interface{}{"A003"}, user)
	checkGetUser(user, err, t)
	gxlog.CInfo("\n\n\nstart to test1 rest - GetUser0")
	ret, err := userProvider1.GetUser0("A003", "Moorse")
	checkGetUserNil(&ret, err, t)
	gxlog.CInfo("\n\n\nstart to test1 rest - GetUsers")
	ret1, err := userProvider.GetUsers([]interface{}{&User{ID: "A002"}})
	checkGetUsers(&ret1, err, t)
	user = &User{}
	err = userProvider1.GetUser2(context.TODO(), []interface{}{1}, user)
	checkGetUser2(user, err, t)
	gxlog.CInfo("\n\n\nstart to test1 rest - GetUser3")
	err = userProvider1.GetUser3()
	checkGetUser3(err, t)
}

func test2(t *testing.T) {

	gxlog.CInfo("\n\n\nstart to test2 rest")
	user := &User{}
	err := userProvider2.GetUser(context.TODO(), []interface{}{"A003"}, user)
	checkGetUser(user, err, t)
	gxlog.CInfo("\n\n\nstart to test2 rest - GetUser0")
	ret, err := userProvider2.GetUser0("A003", "Moorse")
	checkGetUserNil(&ret, err, t)
	gxlog.CInfo("\n\n\nstart to test2 rest - GetUsers")
	ret1, err := userProvider.GetUsers([]interface{}{&User{ID: "A002"}})
	checkGetUsers(&ret1, err, t)
	user = &User{}
	err = userProvider2.GetUser2(context.TODO(), []interface{}{1}, user)
	checkGetUser2(user, err, t)
	gxlog.CInfo("\n\n\nstart to test2 rest - GetUser3")
	err = userProvider2.GetUser3()
	checkGetUser3(err, t)
}
