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
	"os"

	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/agents"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/conf"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/model/ollama"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/tools"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/tools/bookingflight"
)

func main() {
	// test()
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getTools() ([]tools.Tool, error) {
	searchFlightTicketTool := bookingflight.NewSearchFlightTicket("查询机票", "查询指定日期可用的飞机票。")
	purchaseFlightTicketTool := bookingflight.NewPurchaseFlightTicket("购买机票", "购买飞机票。会返回购买结果(result), 和座位号(seat_number)")
	finishPlaceholder := bookingflight.NewFinishPlaceholder("FINISH", "用于表示任务完成的占位符工具")
	agentTools := []tools.Tool{
		searchFlightTicketTool,
		purchaseFlightTicketTool,
		finishPlaceholder,
	}
	return agentTools, nil
}

func run() error {
	cfgLLM := conf.GetConfigLLM()
	llm := ollama.NewLLMOllama(cfgLLM.LLM.Model, cfgLLM.LLM.Url)

	agentTools, _ := getTools()

	cfgPrompt := conf.GetConfigPrompts()
	cot := agents.NewCotAgentRunner(llm, agentTools, 10, cfgPrompt)
	question := "帮我买24年6月1日晚上北京到上海的飞机票"
	_, err := cot.Run(context.Background(), question, nil)
	return err
}
