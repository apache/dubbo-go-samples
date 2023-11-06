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

package test

import (
	"context"
	"testing"

	health "dubbo.apache.org/dubbo-go/v3/protocol/triple/health/triple_health"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	check, err := svc.Check(context.Background(), &health.HealthCheckRequest{Service: "greet.GreetService"})
	assert.Nil(t, err)
	assert.Equal(t, "status:SERVING", check.String())

	watch, err := svc.Watch(context.Background(), &health.HealthCheckRequest{Service: "greet.GreetService"})
	assert.Nil(t, err)
	assert.Equal(t, true, watch.Recv())
	assert.Equal(t, "status:SERVING", watch.Msg().String())
}
