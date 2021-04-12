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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)
import (
	"github.com/apache/dubbo-go/config"
)

func TestGetUser(t *testing.T) {
	o, err := referenceConfig.GetRPCService().(*config.GenericService).Invoke(
		context.TODO(),
		[]interface{}{
			"GetUser",
			[]string{"java.lang.String"},
			[]interface{}{"A003"},
		},
	)

	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}, 0), o)
	resp := o.(map[interface{}]interface{})
	assert.Equal(t, "Alex Stocks", resp["name"])
	assert.Equal(t, int32(18), resp["age"])
	assert.Equal(t, "A001", resp["id"])
}

func TestQueryUser(t *testing.T) {
	user := User{
		ID:   "3213",
		Name: "panty",
		Age:  25,
		Time: time.Now(),
	}

	o, err := referenceConfig.GetRPCService().(*config.GenericService).Invoke(
		context.TODO(),
		[]interface{}{
			"queryUser",
			[]string{"org.apache.dubbo.User"},
			[]interface{}{user},
		},
	)

	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}, 0), o)
	resp := o.(map[interface{}]interface{})
	assert.Equal(t, "panty", resp["name"])
	assert.Equal(t, int32(25), resp["age"])
	assert.Equal(t, "3213", resp["id"])
}
