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
	"log"
	"net/http"
	"runtime/debug"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

import (
	"github.com/apache/dubbo-go-samples/llm/config"
	chat "github.com/apache/dubbo-go-samples/llm/proto"
)

type ChatServer struct {
	llms map[string]*ollama.LLM
}

func NewChatServer() (*ChatServer, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("Error loading config: %v\n", err)
	}

	llmMap := make(map[string]*ollama.LLM)

	for _, model := range cfg.OllamaModels {
		if model == "" {
			continue
		}

		llm, err := ollama.New(
			ollama.WithModel(model),
			ollama.WithServerURL(cfg.OllamaURL),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize model %s: %v", model, err)
		}
		llmMap[model] = llm
		log.Printf("Initialized model: %s", model)
	}

	return &ChatServer{llms: llmMap}, nil
}

func (s *ChatServer) Chat(ctx context.Context, req *chat.ChatRequest, stream chat.ChatService_ChatServer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in Chat: %v\n%s", r, debug.Stack())
			err = fmt.Errorf("internal server error")
		}
	}()

	if len(s.llms) == 0 {
		return fmt.Errorf("no LLM models are initialized")
	}

	if len(req.Messages) == 0 {
		log.Println("Request contains no messages")
		return fmt.Errorf("empty messages in request")
	}

	modelName := req.Model
	var llm *ollama.LLM

	if modelName != "" {
		var ok bool
		llm, ok = s.llms[modelName]
		if !ok {
			return fmt.Errorf("requested model '%s' is not available", modelName)
		}
	} else {
		for name, l := range s.llms {
			modelName = name
			llm = l
			break
		}
		log.Printf("No model specified, using default model: %s", modelName)
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
				log.Printf("GenerateContent failed: %v\n", err)
				return fmt.Errorf("GenerateContent failed: %v", err)
			}
			imgType := http.DetectContentType(decodeByte)
			messageContent.Parts = append(messageContent.Parts, llms.BinaryPart(imgType, decodeByte))
		}

		messages = append(messages, messageContent)
	}

	_, err = llm.GenerateContent(
		ctx,
		messages,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if chunk == nil {
				return nil
			}
			return stream.Send(&chat.ChatResponse{
				Content: string(chunk),
				Model:   modelName,
			})
		}),
	)
	if err != nil {
		log.Printf("GenerateContent failed with model %s: %v\n", modelName, err)
		return fmt.Errorf("GenerateContent failed with model %s: %v", modelName, err)
	}

	return nil
}

func main() {

	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_llm_server"),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress("127.0.0.1:8848"),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20000),
		),
	)

	if err != nil {
		panic(err)
	}
	srv, err := ins.NewServer()

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
