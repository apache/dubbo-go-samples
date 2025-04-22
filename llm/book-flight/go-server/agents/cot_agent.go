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

package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/actions"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/conf"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/model"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/model/ollama"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/prompts"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/tools"
)

type CotAgentRunner struct {
	llm                model.LLM
	tools              []tools.Tool
	maxThoughtSteps    int32
	reactPrompt        string
	intentPrompt       string
	finalPrompt        string
	formatInstructions string
	memoryAgent        []map[string]any
	memoryMsg          []map[string]any
}

func NewCotAgentRunner(llm model.LLM, tools []tools.Tool, maxSteps int32, cfgPrompt conf.CfgPrompts) CotAgentRunner {
	return CotAgentRunner{
		llm:                llm,
		tools:              tools,
		maxThoughtSteps:    maxSteps,
		reactPrompt:        cfgPrompt.ReactPrompt,
		intentPrompt:       cfgPrompt.IntentPrompt,
		finalPrompt:        cfgPrompt.FinalPrompt,
		formatInstructions: cfgPrompt.FormatInstructions,
		memoryAgent:        []map[string]any{},
		memoryMsg:          []map[string]any{},
	}
}

func (cot *CotAgentRunner) Run(ctx context.Context, input string, callopt model.Option, callrst model.CallFunc) (string, error) {
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	opts := model.NewOptions(callopt)

	// Init Memory
	cot.memoryAgent = []map[string]any{}
	cot.memoryMsg = cot.updateMessage(cot.memoryMsg, input, "")

	var task string
	if len(cot.memoryMsg) > 0 {
		task = cot.summaryIntent(timeNow, callopt)
	} else {
		task = input
	}

	// Runner
	var response string
	var action actions.Action

	var idxThoughtStep int32
	var taskState TaskState
	for idxThoughtStep < cot.maxThoughtSteps {
		action, response = cot.thinkStep(task, timeNow, callopt, opts)
		taskState = InitTaskState(action.Name)

		if taskState == TaskCompleted || taskState == TaskUnrelated {
			break
		}

		observation := cot.execAction(action, opts)
		cot.memoryAgent = cot.updateMemory(cot.memoryAgent, response, observation)

		if taskState == TaskInputRequired {
			break
		}

		idxThoughtStep++
	}

	var err error
	reply := "Sorry, failed to complete your task."
	if idxThoughtStep < cot.maxThoughtSteps {
		reply, err = cot.finalStep(task, timeNow, callopt, callrst)

		cot.memoryMsg = cot.updateMessage(cot.memoryMsg, task, reply)
		if taskState == TaskCompleted {
			cot.memoryMsg = []map[string]any{}
		}
	}

	return reply, err
}

func (cot *CotAgentRunner) GetInputCtx(input string) string {
	var ctx string
	for _, msg := range cot.memoryAgent {
		if val, ok := msg["user"]; ok {
			ctx += fmt.Sprintf("\n%v", val)
		}
	}

	return strings.TrimSpace(ctx + "\n" + input)
}

func (cot *CotAgentRunner) summaryIntent(timeNow string, callopt model.Option) string {
	prompt := prompts.CreatePrompt(
		cot.intentPrompt,
		map[string]any{
			"content": cot.memoryMsg,
			"time":    timeNow,
		},
		nil,
	)
	response, _ := cot.llm.Call(context.Background(), prompt, callopt, ollama.WithTemperature(0.0))
	return model.RemoveThink(response)
}

func (cot *CotAgentRunner) thinkStep(
	task string,
	now string,
	callopt model.Option,
	opts model.Options) (actions.Action, string) {
	prompt := prompts.CreatePrompt(
		cot.reactPrompt,
		map[string]any{
			"task":   task,
			"memory": cot.memoryAgent,
			"time":   now,
		},
		cot.tools,
	)
	response, _ := cot.llm.Invoke(context.Background(), prompt, callopt, ollama.WithTemperature(0.0))
	opts.CallOpt("\n")
	response = model.RemoveThink(response)
	return actions.NewAction(response), response
}

func (cot *CotAgentRunner) finalStep(
	task string,
	date string,
	callopt model.Option,
	callrst model.CallFunc) (string, error) {
	prompt := prompts.CreatePrompt(
		cot.finalPrompt,
		map[string]any{
			"task":   task,
			"memory": cot.memoryAgent,
			"time":   date},
		cot.tools,
	)

	reply, err := cot.llm.Call(context.Background(), prompt, callopt, ollama.WithTemperature(0.0))
	reply = model.RemoveThink(reply)

	callrst(reply)
	return reply, err
}

func (cot *CotAgentRunner) execAction(action actions.Action, opts model.Options) string {
	var err error
	var observation string = fmt.Sprintf("Can't find tool: %v.", action.Name)
	for _, tool := range cot.tools {
		if tool.Name() == action.Name {
			strArgs, _ := json.Marshal(action.Args)
			observation, err = tool.Call(context.Background(), string(strArgs))
			opts.CallOpt("\n")
			if err != nil {
				observation = "Validation Error in args: " + string(strArgs)
			}
		}
	}
	return observation
}

func (cot *CotAgentRunner) updateMemory(memory []map[string]any, response string, observation string) []map[string]any {
	return append(memory,
		map[string]any{"input": response, "output": "\nResult:\n" + observation},
	)
}

func (cot *CotAgentRunner) updateMessage(memory []map[string]any, msgUser string, msgAgent string) []map[string]any {
	if msgUser != "" {
		memory = append(memory, map[string]any{"user": msgUser})
	}

	if msgAgent != "" {
		memory = append(memory, map[string]any{"agent": msgAgent})
	}

	return memory
}
