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
	"dubbo.apache.org/dubbo-go/v3/config/generic"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/stretchr/testify/assert"
)

func TestDubboGetUser(t *testing.T) {
	o, err := dubboRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser",
		[]string{"java.lang.String", "java.lang.String"},
		[]hessian.Object{"A003", "Joe"},
	)
	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp := o.(map[interface{}]interface{})
	assert.Equal(t, "Joe", resp["name"])
	assert.Equal(t, int32(48), resp["age"])
	assert.Equal(t, "A003", resp["id"])
}

func TestTripleGetUser(t *testing.T) {
	o, err := tripleRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser",
		[]string{"java.lang.String", "java.lang.String"},
		[]hessian.Object{"A003", "Joe"},
	)
	assert.Nil(t, err)
	assert.IsType(t, make(map[interface{}]interface{}), o)
	resp := o.(map[interface{}]interface{})
	assert.Equal(t, "Joe", resp["name"])
	assert.Equal(t, int32(48), resp["age"])
	assert.Equal(t, "A003", resp["id"])
}
