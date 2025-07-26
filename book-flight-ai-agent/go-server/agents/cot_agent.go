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
)

import (
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/actions"
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/conf"
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/model"
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/model/ollama"
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/prompts"
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/tools"
)

type CotAgentRunner struct {
	llm             model.LLM
	tools           tools.Tools
	maxThoughtSteps int32
	cotPrompts      conf.CfgPrompts
	memoryAgent     []map[string]any
	memoryMsg       []map[string]any
}

func NewCotAgentRunner(
	llm model.LLM,
	tools tools.Tools,
	maxSteps int32,
	cotPrompts conf.CfgPrompts,
) CotAgentRunner {
	return CotAgentRunner{
		llm:             llm,
		tools:           tools,
		maxThoughtSteps: maxSteps,
		cotPrompts:      cotPrompts,
		memoryAgent:     []map[string]any{},
		memoryMsg:       []map[string]any{},
	}
}

func (cot *CotAgentRunner) Run(
	ctx context.Context,
	input string,
	callopt model.Option,
	callrst model.CallFunc,
) (string, error) {
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
		taskState = InitTaskState(action.Method)

		observation := cot.execAction(action, opts)
		cot.memoryAgent = cot.updateMemory(cot.memoryAgent, response, observation)

		if InterruptTask(taskState) {
			break
		}

		idxThoughtStep++
	}

	var err error
	reply := "Sorry, failed to complete your task."
	if idxThoughtStep < cot.maxThoughtSteps {
		reply, err = cot.finalStep(task, input, timeNow, taskState, callopt, callrst)

		cot.memoryMsg = cot.updateMessage(cot.memoryMsg, task, reply)
		if taskState == TaskCompleted || taskState == TaskUnrelated {
			cot.memoryMsg = []map[string]any{}
		}
	}

	return reply, err
}

func (cot *CotAgentRunner) GetInputCtx(input string) string {
	var respBuilder strings.Builder // Use strings.Builder
	for _, msg := range cot.memoryAgent {
		if val, ok := msg["user"]; ok {
			respBuilder.WriteString(fmt.Sprintf("\n%v", val))
		}
	}
	respBuilder.WriteString(fmt.Sprintf("\n%v", input))

	return strings.TrimSpace(respBuilder.String())
}

func (cot *CotAgentRunner) summaryIntent(timeNow string, callopt model.Option) string {
	prompt := prompts.CreatePrompt(
		cot.cotPrompts.IntentPrompt,
		map[string]any{
			"memory": cot.memoryMsg,
			"time":   timeNow,
		},
	)
	response, _ := cot.llm.Call(context.Background(), prompt, callopt, ollama.WithTemperature(0.0))
	return model.RemoveThink(response)
}

func (cot *CotAgentRunner) thinkStep(
	task string,
	now string,
	callopt model.Option,
	opts model.Options,
) (actions.Action, string) {
	prompt := prompts.CreatePrompt(
		cot.cotPrompts.ReactPrompt,
		map[string]any{
			"task":                task,
			"memory":              cot.memoryAgent,
			"time":                now,
			"tools":               cot.tools.ToolsDescription(),
			"format_instructions": conf.GetConfigPrompts().FormatInstructions,
		},
	)
	response, _ := cot.llm.Invoke(context.Background(), prompt, callopt, ollama.WithTemperature(0.0))
	opts.CallOpt("\n")
	response = model.RemoveThink(response)
	return actions.NewAction(response), response
}

func (cot *CotAgentRunner) finalStep(
	task string,
	input string,
	date string,
	taskState TaskState,
	callopt model.Option,
	callrst model.CallFunc,
) (string, error) {
	config := map[string]any{"task": task}
	promptTemplate := cot.cotPrompts.FinalPrompt
	switch taskState {
	case TaskUnrelated:
		promptTemplate = cot.cotPrompts.UnrelatedPrompt
		config["task"] = input
	case TaskInputRequired:
		promptTemplate = cot.cotPrompts.InputPrompt
		config["memory"] = cot.memoryAgent
	default:
		config["memory"] = cot.memoryAgent
		config["time"] = date
	}

	prompt := prompts.CreatePrompt(promptTemplate, config)
	reply, err := cot.llm.Call(context.Background(), prompt, callopt, ollama.WithTemperature(0.0))
	reply = model.RemoveThink(reply)

	callrst(reply)
	return reply, err
}

func (cot *CotAgentRunner) execAction(action actions.Action, opts model.Options) string {
	var err error
	var observation string = fmt.Sprintf("Can't find tool: %v.", action.Method)
	tool := cot.tools.QueryTool(action.Method)
	if tool != nil {
		strArgs, _ := json.Marshal(action.Params)
		observation, err = tool.Call(context.Background(), string(strArgs))
		opts.CallOpt("\n")
		if err != nil {
			observation = "Validation Error in args: " + string(strArgs)
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
		memory = append(memory, map[string]any{"Human": msgUser})
	}

	if msgAgent != "" {
		memory = append(memory, map[string]any{"Agent": msgAgent})
	}

	return memory
}
