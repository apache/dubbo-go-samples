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
package bailian

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

import (
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/model"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BailianRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

type BailianResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type LLMBailian struct {
	Model     string
	Url       string
	ApiKey    string
	MaxTokens int
	options   []any
}

func NewLLMBailian(model string, url string, apiKey string) *LLMBailian {
	return &LLMBailian{
		Model:     model,
		Url:       url,
		ApiKey:    apiKey,
		MaxTokens: 2048,
		options:   []any{},
	}
}

func (llm *LLMBailian) Call(ctx context.Context, input string, opts ...model.Option) (string, error) {
	return llm.Invoke(ctx, input, opts...)
}

func (llm *LLMBailian) Stream(ctx context.Context, input string, opts ...model.Option) (string, error) {
	// 百炼API的流式调用实现
	// 简化版本，实际上应该使用流式API
	return llm.Invoke(ctx, input, opts...)
}

func (llm *LLMBailian) Invoke(ctx context.Context, input string, opts ...model.Option) (string, error) {
	options := model.NewOptions(opts...)

	// 解析选项
	temperature := 0.7
	for _, opt := range llm.options {
		if temp, ok := opt.(WithTemperature); ok {
			temperature = float64(temp)
		}
	}

	// 构建请求体
	reqBody := BailianRequest{
		Model: llm.Model,
		Messages: []Message{
			{Role: "user", Content: input},
		},
		Temperature: temperature,
		MaxTokens:   llm.MaxTokens,
		Stream:      false,
	}

	// 将请求体转换为JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", llm.Url+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+llm.ApiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var bailianResp BailianResponse
	if err := json.Unmarshal(respBody, &bailianResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// 检查是否有结果
	if len(bailianResp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	// 获取结果
	result := bailianResp.Choices[0].Message.Content
	result = strings.TrimSpace(result)

	// 调用回调函数（如果有）
	if options.CallOpt != nil {
		options.CallOpt(result)
	}

	return result, nil
}