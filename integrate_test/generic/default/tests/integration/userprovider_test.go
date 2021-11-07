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

// nolint
import (
	"dubbo.apache.org/dubbo-go/v3/config/generic"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/stretchr/testify/assert"
)

func TestGetUser1(t *testing.T) {
	o, err := dubboRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser1",
		[]string{"java.lang.String"},
		[]hessian.Object{"A003"},
	)
	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp := o.(map[interface{}]interface{})
	assert.Equal(t, "Joe", resp["name"])
	assert.Equal(t, int32(48), resp["age"])
	assert.Equal(t, "A003", resp["iD"])

	o, err = tripleRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser1",
		[]string{"java.lang.String"},
		[]hessian.Object{"A003"},
	)
	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp = o.(map[interface{}]interface{})
	assert.Equal(t, "Joe", resp["name"])
	assert.Equal(t, int32(48), resp["age"])
	assert.Equal(t, "A003", resp["iD"])
}

func TestGetUser2(t *testing.T) {
	o, err := dubboRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser2",
		[]string{"java.lang.String", "java.lang.String"},
		[]hessian.Object{"A003", "lily"},
	)
	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp := o.(map[interface{}]interface{})
	assert.Equal(t, "lily", resp["name"])
	assert.Equal(t, int32(48), resp["age"])
	assert.Equal(t, "A003", resp["iD"])

	o, err = tripleRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser2",
		[]string{"java.lang.String", "java.lang.String"},
		[]hessian.Object{"A003", "lily"},
	)
	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp = o.(map[interface{}]interface{})
	assert.Equal(t, "lily", resp["name"])
	assert.Equal(t, int32(48), resp["age"])
	assert.Equal(t, "A003", resp["iD"])
}

func TestGetUser3(t *testing.T) {
	o, err := dubboRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser3",
		[]string{"int"},
		[]hessian.Object{1},
	)
	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp := o.(map[interface{}]interface{})
	assert.Equal(t, "Alex Stocks", resp["name"])
	assert.Equal(t, int32(18), resp["age"])
	assert.Equal(t, "1", resp["iD"])

	o, err = tripleRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser3",
		[]string{"int"},
		[]hessian.Object{1},
	)
	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp = o.(map[interface{}]interface{})
	assert.Equal(t, "Alex Stocks", resp["name"])
	assert.Equal(t, int32(18), resp["age"])
	assert.Equal(t, "1", resp["iD"])
}

func TestGetUser4(t *testing.T) {
	o, err := dubboRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser4",
		[]string{"int", "java.lang.String"},
		[]hessian.Object{1, "zhangsan"},
	)
	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp := o.(map[interface{}]interface{})
	assert.Equal(t, "zhangsan", resp["name"])
	assert.Equal(t, int32(18), resp["age"])
	assert.Equal(t, "1", resp["iD"])

	o, err = tripleRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser4",
		[]string{"int", "java.lang.String"},
		[]hessian.Object{1, "zhangsan"},
	)
	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp = o.(map[interface{}]interface{})
	assert.Equal(t, "zhangsan", resp["name"])
	assert.Equal(t, int32(18), resp["age"])
	assert.Equal(t, "1", resp["iD"])
}

// TODO: waiting for fixing the bug that pass empty array with basic types properly
//func TestGetOneUser(t *testing.T) {
//	o, err := dubboRefConf.GetRPCService().(*generic.GenericService).Invoke(
//		context.TODO(),
//		"GetOneUser",
//		[]string{},
//		[]hessian.Object{},
//	)
//	assert.Nil(t, err)
//	assert.IsType(t, make(map[interface{}]interface{}, 0), o)
//	resp := o.(map[interface{}]interface{})
//	assert.Equal(t, "xavierniu", resp["name"])
//	assert.Equal(t, int32(24), resp["age"])
//	assert.Equal(t, "1000", resp["iD"])
//
// TODO: Triple protocol test is required.
//}

func TestGetUsers(t *testing.T) {
	o, err := dubboRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUsers",
		[]string{"java.util.List"},
		[]hessian.Object{
			[]hessian.Object{
				"001", "002", "003", "004",
			},
		},
	)
	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp := o.(map[interface{}]interface{})
	assert.Equal(t, "other-zhangsan", resp["users"].([]interface{})[0].(map[interface{}]interface{})["name"])
	assert.Equal(t, "other-lisi", resp["users"].([]interface{})[1].(map[interface{}]interface{})["name"])
	assert.Equal(t, "other-lily", resp["users"].([]interface{})[2].(map[interface{}]interface{})["name"])
	assert.Equal(t, "other-lisa", resp["users"].([]interface{})[3].(map[interface{}]interface{})["name"])

	o, err = tripleRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUsers",
		[]string{"java.util.List"},
		[]hessian.Object{
			[]hessian.Object{
				"001", "002", "003", "004",
			},
		},
	)
	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp = o.(map[interface{}]interface{})
	assert.Equal(t, "other-zhangsan", resp["users"].([]interface{})[0].(map[interface{}]interface{})["name"])
	assert.Equal(t, "other-lisi", resp["users"].([]interface{})[1].(map[interface{}]interface{})["name"])
	assert.Equal(t, "other-lily", resp["users"].([]interface{})[2].(map[interface{}]interface{})["name"])
	assert.Equal(t, "other-lisa", resp["users"].([]interface{})[3].(map[interface{}]interface{})["name"])
}

func TestQueryUser(t *testing.T) {
	o, err := dubboRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"queryUser",
		[]string{"org.apache.dubbo.User"},
		[]hessian.Object{map[string]hessian.Object{
			"iD":   "3213",
			"name": "panty",
			"age":  25,
			"time": time.Now(),
		}},
	)

	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp := o.(map[interface{}]interface{})
	assert.Equal(t, "panty", resp["name"])
	assert.Equal(t, int32(25), resp["age"])
	assert.Equal(t, "3213", resp["iD"])

	o, err = tripleRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"queryUser",
		[]string{"org.apache.dubbo.User"},
		[]hessian.Object{map[string]hessian.Object{
			"iD":   "3213",
			"name": "panty",
			"age":  25,
			"time": time.Now(),
		}},
	)

	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp = o.(map[interface{}]interface{})
	assert.Equal(t, "panty", resp["name"])
	assert.Equal(t, int32(25), resp["age"])
	assert.Equal(t, "3213", resp["iD"])
}

func TestQueryUsers(t *testing.T) {
	o, err := dubboRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"queryUsers",
		[]string{"java.util.List"},
		[]hessian.Object{
			[]hessian.Object{
				map[string]hessian.Object{
					"id":    "3212",
					"name":  "XavierNiu",
					"age":   24,
					"time":  time.Now(),
					"class": "org.apache.dubbo.User",
				},
				map[string]hessian.Object{
					"iD":    "3213",
					"name":  "zhangsan",
					"age":   21,
					"time":  time.Now(),
					"class": "org.apache.dubbo.User",
				},
			},
		},
	)

	assert.Nil(t, err)
	resp, ok := o.(map[interface{}]interface{})
	assert.True(t, ok)
	assert.Equal(t, "XavierNiu", resp["users"].([]interface{})[0].(map[interface{}]interface{})["name"])
	assert.Equal(t, "zhangsan", resp["users"].([]interface{})[1].(map[interface{}]interface{})["name"])

	o, err = tripleRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"queryUsers",
		[]string{"org.apache.dubbo.User"},
		[]hessian.Object{
			[]hessian.Object{
				map[string]hessian.Object{
					"id":    "3212",
					"name":  "XavierNiu",
					"age":   24,
					"time":  time.Now(),
					"class": "org.apache.dubbo.User",
				},
				map[string]hessian.Object{
					"iD":    "3213",
					"name":  "zhangsan",
					"age":   21,
					"time":  time.Now(),
					"class": "org.apache.dubbo.User",
				},
			},
		},
	)

	assert.Nil(t, err)
	resp, ok = o.(map[interface{}]interface{})
	assert.True(t, ok)
	assert.Equal(t, "XavierNiu", resp["users"].([]interface{})[0].(map[interface{}]interface{})["name"])
	assert.Equal(t, "zhangsan", resp["users"].([]interface{})[1].(map[interface{}]interface{})["name"])
}

// TODO: Waiting for hessian-go bugfix
//func TestQueryAll(t *testing.T) {
//	o, err := dubboRefConf.GetRPCService().(*generic.GenericService).Invoke(
//		context.TODO(),
//			"queryAll",
//			[]hessian.Object{},
//			[]hessian.Object{},
//	)
//
//	assert.Nil(t, err)
//	assert.IsType(t, make(map[interface{}]interface{}, 0), o)
//	resp := o.(map[interface{}]interface{})
//	assert.Equal(t, "Joe", resp[0].(*pkg.User).Name)
//	assert.Equal(t, "Wen", resp[1].(*pkg.User).Name)
//}
