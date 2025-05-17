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

package handlers

import (
	"context"
	"io"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"time"
)

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

import (
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-client/frontend/service"
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/conf"
	chat "github.com/apache/dubbo-go-samples/book-flight-ai-agent/proto"
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

	value := session.Get("current_context")
	if value == nil {
		value = h.ctxManager.CreateContext()
		session.Set("current_context", value)
		session.Save()
	}

	ctxID, _ := session.Get("current_context").(string)
	var req struct {
		Message string `json:"message"`
		Bin     string `json:"bin"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var img string
	if len(req.Bin) > 0 {
		re := regexp.MustCompile(`^data:image/([a-zA-Z]+);base64,([^"]+)$`)
		// this regex does not support file types like svg
		matches := re.FindStringSubmatch(req.Bin)

		if len(matches) != 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid base64 data format"})
			return
		}

		img = matches[2]
	}

	messages := h.ctxManager.GetHistory(ctxID)
	messages = append(messages, &chat.ChatMessage{
		Role:    "human",
		Content: req.Message,
		Bin:     []byte(img),
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

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "close")

	responseCh := make(chan string, 100) // use buffer
	responseRc := make(chan string, 100) // use buffer

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in stream processing: %v\n%s", r, debug.Stack())
			}
			close(responseCh)
			close(responseRc)
		}()

		for {
			select {
			case <-c.Request.Context().Done(): // client disconnect
				log.Println("Client disconnected, stopping stream processing")
				return
			default:
				if !stream.Recv() {
					if err := stream.Err(); err != nil {
						log.Printf("Stream receive error: %v", err)
					}
					return
				}
				content := stream.Msg().Content
				record := stream.Msg().Record
				if content != "" {
					responseCh <- content
				}
				if record != "" {
					responseRc <- record
				}
			}
		}
	}()

	// SSE stream output
	timeout := conf.GetEnvironment().TimeOut
	c.Stream(func(w io.Writer) bool {
		select {
		case chunk, ok := <-responseCh:
			if !ok {
				return false
			}
			c.SSEvent("message", gin.H{"content": chunk})
			return true
		case chunk, ok := <-responseRc:
			if !ok {
				return false
			}
			c.SSEvent("message", gin.H{"record": chunk})
			return true
		case <-time.After(time.Duration(timeout) * time.Second):
			log.Println("Stream time out")
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
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"context_id": newCtxID,
	})
}

func (h *ChatHandler) ListContexts(c *gin.Context) {
	session := sessions.Default(c)
	currentCtx := session.Get("current_context").(string)

	contexts := h.ctxManager.List()

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

	exists := h.ctxManager.Consists(req.ContextID)

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "context not found"})
		return
	}

	session := sessions.Default(c)
	session.Set("current_context", req.ContextID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "context switched",
	})
}
