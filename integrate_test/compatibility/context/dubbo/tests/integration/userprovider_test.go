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

func TestGetUser(t *testing.T) {
	// add field that client wants to send
	atta := make(map[string]any)
	atta["string-value"] = "string-demo"
	atta["int-value"] = 1231242
	atta["user-defined-value"] = ContextContent{InterfaceName: "test.interface.name"}
	reqContext := context.WithValue(context.Background(), constant.DubboCtxKey("attachment"), atta)

	// invoke with reqContext
	rspContext, err := userProvider.GetContext(reqContext)

	// assert dubbo-go fields
	assert.Nil(t, err)
	assert.Equal(t, "org.apache.dubbo.UserProvider", rspContext.InterfaceName)
	assert.Equal(t, "org.apache.dubbo.UserProvider", rspContext.Path)
	assert.NotNil(t, rspContext.LocalAddr)
	assert.NotNil(t, rspContext.RemoteAddr)
	assert.NotNil(t, rspContext.DubboVersion)

	// assert user defined fields
	assert.Equal(t, "test.interface.name", rspContext.UserDefinedStrVal)
	assert.Equal(t, "string-demo", rspContext.CtxStrVal)
	assert.Equal(t, int64(1231242), rspContext.CtxIntVal)
}
