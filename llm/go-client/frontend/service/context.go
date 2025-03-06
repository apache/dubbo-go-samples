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

package service

import (
	chat "github.com/apache/dubbo-go-samples/llm/proto"
	"strconv"
	"sync"
)

type ContextManager struct {
	Contexts map[string][]*chat.ChatMessage
	Mu       sync.RWMutex
}

var nowID uint8 = 0

func NewContextManager() *ContextManager {
	return &ContextManager{
		Contexts: make(map[string][]*chat.ChatMessage),
	}
}

func (m *ContextManager) CreateContext() string {
	m.Mu.Lock()
	ctxID := nowID
	nowID++
	defer m.Mu.Unlock()
	m.Contexts[strconv.Itoa(int(ctxID))] = []*chat.ChatMessage{}
	return strconv.Itoa(int(ctxID))
}

func (m *ContextManager) GetHistory(ctxID string) []*chat.ChatMessage {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	return m.Contexts[ctxID]
}

func (m *ContextManager) AppendMessage(ctxID string, msg *chat.ChatMessage) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	if len(m.Contexts[ctxID]) >= 10 {
		m.Contexts[ctxID] = m.Contexts[ctxID][1:]
	}
	m.Contexts[ctxID] = append(m.Contexts[ctxID], msg)
}
