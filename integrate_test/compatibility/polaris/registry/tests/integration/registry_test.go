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

func TestPolarisRegistry(t *testing.T) {
	user, err := userProvider.GetUser(context.TODO(), &User{Name: "Alex001"})
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.ID, "A001")
	assert.Equal(t, user.Name, "Alex Stocks")
	assert.Equal(t, user.Age, 18)

	user, err = userProviderWithCustomRegistryGroupAndVersion.GetUser(context.TODO(), &User{Name: "Alex001"})
	assert.Nil(t, err)
	assert.Equal(t, user.Name, "Alex Stocks from UserProviderWithCustomGroupAndVersion")
	assert.Equal(t, user.Age, 18)
}
