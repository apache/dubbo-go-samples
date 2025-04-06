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
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

import (
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	"github.com/apache/dubbo-go-samples/llm/config"
	chat "github.com/apache/dubbo-go-samples/llm/proto"
)

type ChatContext struct {
	ID      string
	History []*chat.ChatMessage
}

var (
	contexts        = make(map[string]*ChatContext)
	currentCtxID    string
	contextOrder    []string
	maxID           uint8 = 0
	availableModels []string
	currentModel    string
)

func handleCommand(cmd string) (resp string) {
	cmd = strings.TrimSpace(cmd)
	resp = ""
	switch {
	case cmd == "/?" || cmd == "/help":
		resp += "Available commands:\n"
		resp += "/? help        - Show this help\n"
		resp += "/list          - List all contexts\n"
		resp += "/cd <context>  - Switch context\n"
		resp += "/new           - Create new context\n"
		resp += "/models        - List available models\n"
		resp += "/model <name>  - Switch to specified model"
		return resp
	case cmd == "/list":
		fmt.Println("Stored contexts (max 3):")
		for _, ctxID := range contextOrder {
			resp += fmt.Sprintf("- %s\n", ctxID)
		}
		return resp
	case strings.HasPrefix(cmd, "/cd "):
		target := strings.TrimPrefix(cmd, "/cd ")
		if ctx, exists := contexts[target]; exists {
			currentCtxID = ctx.ID
			resp += fmt.Sprintf("Switched to context: %s\n", target)
		} else {
			resp += "Context not found"
		}
		return resp
	case cmd == "/new":
		newID := createContext()
		currentCtxID = newID
		resp += fmt.Sprintf("Created new context: %s\n", newID)
		return resp
	case cmd == "/models":
		resp += "Available models:\n"
		for _, model := range availableModels {
			marker := " "
			if model == currentModel {
				marker = "*"
			}
			resp += fmt.Sprintf("\n%s %s", marker, model)
		}
		return resp
	case strings.HasPrefix(cmd, "/model "):
		modelName := strings.TrimPrefix(cmd, "/model ")
		modelFound := false
		for _, model := range availableModels {
			if model == modelName {
				currentModel = model
				modelFound = true
				break
			}
		}
		if modelFound {
			resp += fmt.Sprintf("Switched to model: %s", modelName)
		} else {
			resp += fmt.Sprintf("Model '%s' not found. Use /models to see available models.", modelName)
		}
		return resp
	default:
		return "Invalid command, use /? for help"
	}
}

func createContext() string {
	id := fmt.Sprintf("ctx%d", maxID)
	maxID++
	contexts[id] = &ChatContext{
		ID:      id,
		History: []*chat.ChatMessage{},
	}
	contextOrder = append(contextOrder, id)

	// up to 3 context
	if len(contextOrder) > 3 {
		delete(contexts, contextOrder[0])
		contextOrder = contextOrder[1:]
	}
	return id
}

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	availableModels = cfg.OllamaModels
	currentModel = cfg.DefaultModel()

	currentCtxID = createContext()

	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}

	svc, err := chat.NewChatService(cli)
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		return
	}

	fmt.Printf("\nSend a message (/? for help) - Using model: %s\n", currentModel)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n> ")
		scanner.Scan()
		input := scanner.Text()

		// handle command
		if strings.HasPrefix(input, "/") {
			fmt.Println(handleCommand(input))
			continue
		}

		func() {
			currentCtx := contexts[currentCtxID]
			currentCtx.History = append(currentCtx.History,
				&chat.ChatMessage{
					Role:    "human",
					Content: input,
					Bin:     nil,
				})

			stream, err := svc.Chat(context.Background(), &chat.ChatRequest{
				Messages: currentCtx.History,
				Model:    currentModel,
			})
			if err != nil {
				panic(err)
			}
			defer func(stream chat.ChatService_ChatClient) {
				err := stream.Close()
				if err != nil {
					fmt.Printf("Error closing stream: %v\n", err)
				}
			}(stream)

			resp := ""

			for stream.Recv() {
				msg := stream.Msg()
				c := msg.Content
				resp += c
				fmt.Print(c)
			}

			if err := stream.Err(); err != nil {
				fmt.Printf("Stream error: %v\n", err)
				return
			}

			currentCtx.History = append(currentCtx.History,
				&chat.ChatMessage{
					Role:    "ai",
					Content: resp,
					Bin:     nil,
				})
		}()
	}
}
