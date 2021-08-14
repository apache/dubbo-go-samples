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
	"dubbo.apache.org/dubbo-go/v3/config"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	o, err := referenceConfig.GetRPCService().(*config.GenericService).Invoke(
		context.TODO(),
		[]interface{}{
			"GetUser",
			[]string{"java.lang.String"},
			[]hessian.Object{"A003"},
		},
	)
	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}, 0), o)
	resp := o.(map[interface{}]interface{})
	assert.Equal(t, "Alex Stocks", resp["name"])
	assert.Equal(t, int32(18), resp["age"])
	assert.Equal(t, "A001", resp["iD"])
}

func TestQueryUser(t *testing.T) {
	o, err := referenceConfig.GetRPCService().(*config.GenericService).Invoke(
		context.TODO(),
		[]interface{}{
			"queryUser",
			[]string{"org.apache.dubbo.User"},
			[]hessian.Object{map[string]hessian.Object{
				"iD":   "3213",
				"name": "panty",
				"age":  25,
				"time": time.Now(),
			}},
		},
	)

	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}, 0), o)
	resp := o.(map[interface{}]interface{})
	assert.Equal(t, "panty", resp["name"])
	assert.Equal(t, int32(25), resp["age"])
	assert.Equal(t, "3213", resp["iD"])
}

func TestQueryUsers(t *testing.T) {
	o, err := referenceConfig.GetRPCService().(*config.GenericService).Invoke(
		context.TODO(),
		[]interface{}{
			"QueryUsers",
			[]string{"java.lang.Array"},
			[]hessian.Object{
				[]hessian.Object{
					map[string]hessian.Object{
						"iD":   "3213",
						"name": "panty",
						"age":  25,
						"time": time.Now(),
					},
					map[string]hessian.Object{
						"iD":   "3212",
						"name": "XavierNiu",
						"age":  24,
						"time": time.Now().Add(4),
					},
				},
			},
		},
	)

	assert.Nil(t, err)
	resp, ok := o.(map[interface{}]interface{})
	assert.True(t, ok)
	users, ok := resp["users"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, "panty", users[0].(map[interface{}]interface{})["name"])
	assert.Equal(t, int32(25), users[0].(map[interface{}]interface{})["age"])
	assert.Equal(t, "3213", users[0].(map[interface{}]interface{})["iD"])
	assert.Equal(t, "XavierNiu", users[1].(map[interface{}]interface{})["name"])
	assert.Equal(t, int32(24), users[1].(map[interface{}]interface{})["age"])
	assert.Equal(t, "3212", users[1].(map[interface{}]interface{})["iD"])
}

func TestGetOneUser(t *testing.T) {
	o, err := referenceConfig.GetRPCService().(*config.GenericService).Invoke(
		context.TODO(),
		[]interface{}{
			"GetOneUser",
			[]string{"org.apache.dubbo.User"},
			[]hessian.Object{},
		},
	)

	assert.Nil(t, err)
	resp, ok := o.(map[interface{}]interface{})
	assert.True(t, ok)
	assert.Equal(t, "xavierniu", resp["name"])
	assert.Equal(t, int32(24), resp["age"])
	assert.Equal(t, "1000", resp["iD"])
}
