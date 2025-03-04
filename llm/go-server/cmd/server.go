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
	"log"
	"runtime/debug"
	"time"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"

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
	llm, err := ollama.New(ollama.WithModel("deepseek-r1:1.5b"))
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

	// timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if s.llm == nil {
		return fmt.Errorf("LLM is not initialized")
	}

	if len(req.Messages) == 0 {
		log.Println("Request contains no messages")
		return fmt.Errorf("empty messages in request")
	}

	var messages []llms.MessageContent
	for i, msg := range req.Messages {
		msgType := llms.ChatMessageTypeHuman
		if msg.Role == "ai" {
			msgType = llms.ChatMessageTypeAI
		}

		messageContent := llms.TextParts(msgType, msg.Content)
		if err != nil {
			log.Printf("Invalid message content at index %d: %v", i, err)
			return fmt.Errorf("invalid message content at index %d", i)
		}

		messages = append(messages, messageContent)
	}
	log.Printf("Messages constructed successfully: %+v", messages)

	buffer := make(chan []byte, 100) // 使用缓冲区
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in producer: %v\n%s", r, debug.Stack())
			}
			close(buffer)
		}()

		log.Println("Starting GenerateContent...")
		_, err := s.llm.GenerateContent(
			ctx,
			messages,
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				if ctx.Err() != nil {
					log.Println("Context canceled in callback")
					return ctx.Err()
				}

				if chunk == nil {
					log.Println("Received nil chunk in callback")
					return fmt.Errorf("nil chunk received")
				}

				select {
				case buffer <- chunk:
					return nil
				case <-ctx.Done():
					log.Println("Context canceled in producer")
					return ctx.Err()
				}
			}),
			llms.WithTemperature(0.7),
			llms.WithMaxTokens(500), // 限制生成长度
		)
		if err != nil {
			log.Printf("GenerateContent failed: %v", err)
		}
		log.Println("GenerateContent completed")
	}()

	// 从缓冲区读取数据并发送到客户端
	for chunk := range buffer {
		if ctx.Err() != nil {
			log.Println("Context canceled while sending chunks")
			break
		}

		err := stream.Send(&chat.ChatResponse{
			Content: string(chunk),
		})
		if err != nil {
			log.Printf("Failed to send chunk to stream: %v", err)
			break
		}
	}

	return nil
}

func main() {
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
