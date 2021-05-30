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

package pkg

import (
	"context"
	"time"
)

import (
	"github.com/apache/dubbo-go/config"
	"github.com/dubbogo/gost/log"
)

func init() {
	config.SetProviderService(new(UserProvider))
}

type UserProvider struct {
}

func (u *UserProvider) GetUser(ctx context.Context, req []interface{}) (map[string]interface{}, error) {

	gxlog.CInfo("req:%#v", req)

	umap := make(map[string]interface{})
	umap["ID"] = "A001"
	umap["Name"] = "Alex Stocks"
	umap["Age"] = 18
	umap["Time"] = time.Now()

	gxlog.CInfo("rsp:%#v", umap)

	return umap, nil
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}
