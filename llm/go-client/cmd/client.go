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
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	chat "github.com/apache/dubbo-go-samples/llm/proto"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	svc, err := chat.NewChatService(cli)
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		return
	}

	stream, err := svc.Chat(context.Background(), &chat.ChatRequest{
		Prompt: "Write a simple function to calculate fibonacci sequence in Go",
	})
	if err != nil {
		fmt.Printf("Error calling service: %v\n", err)
		return
	}
	defer stream.Close()

	for stream.Recv() {
		fmt.Print(stream.Msg().Content)
	}

	if err := stream.Err(); err != nil {
		fmt.Printf("Stream error: %v\n", err)
		return
	}
}
