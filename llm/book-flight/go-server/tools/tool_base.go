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
	idTool       string
	name         string
	description  string
	inputFormat  string
	outputFormat string
	introduction string
}

func NewBaseTool(name, description, inputFormat, outputFormat, introduction, idTool string) BaseTool {
	return BaseTool{
		name:         name,
		description:  description,
		inputFormat:  inputFormat,
		outputFormat: outputFormat,
		introduction: introduction,
		idTool:       idTool,
	}
}

func (b BaseTool) IdTool() string       { return b.idTool }
func (b BaseTool) Name() string         { return b.name }
func (b BaseTool) Description() string  { return b.inputFormat + " - " + b.description }
func (b BaseTool) InputFormat() string  { return b.inputFormat }
func (b BaseTool) OutputFormat() string { return b.outputFormat }
func (b BaseTool) Introduction() string { return b.introduction }

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

		// 如果 JSON 标签为空，使用结构体字段名
		if jsonKey == "" {
			jsonKey = field.Name
		}

		keys = append(keys, jsonKey)
	}

	rst := "(" + strings.Join(keys, ", ") + ")"
	return rst
}
