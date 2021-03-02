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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChinese(t *testing.T) {
	have, err := chinese.Have()
	assert.Nil(t, err)
	assert.Equal(t, "I'm Chinese and I have a Dog", have)
	hear, err := chinese.Hear()
	assert.Nil(t, err)
	assert.Equal(t, "I'm Chinese and I heard a Tiger yells like Tiger Tiger!", hear)
}

func TestAmerican(t *testing.T) {
	have, err := american.Have()
	assert.Nil(t, err)
	assert.Equal(t, "I'm American and I have a Cat", have)
	hear, err := american.Hear()
	assert.Nil(t, err)
	assert.Equal(t, "I'm American and I heard a Lion yells like Lion Lion!", hear)
}
