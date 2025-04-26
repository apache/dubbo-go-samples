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

package tools

import (
	"reflect"
	"strings"
)

type BaseTool struct {
	id             string
	name           string
	description    string
	requestParams  string
	responseParams string
	introduction   string
}

func NewBaseTool(name string, description string, introduction, id string) BaseTool {
	return BaseTool{
		name:         name,
		description:  description,
		introduction: introduction,
		id:           id,
	}
}

func (b BaseTool) Id() string             { return b.id }
func (b BaseTool) Name() string           { return b.name }
func (b BaseTool) Description() string    { return b.requestParams + " - " + b.description }
func (b BaseTool) RequestParams() string  { return b.requestParams }
func (b BaseTool) ResponseParams() string { return b.responseParams }
func (b BaseTool) Introduction() string   { return b.introduction }

func GetStructKeys(obj interface{}) string {
	if obj == nil {
		return "()"
	}

	t := reflect.TypeOf(obj)
	keys := []string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		jsonKey := strings.Split(jsonTag, ",")[0]

		// If the JSON tag is empty, use the structure field name
		if jsonKey == "" {
			jsonKey = field.Name
		}

		keys = append(keys, jsonKey)
	}

	rst := "(" + strings.Join(keys, ", ") + ")"
	return rst
}

// Toolkit is the manager of the toolkit, mainly providing descriptions of
// the tools and detailed descriptions of the toolkit.
type Toolkit struct {
	tools       []Tool
	toolMap     map[string]*Tool
	description string
}

func NewToolkit(tools []Tool, description string) Toolkit {
	toolkit := Toolkit{
		tools: tools,
		toolMap: func() map[string]*Tool {
			toolMap := make(map[string]*Tool)
			for _, tool := range tools {
				// 直接使用当前的 tool 指针
				toolMap[tool.Name()] = &tool
			}
			return toolMap
		}(),
		description: description,
	}

	return toolkit
}

// Return Toolkit Description
func (t Toolkit) Description() string {
	return t.description
}

// Returns description of all tools
func (t Toolkit) ToolsDescription() string {
	description := ""
	for _, tool := range t.tools {
		description += tool.Name() + tool.Description() + "\n"
	}
	return description
}

// Query Tool in the Toolkit
func (t Toolkit) QueryTool(method string) Tool {
	value, ok := t.toolMap[method]
	if !ok {
		return nil
	}
	return *value
}
