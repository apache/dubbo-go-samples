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
	"fmt"
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

func NewBaseTool(name, description, requestParams, id string) BaseTool {
	return BaseTool{
		name:          name,
		description:   description,
		requestParams: requestParams,
		id:            id,
	}
}

func (b BaseTool) ID() string   { return b.id }
func (b BaseTool) Name() string { return b.name }
func (b BaseTool) Description() string {
	return b.name + b.RequestParams() + " - " + b.description + "\n"
}
func (b *BaseTool) RequestParams() string {
	if b.requestParams == "" {
		b.requestParams = scanStructKeys(b)
	}
	return b.requestParams
}
func (b *BaseTool) ResponseParams() string { return b.responseParams }
func (b BaseTool) Introduction() string    { return b.introduction }

func scanStructKeys(obj interface{}) string {
	if obj == nil {
		return "()"
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		// If it is a pointer, get the value it points to
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		fmt.Println("scanStructKeys: not a struct type")
		return "()"
	}

	keys := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		jsonKey := strings.Split(jsonTag, ",")[0]

		if field.Name == "BaseTool" {
			continue
		}

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
				// Use the current tool pointer directly
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
	var descBuilder strings.Builder // Use strings.Builder
	for _, tool := range t.tools {
		descBuilder.WriteString(tool.Description())
	}
	return descBuilder.String()
}

// Query Tool in the Toolkit
func (t Toolkit) QueryTool(method string) Tool {
	value, ok := t.toolMap[method]
	if !ok {
		return nil
	}
	return *value
}

// CreateTool
func CreateTool[T any](name, description, id string) (*T, error) {
	tool := new(T)
	base := NewBaseTool(name, description, scanStructKeys(tool), id)

	v := reflect.ValueOf(tool).Elem()
	field := v.FieldByName("BaseTool")
	if !field.IsValid() {
		return nil, fmt.Errorf("CreateTool: %T does not have an embedded BaseTool field", tool)
	}
	if !field.CanSet() {
		return nil, fmt.Errorf("CreateTool: cannot set BaseTool field on %T", tool)
	}

	field.Set(reflect.ValueOf(base))
	return tool, nil
}
