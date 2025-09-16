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

package main

import (
	"context"
	"fmt"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
)

func init() {
	extension.SetFilter("myServerFilter", NewMyServerFilter)
}

func NewMyServerFilter() filter.Filter {
	return &MyServerFilter{}
}

type MyServerFilter struct {
}

func (f *MyServerFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	fmt.Println("MyServerFilter Invoke is called, method Name = ", invocation.MethodName())
	fmt.Printf("request attachments = %s\n", invocation.Attachments())
	return invoker.Invoke(ctx, invocation)
}
func (f *MyServerFilter) OnResponse(ctx context.Context, result protocol.Result, invoker protocol.Invoker, protocol protocol.Invocation) protocol.Result {
	fmt.Println("MyServerFilter OnResponse is called")
	myAttachmentMap := make(map[string]any)
	myAttachmentMap["key1"] = "value1"
	myAttachmentMap["key2"] = []string{"value1", "value2"}
	result.SetAttachments(myAttachmentMap)
	return result
}
