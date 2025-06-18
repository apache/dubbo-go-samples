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

package actions

import (
	"encoding/json"
	"regexp"
)

import (
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/mcp"
)

type Action mcp.RequestRPC

func NewAction(text string) Action {
	re := regexp.MustCompile("```json[^`]*```")
	matches := re.FindAllString(text, -1)
	action := Action{}
	if len(matches) > 0 {
		match := matches[0][7 : len(matches[0])-3]
		json.Unmarshal([]byte(match), &action)
	}
	return action
}
