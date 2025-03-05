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
