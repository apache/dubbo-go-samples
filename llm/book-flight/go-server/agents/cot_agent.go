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

	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/actions"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/conf"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/model"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/model/ollama"
	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/tools"
)

type CotAgentRunner struct {
	llm                model.LLM
	tools              []tools.Tool
	maxThoughtSteps    int32
	reactPrompt        string
	finalPrompt        string
	formatInstructions string
}

func NewCotAgentRunner(llm model.LLM, tools []tools.Tool, maxSteps int32, cfgPrompt conf.CfgPrompts) CotAgentRunner {
	return CotAgentRunner{
		llm:                llm,
		tools:              tools,
		maxThoughtSteps:    maxSteps,
		reactPrompt:        cfgPrompt.ReactPrompt,
		finalPrompt:        cfgPrompt.FinalPrompt,
		formatInstructions: cfgPrompt.FormatInstructions,
	}
}

func (cot *CotAgentRunner) Run(ctx context.Context, input string, callopt model.Option) (string, error) {
	agentMemory := []map[string]any{}

	idxThoughtStep := 0
	for idxThoughtStep < int(cot.maxThoughtSteps) {
		action, response := cot.step(input, agentMemory, callopt)

		if action.Name == "FINISH" {
			break
		}

		observation := cot.execAction(action)
		agentMemory = cot.updateMemory(agentMemory, response, observation)

		idxThoughtStep++
	}

	reply := "Sorry, failed to complete your task."
	var err error = nil
	if idxThoughtStep < int(cot.maxThoughtSteps) {
		prompt := conf.Prompt(cot.finalPrompt, map[string]any{
			"task_description": input,
			"memory":           agentMemory},
			cot.tools,
		)
		reply, err = cot.llm.Call(context.Background(), prompt, callopt, ollama.WithTemperature(0.0))
	}

	return reply, err
}

func (cot *CotAgentRunner) step(taskDescription string, memory []map[string]any, callopt model.Option) (actions.Action, string) {
	prompt := conf.Prompt(
		cot.reactPrompt,
		map[string]any{"task_description": taskDescription, "memory": memory},
		cot.tools,
	)
	response, _ := cot.llm.Invoke(context.Background(), prompt, callopt, ollama.WithTemperature(0.0))

	response = model.RemoveThink(response)
	return actions.NewAction(response), response
}

func (cot *CotAgentRunner) execAction(action actions.Action) string {
	var err error = nil
	var observation string = fmt.Sprintf("Can't find tool: %v.", action.Name)
	for _, tool := range cot.tools {
		if tool.Name() == action.Name {
			strArgs, _ := json.Marshal(action.Args)
			observation, err = tool.Call(context.Background(), string(strArgs))
			if err != nil {
				observation = "Validation Error in args: " + string(strArgs)
			}
		}
	}
	return observation
}

func (cot *CotAgentRunner) updateMemory(memory []map[string]any, response string, observation string) []map[string]any {
	return append(memory,
		map[string]any{"input": response},
		map[string]any{"output": "\nResult:\n" + observation},
	)
}
