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

func (s *ChatServer) Chat(ctx context.Context, req *chat.ChatRequest, stream chat.ChatService_ChatServer) error {
	callback := func(ctx context.Context, chunk []byte) error {
		return stream.Send(&chat.ChatResponse{
			Content: string(chunk),
		})
	}
	_, err := s.llm.GenerateContent(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, req.Prompt),
		},
		llms.WithStreamingFunc(callback),
	)
	return err
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
