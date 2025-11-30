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

type ReactAgentRunner struct {
	llm             model.LLM
	tools           tools.Tools
	maxThoughtSteps int32
	reactPrompts    conf.CfgPrompts
	memoryAgent     []map[string]any
	memoryMsg       []map[string]any
}

func NewReactAgentRunner(
	llm model.LLM,
	tools tools.Tools,
	maxSteps int32,
	reactPrompts conf.CfgPrompts,
) ReactAgentRunner {
	return ReactAgentRunner{
		llm:             llm,
		tools:           tools,
		maxThoughtSteps: maxSteps,
		reactPrompts:    reactPrompts,
		memoryAgent:     []map[string]any{},
		memoryMsg:       []map[string]any{},
	}
}

func (react *ReactAgentRunner) Run(
	ctx context.Context,
	input string,
	callopt model.Option,
	callrst model.CallFunc,
) (string, error) {
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	opts := model.NewOptions(callopt)

	// Init Memory
	react.memoryAgent = []map[string]any{}
	react.memoryMsg = react.updateMessage(react.memoryMsg, input, "")

	var task string
	if len(react.memoryMsg) > 0 {
		task = react.summaryIntent(timeNow, callopt)
	} else {
		task = input
	}

	// ReAct Loop: Reasoning -> Acting -> Observation
	var response string
	var action actions.Action

	var idxThoughtStep int32
	var taskState TaskState
	for idxThoughtStep < react.maxThoughtSteps {
		// Reasoning: Think about what to do next
		action, response = react.reasoningStep(task, timeNow, callopt, opts)
		taskState = InitTaskState(action.Method)

		// Acting: Execute the action
		observation := react.execAction(action, opts)

		// Observation: Update memory with the result
		react.memoryAgent = react.updateMemory(react.memoryAgent, response, observation)

		if InterruptTask(taskState) {
			break
		}

		idxThoughtStep++
	}

	var err error
	reply := "Sorry, failed to complete your task."
	if idxThoughtStep < react.maxThoughtSteps {
		reply, err = react.finalStep(task, input, timeNow, taskState, callopt, callrst)

		react.memoryMsg = react.updateMessage(react.memoryMsg, task, reply)
		if taskState == TaskCompleted || taskState == TaskUnrelated {
			react.memoryMsg = []map[string]any{}
		}
	}

	return reply, err
}

func (react *ReactAgentRunner) GetInputCtx(input string) string {
	var respBuilder strings.Builder // Use strings.Builder
	for _, msg := range react.memoryAgent {
		if val, ok := msg["user"]; ok {
			respBuilder.WriteString(fmt.Sprintf("\n%v", val))
		}
	}
	respBuilder.WriteString(fmt.Sprintf("\n%v", input))

	return strings.TrimSpace(respBuilder.String())
}

func (react *ReactAgentRunner) summaryIntent(timeNow string, callopt model.Option) string {
	prompt := prompts.CreatePrompt(
		react.reactPrompts.IntentPrompt,
		map[string]any{
			"memory": react.memoryMsg,
			"time":   timeNow,
		},
	)
	response, _ := react.llm.Call(context.Background(), prompt, callopt, ollama.WithTemperature(0.0))
	return model.RemoveThink(response)
}

func (react *ReactAgentRunner) reasoningStep(
	task string,
	now string,
	callopt model.Option,
	opts model.Options,
) (actions.Action, string) {
	prompt := prompts.CreatePrompt(
		react.reactPrompts.ReactPrompt,
		map[string]any{
			"task":                task,
			"memory":              react.memoryAgent,
			"time":                now,
			"tools":               react.tools.ToolsDescription(),
			"format_instructions": conf.GetConfigPrompts().FormatInstructions,
		},
	)
	response, _ := react.llm.Invoke(context.Background(), prompt, callopt, ollama.WithTemperature(0.0))
	opts.CallOpt("\n")
	response = model.RemoveThink(response)
	return actions.NewAction(response), response
}

func (react *ReactAgentRunner) finalStep(
	task string,
	input string,
	date string,
	taskState TaskState,
	callopt model.Option,
	callrst model.CallFunc,
) (string, error) {
	config := map[string]any{"task": task}
	promptTemplate := react.reactPrompts.FinalPrompt
	switch taskState {
	case TaskUnrelated:
		promptTemplate = react.reactPrompts.UnrelatedPrompt
		config["task"] = input
	case TaskInputRequired:
		promptTemplate = react.reactPrompts.InputPrompt
		config["memory"] = react.memoryAgent
	default:
		config["memory"] = react.memoryAgent
		config["time"] = date
	}

	prompt := prompts.CreatePrompt(promptTemplate, config)
	reply, err := react.llm.Call(context.Background(), prompt, callopt, ollama.WithTemperature(0.0))
	reply = model.RemoveThink(reply)

	callrst(reply)
	return reply, err
}

func (react *ReactAgentRunner) execAction(action actions.Action, opts model.Options) string {
	var err error
	var observation string = fmt.Sprintf("Can't find tool: %v.", action.Method)
	tool := react.tools.QueryTool(action.Method)
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

func (react *ReactAgentRunner) updateMemory(memory []map[string]any, response string, observation string) []map[string]any {
	return append(memory,
		map[string]any{"input": response, "output": "\nResult:\n" + observation},
	)
}

func (react *ReactAgentRunner) updateMessage(memory []map[string]any, msgUser string, msgAgent string) []map[string]any {
	if msgUser != "" {
		memory = append(memory, map[string]any{"Human": msgUser})
	}

	if msgAgent != "" {
		memory = append(memory, map[string]any{"Agent": msgAgent})
	}

	return memory
}
