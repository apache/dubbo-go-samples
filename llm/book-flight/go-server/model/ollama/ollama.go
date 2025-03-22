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
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/model"
	"github.com/ollama/ollama/api"
)

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
	model     string
	url       string
	scheam    string // [scheam:]//host
	host      string // host:port
	Prompt    string
	stream    *bool
	suffix    string
	images    []string
	format    string
	system    string
	template  string
	raw       bool
	keepAlive string
	options   []any
}

func NewLLMOllama(model string, url string) *LLMOllama {
	scheam_host := strings.Split(url, "://")
	return &LLMOllama{
		model:  model,
		url:    url,
		scheam: scheam_host[0],
		host:   scheam_host[1],
	}
}

func (llm *LLMOllama) Call(ctx context.Context, input string, opts ...model.Option) (string, error) {
	client := api.NewClient(&url.URL{Scheme: llm.scheam, Host: llm.host}, http.DefaultClient)

	optss := make(map[string]any)
	if len(opts) > 0 {
		for _, opt := range opts {
			for k, v := range opt {
				optss[k] = v
			}
		}
	}

	// By default, GenerateRequest is streaming.
	req := &api.GenerateRequest{
		Model:   llm.model,
		Prompt:  input,
		Stream:  llm.stream,
		Suffix:  llm.suffix,
		Options: optss,
	}

	response := ""
	// ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.

		// In streaming mode, responses are partial so we call fmt.Print (and not
		// Println) in order to avoid spurious newlines being introduced. The
		// model will insert its own newlines if it wants.
		fmt.Print(resp.Response)
		response += resp.Response
		return nil
	}

	err := client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	fmt.Println()
	return response, nil
}

func (llm *LLMOllama) Stream(ctx context.Context, input string, opts ...model.Option) (string, error) {
	return llm.Call(ctx, input, opts...)
}

func (llm *LLMOllama) Invoke(ctx context.Context, input string, opts ...model.Option) (string, error) {
	var stream = new(bool)
	llm.stream, stream = stream, llm.stream
	rst, err := llm.Call(ctx, input, opts...)
	llm.stream, stream = stream, llm.stream
	return rst, err
}
