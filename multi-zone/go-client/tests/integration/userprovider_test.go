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
	"dubbo.apache.org/dubbo-go/v3/common/constant"

	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-samples/multi-zone/go-client/pkg"
)

func TestGetUser(t *testing.T) {
	ctx := context.WithValue(context.Background(), constant.REGISTRY_KEY+"."+constant.ZONE_FORCE_KEY, true)

	var hz, sh int
	user := &pkg.User{}
	for i := 0; i < 50; i++ {
		err := userProvider.GetUser(ctx, []interface{}{i}, user)
		if err != nil {
			panic(err)
		}
		if "dev-hz" == user.ID {
			hz++
		}
		if "dev-sh" == user.ID {
			sh++
		}
	}

	assert.Greater(t, hz, 0)
	assert.Greater(t, sh, 0)
	assert.Equal(t, 50, hz+sh)
}
