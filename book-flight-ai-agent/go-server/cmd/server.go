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
	"strings"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
)

import (
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/agents"
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/conf"
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/model/ollama"
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/tools"
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/tools/bookingflight"
	chat "github.com/apache/dubbo-go-samples/book-flight-ai-agent/proto"
)

var cfgEnv = conf.GetEnvironment()

func getTools() tools.Tools {
	var t tools.Tool
	var err error
	var tool_list []tools.Tool

	t, err = tools.CreateTool[bookingflight.SearchFlightTicketTool]("查询机票", "查询指定日期可用的飞机票。", "")
	if err == nil {
		tool_list = append(tool_list, t)
	}
	t, err = tools.CreateTool[bookingflight.PurchaseFlightTicketTool](
		"购买机票", "购买飞机票。会返回购买结果(result), 和座位号(seat_number)", "")
	if err == nil {
		tool_list = append(tool_list, t)
	}

	return agents.CreateToolkit(
		"订机票工具包，查询/预订机票功能。",
		agents.TaskCompleted|agents.TaskInputRequired|agents.TaskFailed|agents.TaskUnrelated,
		tool_list,
	)
}

type ChatServer struct {
	llm *ollama.LLMOllama
	cot agents.CotAgentRunner
}

func NewChatServer() (*ChatServer, error) {
	llm := ollama.NewLLMOllama(cfgEnv.Model, cfgEnv.Url)
	cot := agents.NewCotAgentRunner(llm, getTools(), 10, conf.GetConfigPrompts())
	return &ChatServer{llm: llm, cot: cot}, nil
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

	respFunc := func(resp string) error {
		return stream.Send(&chat.ChatResponse{
			Record: resp,
		})
	}

	rstFunc := func(resp string) error {
		return stream.Send(&chat.ChatResponse{
			Content: "\n" + strings.TrimSpace(resp),
		})
	}

	input := req.Messages[len(req.Messages)-1].Content
	_, err = s.cot.Run(ctx, input, ollama.WithStreamingFunc(respFunc), rstFunc)
	if err != nil {
		log.Printf("Run failed: %v", err)
	}

	return nil
}

func main() {
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(cfgEnv.PortClient),
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
