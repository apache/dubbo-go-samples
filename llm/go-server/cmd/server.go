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
	"os"
	"runtime/debug"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"

	"github.com/joho/godotenv"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

import (
	chat "github.com/apache/dubbo-go-samples/llm/proto"
)

type ChatServer struct {
	llm *ollama.LLM
}

func NewChatServer() (*ChatServer, error) {
	llm, err := ollama.New(
		ollama.WithModel(os.Getenv("OLLAMA_MODEL")),
		ollama.WithServerURL(os.Getenv("OLLAMA_URL")),
	)
	if err != nil {
		return nil, err
	}
	return &ChatServer{llm: llm}, nil
}

func (s *ChatServer) Chat(ctx context.Context, req *chat.ChatRequest, stream chat.ChatService_ChatServer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in Chat: %v\n%s", r, debug.Stack())
			err = fmt.Errorf("internal server error")
		}
	}()

	if s.llm == nil {
		return fmt.Errorf("LLM is not initialized")
	}

	if len(req.Messages) == 0 {
		log.Println("Request contains no messages")
		return fmt.Errorf("empty messages in request")
	}

	var messages []llms.MessageContent
	for _, msg := range req.Messages {
		msgType := llms.ChatMessageTypeHuman
		if msg.Role == "ai" {
			msgType = llms.ChatMessageTypeAI
		}

		messageContent := llms.MessageContent{
			Role: msgType,
			Parts: []llms.ContentPart{
				llms.TextContent{msg.Content},
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

	_, err = s.llm.GenerateContent(
		ctx,
		messages,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if chunk == nil {
				return nil
			}
			return stream.Send(&chat.ChatResponse{
				Content: string(chunk),
			})
		}),
	)
	if err != nil {
		log.Printf("GenerateContent failed: %v\n", err)
		return fmt.Errorf("GenerateContent failed: %v", err)
	}

	return nil
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}

	_, exist := os.LookupEnv("OLLAMA_MODEL")

	if !exist {
		fmt.Println("OLLAMA_MODEL is not set")
		return
	}

	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000),
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
