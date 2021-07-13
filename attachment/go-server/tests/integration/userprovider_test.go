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
	"github.com/apache/dubbo-go/common/constant"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	user := &User{}
	timestamp := time.Now()
	ctx := context.WithValue(context.Background(), constant.AttachmentKey, map[string]interface{}{"timestamp": timestamp})
	err := userProvider.GetUser(ctx, []interface{}{"A001"}, user)

	assert.Nil(t, err)
	assert.Equal(t, "A001", user.ID)
	assert.Equal(t, "Alex Stocks", user.Name)
	assert.Equal(t, int32(18), user.Age)
	assert.NotNil(t, user.Time)
	assert.True(t, user.Time.Before(timestamp))

	t.Logf("consumer timestamp: %v", timestamp)
	t.Logf("provider timestamp: %v", user.Time)
}
