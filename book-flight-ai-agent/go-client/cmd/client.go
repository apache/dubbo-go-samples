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
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/conf"
	chat "github.com/apache/dubbo-go-samples/book-flight-ai-agent/proto"
)

var cfgEnv = conf.GetEnvironment()

const (
	maxHistoryLength = 20 // Example maximum history length
	maxContextCount  = 3  // Example maximum number of contexts
)

type ChatContext struct {
	ID      string
	History []*chat.ChatMessage
}

var (
	contexts     = make(map[string]*ChatContext)
	currentCtxID string
	contextOrder []string
	maxID        uint8 = 0
)

func handleCommand(cmd string) (resp string) {
	cmd = strings.TrimSpace(cmd)

	var respBuilder strings.Builder // Use strings.Builder
	switch {
	case cmd == "/?" || cmd == "/help":
		respBuilder.WriteString("Available commands:\n")
		respBuilder.WriteString("/? help        - Show this help\n")
		respBuilder.WriteString("/list          - List all contexts\n")
		respBuilder.WriteString("/cd <context>  - Switch context\n")
		respBuilder.WriteString("/new           - Create new context")
	case cmd == "/list":
		fmt.Println("Stored contexts (max 3):")
		for _, ctxID := range contextOrder {
			respBuilder.WriteString(fmt.Sprintf("- %s\n", ctxID))
		}
	case strings.HasPrefix(cmd, "/cd "):
		target := strings.TrimPrefix(cmd, "/cd ")
		if ctx, exists := contexts[target]; exists {
			currentCtxID = ctx.ID
			respBuilder.WriteString(fmt.Sprintf("Switched to context: %s\n", target))
		} else {
			respBuilder.WriteString("Context not found")
		}
	case cmd == "/new":
		newID := createContext()
		currentCtxID = newID
		respBuilder.WriteString(fmt.Sprintf("Created new context: %s\n", newID))
	default:
		respBuilder.WriteString("Available commands:\n")
		respBuilder.WriteString("/? help        - Show this help\n")
		respBuilder.WriteString("/list          - List all contexts\n")
		respBuilder.WriteString("/cd <context>  - Switch context\n")
		respBuilder.WriteString("/new           - Create new context")
	}

	return respBuilder.String()
}

func createContext() string {
	id := fmt.Sprintf("ctx%d", maxID)
	maxID++
	contexts[id] = &ChatContext{
		ID:      id,
		History: []*chat.ChatMessage{},
	}
	contextOrder = append(contextOrder, id)

	if len(contextOrder) > maxContextCount {
		delete(contexts, contextOrder[0])
		contextOrder = contextOrder[1:]
	}
	return id
}

func addMessageToHistory(history []*chat.ChatMessage, newMessage *chat.ChatMessage) []*chat.ChatMessage {
	history = append(history, newMessage)
	if len(history) > maxHistoryLength {
		// Remove the oldest message
		history = history[1:]
	}
	return history
}

func main() {
	currentCtxID = createContext()

	cli, err := client.NewClient(
		client.WithClientURL(cfgEnv.UrlClient),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	svc, err := chat.NewChatService(cli)
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		return
	}

	fmt.Print("\nSend a message (/? for help)")
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
			currentCtx.History = addMessageToHistory(
				currentCtx.History,
				&chat.ChatMessage{
					Role:    "human",
					Content: input,
					Bin:     nil,
				})

			stream, err := svc.Chat(context.Background(), &chat.ChatRequest{
				Messages: currentCtx.History,
			})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			defer stream.Close()

			var respBuilder strings.Builder // Use strings.Builder
			for stream.Recv() {
				c := stream.Msg().Content
				respBuilder.WriteString(c) // Append to the builder
				fmt.Print(c)
			}
			resp := respBuilder.String() // Get the final string

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
