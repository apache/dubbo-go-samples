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

package prompts

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/conf"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/tools"
)

func CreatePrompt(prompt string, ctx map[string]any, tools []tools.Tool) string {
	// ctx
	for k, v := range ctx {
		switch v.(type) {
		case map[string]any:
		case []map[string]any:
			vstr, _ := json.Marshal(v)
			prompt = strings.ReplaceAll(prompt, "{"+k+"}", string(vstr))
		default:
			prompt = strings.ReplaceAll(prompt, "{"+k+"}", fmt.Sprintln(v))
		}
	}
	// Tools
	tools_description := ""
	for _, tool := range tools {
		tools_description += tool.Name() + tool.Description()
	}
	prompt = strings.ReplaceAll(prompt, "{tools}", tools_description)
	// format_instructions
	configPrompts := conf.GetConfigPrompts()
	prompt = strings.ReplaceAll(prompt, "{format_instructions}", configPrompts.FormatInstructions)
	return prompt
}
