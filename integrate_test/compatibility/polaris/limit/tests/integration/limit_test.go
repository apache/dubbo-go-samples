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
	"github.com/stretchr/testify/assert"
)

func TestPolarisLimit(t *testing.T) {

	var successCount, failCount int64
	for i := 0; i < 10; i++ {
		time.Sleep(50 * time.Millisecond)
		_, err := userProvider.GetUser(context.TODO(), &User{Name: "Alex03"})
		if err != nil {
			failCount++
		} else {
			successCount++
		}
	}
	//current limiting effect
	assert.Equal(t, true, failCount > 0)

}
