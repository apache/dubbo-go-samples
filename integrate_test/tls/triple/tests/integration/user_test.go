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
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"github.com/dubbogo/gost/log/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGreeter(t *testing.T) {
	user, err := userProvider.GetUser(context.TODO(), &User{Name: "zlber"}, &User{Name: "zlber2"}, "testName")
	assert.Nil(t, err)
	assert.Equal(t, 18, user.Age)
	assert.Equal(t, "12345", user.ID)
	assert.Equal(t, "zlberzlber2", user.Name)

	logger.Infof("response result: %v\n", user)
}
