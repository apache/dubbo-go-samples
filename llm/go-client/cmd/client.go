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
	chat "github.com/apache/dubbo-go-samples/llm/proto"
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
	resp = ""
	switch {
	case cmd == "/?" || cmd == "/help":
		resp += "Available commands:\n"
		resp += "/? help        - Show this help\n"
		resp += "/list          - List all contexts\n"
		resp += "/cd <context>  - Switch context\n"
		resp += "/new           - Create new context"
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
	default:
		resp += "Available commands:\n"
		resp += "/? help        - Show this help\n"
		resp += "/list          - List all contexts\n"
		resp += "/cd <context>  - Switch context\n"
		resp += "/new           - Create new context"
		return resp
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
			currentCtx.History = append(currentCtx.History,
				&chat.ChatMessage{
					Role:    "human",
					Content: input,
					Bin:     nil,
				})

			stream, err := svc.Chat(context.Background(), &chat.ChatRequest{
				Messages: currentCtx.History,
			})
			if err != nil {
				panic(err)
			}
			defer stream.Close()

			resp := ""

			for stream.Recv() {
				c := stream.Msg().Content
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
