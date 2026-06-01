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

package api

import (
	"net/http"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	restconfig "dubbo.apache.org/dubbo-go/v3/protocol/rest/config"
)

const RestPath = "/api/v1/users/{userID}/greeting"

func InstallProviderRestConfig() {
	restconfig.SetRestProviderServiceConfigMap(map[string]*restconfig.RestServiceConfig{
		InterfaceName: newRestServiceConfig(constant.DefaultRestServer, ""),
	})
}

func InstallConsumerRestConfig() {
	restconfig.SetRestConsumerServiceConfigMap(map[string]*restconfig.RestServiceConfig{
		InterfaceName: newRestServiceConfig("", constant.DefaultRestClient),
	})
}

func newRestServiceConfig(serverType string, clientType string) *restconfig.RestServiceConfig {
	return &restconfig.RestServiceConfig{
		InterfaceName: InterfaceName,
		Server:        serverType,
		Client:        clientType,
		RestMethodConfigsMap: map[string]*restconfig.RestMethodConfig{
			MethodGetGreeting: {
				InterfaceName:  InterfaceName,
				MethodName:     MethodGetGreeting,
				Path:           RestPath,
				MethodType:     http.MethodPost,
				PathParamsMap:  map[int]string{0: "userID"},
				QueryParamsMap: map[int]string{1: "name"},
				HeadersMap:     map[int]string{2: "X-Trace-ID"},
				Body:           3,
				Consumes:       "application/json",
				Produces:       "application/json",
			},
		},
	}
}
