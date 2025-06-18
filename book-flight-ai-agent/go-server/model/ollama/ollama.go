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
package ollama

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strings"
)

import (
	"github.com/ollama/ollama/api"
)

import (
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/model"
)

type (
	ImageData = api.ImageData
	LLMUrl    struct {
		scheam string // [scheam:]//host
		host   string // host:port
	}
)

func NewURL(url string) LLMUrl {
	scheam_host := strings.Split(url, "://")
	return LLMUrl{scheam: scheam_host[0], host: scheam_host[1]}
}

/*
Options:

	num_keep:
	seed:
	num_predict:
	top_k:
	top_p:
	min_p:
	typical_p:
	repeat_last_n:
	temperature:
	repeat_penalty:
	presence_penalty:
	frequency_penalty:
	mirostat:
	mirostat_tau:
	mirostat_eta:
	penalize_newline:
	stop:
	nama:
	num_ctx:
	num_batch:
	num_gpu:
	main_gpu:
	low_vram:
	vocab_only:
	vuse_mmap:
	use_mlock:
	num_thread:
*/
type LLMOllama struct {
	llmUrl    LLMUrl
	Model     string
	Url       string
	Prompt    string
	stream    *bool
	suffix    string
	images    []ImageData
	format    string
	system    string
	template  string
	raw       bool
	keepAlive string
	options   []any
}

func NewLLMOllama(model string, url string) *LLMOllama {
	return &LLMOllama{
		llmUrl: NewURL(url),
		Model:  model,
		Url:    url,
	}
}

func (llm *LLMOllama) Call(ctx context.Context, input string, opts ...model.Option) (string, error) {
	client := api.NewClient(&url.URL{Scheme: llm.llmUrl.scheam, Host: llm.llmUrl.host}, http.DefaultClient)

	optss := model.NewOptions(opts...)

	// By default, GenerateRequest is streaming.
	req := &api.GenerateRequest{
		Model:   llm.Model,
		Prompt:  input,
		Stream:  llm.stream,
		Suffix:  llm.suffix,
		Options: optss.Opts,
	}

	var respBuilder strings.Builder // Use strings.Builder
	respFunc := func(resp api.GenerateResponse) error {
		respBuilder.WriteString(resp.Response)
		return optss.CallOpt(resp.Response)
	}

	err := client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return respBuilder.String(), nil
}

func (llm *LLMOllama) Stream(ctx context.Context, input string, opts ...model.Option) (string, error) {
	return llm.Call(ctx, input, opts...)
}

func (llm *LLMOllama) Invoke(ctx context.Context, input string, opts ...model.Option) (string, error) {
	client := api.NewClient(&url.URL{Scheme: llm.llmUrl.scheam, Host: llm.llmUrl.host}, http.DefaultClient)

	// Messages
	messages := []api.Message{
		api.Message{
			Role:    "system",
			Content: llm.system,
		},
		api.Message{
			Role:    "user",
			Content: "Provide very brief, concise responses",
		},
		api.Message{
			Role:    "assistant",
			Content: "Provide very brief, concise responses",
		},
		api.Message{
			Role:    "user",
			Content: input,
			Images:  llm.images,
		},
	}

	// Options
	optss := model.NewOptions(opts...)

	// ChatRequest
	req := &api.ChatRequest{
		Model:    llm.Model,
		Stream:   llm.stream,
		Messages: messages,
		Options:  optss.Opts,
	}

	var respBuilder strings.Builder // Use strings.Builder
	respFunc := func(resp api.ChatResponse) error {
		respBuilder.WriteString(resp.Message.Content)
		if optss.CallOpt == nil {
			return nil
		}
		return optss.CallOpt(resp.Message.Content)
	}

	err := client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return respBuilder.String(), nil
}
