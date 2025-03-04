package handlers

import (
	"context"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

import (
	"github.com/apache/dubbo-go-samples/llm/frontend/service"
	chat "github.com/apache/dubbo-go-samples/llm/proto"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	svc        chat.ChatService
	ctxManager *service.ContextManager
}

func NewChatHandler(svc chat.ChatService, mgr *service.ContextManager) *ChatHandler {
	return &ChatHandler{
		svc:        svc,
		ctxManager: mgr,
	}
}

func (h *ChatHandler) Index(c *gin.Context) {
	session := sessions.Default(c)
	ctxID := session.Get("current_context")
	if ctxID == nil {
		ctxID = h.ctxManager.CreateContext()
		session.Set("current_context", ctxID)
		session.Save()
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "LLM Chat",
	})
}

func (h *ChatHandler) Chat(c *gin.Context) {
	session := sessions.Default(c)
	ctxID, ok := session.Get("current_context").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session context"})
		return
	}

	var req struct {
		Message string `json:"message"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	messages := h.ctxManager.GetHistory(ctxID)
	messages = append(messages, &chat.ChatMessage{
		Role:    "human",
		Content: req.Message,
	})

	stream, err := h.svc.Chat(context.Background(), &chat.ChatRequest{
		Messages: messages,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		if err := stream.Close(); err != nil {
			log.Println("Error closing stream:", err)
		}
	}()

	// 设置响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "close")

	responseCh := make(chan string)
	//defer close(responseCh)

	// 流处理协程
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in stream processing: %v\n%s", r, debug.Stack())
			}
			close(responseCh)
		}()

		var fullResponse strings.Builder
		for stream.Recv() {
			content := stream.Msg().Content
			responseCh <- content
			fullResponse.WriteString(content)
		}
	}()

	// SSE stream output
	c.Stream(func(w io.Writer) bool {
		select {
		case chunk, ok := <-responseCh:
			if !ok {
				return false
			}
			c.SSEvent("message", gin.H{"content": chunk})
			return true
		case <-time.After(30 * time.Second):
			log.Println("Stream timed out")
			return false
		case <-c.Request.Context().Done():
			log.Println("Client disconnected")
			return false
		}
	})
}

func (h *ChatHandler) NewContext(c *gin.Context) {
	session := sessions.Default(c)
	newCtxID := h.ctxManager.CreateContext()
	session.Set("current_context", newCtxID)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"context_id": newCtxID,
	})
}

func (h *ChatHandler) ListContexts(c *gin.Context) {
	session := sessions.Default(c)
	currentCtx := session.Get("current_context").(string)

	h.ctxManager.Mu.RLock()
	defer h.ctxManager.Mu.RUnlock()

	contexts := make([]string, 0, len(h.ctxManager.Contexts))
	for ctxID := range h.ctxManager.Contexts {
		contexts = append(contexts, ctxID)
	}

	c.JSON(http.StatusOK, gin.H{
		"current":  currentCtx,
		"contexts": contexts,
	})
}

func (h *ChatHandler) SwitchContext(c *gin.Context) {
	var req struct {
		ContextID string `json:"context_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.ctxManager.Mu.RLock()
	defer h.ctxManager.Mu.RUnlock()

	if _, exists := h.ctxManager.Contexts[req.ContextID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "context not found"})
		return
	}

	session := sessions.Default(c)
	session.Set("current_context", req.ContextID)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "context switched",
	})
}
