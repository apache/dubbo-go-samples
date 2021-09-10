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

func TestTest(t *testing.T)  {

	ctx := context.Background()

	// test Echo
	echo, err := userProvider.Echo(ctx, "Phil")
	assert.Nil(t, err)
	assert.Equal(t, "Phil", echo)

	// test GetUser
	user, err := userProvider.GetUser(ctx, "A003")
	assert.Nil(t, err)
	assert.Equal(t, "Moorse", user.Name)
	assert.Equal(t, int64(30), user.Age)
	assert.Equal(t, "MAN", user.Sex)

	user0, err := userProvider.GetUser0("A003", "Moorse")
	assert.Nil(t, err)
	assert.Equal(t, "Moorse", user0.Name)
	assert.Equal(t, int64(30), user0.Age)
	assert.Equal(t, "MAN", user0.Sex)

	users, err := userProvider.GetUsers([]interface{}{[]interface{}{"A002", "A003"}})
	assert.Nil(t, err)
	assert.Equal(t, "Lily", users[0].Name)
	assert.Equal(t, int64(20), users[0].Age)
	assert.Equal(t, "WOMAN", users[0].Sex)

}