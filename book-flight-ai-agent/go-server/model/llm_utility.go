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

package model

import "regexp"

func RemoveThink(text string) string {
	re := regexp.MustCompile(`<think>[\s\S]*</think>`)
	matches := re.FindAllString(text, -1)
	if len(matches) > 0 {
		text = text[len(matches[0]):]
	}
	return text
}

func MergeMaps[K comparable, V any](m1, m2 map[K]V) map[K]V {
	merged := make(map[K]V)

	for k, v := range m1 {
		merged[k] = v
	}

	for k, v := range m2 {
		merged[k] = v
	}

	return merged
}
