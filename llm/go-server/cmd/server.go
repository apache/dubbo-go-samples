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
	"encoding/base64"
	"fmt"
	"net/http"
	"runtime/debug"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"dubbo.apache.org/dubbo-go/v3/server"
	"github.com/dubbogo/gost/log/logger"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

import (
	"github.com/apache/dubbo-go-samples/llm/config"
	chat "github.com/apache/dubbo-go-samples/llm/proto"
)

var cfg *config.Config

type ChatServer struct {
	llm *ollama.LLM
}

func NewChatServer() (*ChatServer, error) {
	if cfg.ModelName == "" {
		return nil, fmt.Errorf("MODEL_NAME environment variable is not set")
	}

	llm, err := ollama.New(
		ollama.WithModel(cfg.ModelName),
		ollama.WithServerURL(cfg.OllamaURL),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize model %s: %v", cfg.ModelName, err)
	}
	logger.Infof("Initialized model: %s", cfg.ModelName)

	return &ChatServer{llm: llm}, nil
}

func (s *ChatServer) Chat(ctx context.Context, req *chat.ChatRequest, stream chat.ChatService_ChatServer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("panic in Chat: %v\n%s", r, debug.Stack())
			err = fmt.Errorf("internal server error")
		}
	}()

	if s.llm == nil {
		return fmt.Errorf("LLM model is not initialized")
	}

	if len(req.Messages) == 0 {
		logger.Info("Request contains no messages")
		return fmt.Errorf("empty messages in request")
	}

	var messages []llms.MessageContent
	for _, msg := range req.Messages {
		var msgType llms.ChatMessageType
		switch msg.Role {
		case "human":
			msgType = llms.ChatMessageTypeHuman
		case "ai":
			msgType = llms.ChatMessageTypeAI
		case "system":
			msgType = llms.ChatMessageTypeSystem
		}

		messageContent := llms.MessageContent{
			Role: msgType,
			Parts: []llms.ContentPart{
				llms.TextContent{Text: msg.Content},
			},
		}

		if msg.Bin != nil && len(msg.Bin) != 0 {
			decodeByte, err := base64.StdEncoding.DecodeString(string(msg.Bin))
			if err != nil {
				logger.Errorf("GenerateContent failed: %v\n", err)
				return fmt.Errorf("GenerateContent failed: %v", err)
			}
			imgType := http.DetectContentType(decodeByte)
			messageContent.Parts = append(messageContent.Parts, llms.BinaryPart(imgType, decodeByte))
		}

		messages = append(messages, messageContent)
	}

	_, err = s.llm.GenerateContent(
		ctx,
		messages,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if chunk == nil {
				return nil
			}
			return stream.Send(&chat.ChatResponse{
				Content: string(chunk),
				Model:   cfg.ModelName,
			})
		}),
	)
	if err != nil {
		logger.Errorf("GenerateContent failed with model %s: %v\n", cfg.ModelName, err)
		return fmt.Errorf("GenerateContent failed with model %s: %v", cfg.ModelName, err)
	}

	logger.Infof("GenerateContent successfully with model: %s", cfg.ModelName)

	return nil
}

func main() {
	var err error
	cfg, err = config.GetConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	srv, err := server.NewServer(
		server.WithServerRegistry(
			registry.WithNacos(),
			registry.WithAddress(cfg.NacosURL),
		),
		server.WithServerProtocol(
			protocol.WithTriple(),
			protocol.WithPort(cfg.ServerPort),
		),
	)

	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		return
	}

	chatServer, err := NewChatServer()
	if err != nil {
		fmt.Printf("Error creating chat server: %v\n", err)
		return
	}

	if err := chat.RegisterChatServiceHandler(srv, chatServer); err != nil {
		fmt.Printf("Error registering handler: %v\n", err)
		return
	}

	if err := srv.Serve(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
}
