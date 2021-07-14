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
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestCat(t *testing.T) {
	name, err := cat.GetName()
	assert.Nil(t, err)
	assert.Equal(t, "Cat", name)
	id, err := cat.GetId()
	assert.Nil(t, err)
	assert.Equal(t, 1, id)
}

func TestTiger(t *testing.T) {
	id, err := tiger.GetId()
	assert.Nil(t, err)
	assert.Equal(t, 3, id)
}

func TestDog(t *testing.T) {
	yell, err := dog.Yell()
	assert.Nil(t, err)
	assert.Equal(t, "Woof Woof!", yell)
}

func TestLion(t *testing.T) {
	name, err := lion.GetName()
	assert.Nil(t, err)
	assert.Equal(t, "Lion", name)
}
