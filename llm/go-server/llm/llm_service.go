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

package llm

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

// LLMProvider represents different LLM service providers
type LLMProvider string

const (
	ProviderOllama      LLMProvider = "ollama"
	ProviderOpenAI      LLMProvider = "openai"
	ProviderAnthropic   LLMProvider = "anthropic"
	ProviderAzureOpenAI LLMProvider = "azure-openai"
)

// LLMService wraps langchaingo LLM with provider information
type LLMService struct {
	llm      llms.Model
	provider LLMProvider
	model    string
}

// NewLLMService creates a new LLM service based on provider configuration
func NewLLMService(provider LLMProvider, model string, baseURL string, apiKey string) (*LLMService, error) {
	var llm llms.Model
	var err error

	switch provider {
	case ProviderOllama:
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		llm, err = ollama.New(
			ollama.WithModel(model),
			ollama.WithServerURL(baseURL),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Ollama model %s: %v", model, err)
		}

	case ProviderOpenAI:
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}
		if apiKey == "" {
			return nil, fmt.Errorf("OpenAI API key is required")
		}
		
		opts := []openai.Option{
			openai.WithModel(model),
			openai.WithToken(apiKey),
		}
		if baseURL != "" {
			opts = append(opts, openai.WithBaseURL(baseURL))
		}
		
		llm, err = openai.New(opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize OpenAI model %s: %v", model, err)
		}

	case ProviderAnthropic:
		if apiKey == "" {
			apiKey = os.Getenv("ANTHROPIC_API_KEY")
		}
		if apiKey == "" {
			return nil, fmt.Errorf("Anthropic API key is required")
		}
		
		opts := []anthropic.Option{
			anthropic.WithModel(model),
			anthropic.WithToken(apiKey),
		}
		if baseURL != "" {
			opts = append(opts, anthropic.WithBaseURL(baseURL))
		}
		
		llm, err = anthropic.New(opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Anthropic model %s: %v", model, err)
		}

	case ProviderAzureOpenAI:
		if apiKey == "" {
			apiKey = os.Getenv("AZURE_OPENAI_API_KEY")
		}
		if apiKey == "" {
			return nil, fmt.Errorf("Azure OpenAI API key is required")
		}
		
		// Azure OpenAI uses OpenAI client with specific configuration
		opts := []openai.Option{
			openai.WithModel(model),
			openai.WithToken(apiKey),
		}
		if baseURL != "" {
			opts = append(opts, openai.WithBaseURL(baseURL))
		}
		
		llm, err = openai.New(opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Azure OpenAI model %s: %v", model, err)
		}

	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", provider)
	}

	return &LLMService{
		llm:      llm,
		provider: provider,
		model:    model,
	}, nil
}

// GenerateContent generates content using the LLM with streaming support
func (s *LLMService) GenerateContent(ctx context.Context, messages []llms.MessageContent, callback func(ctx context.Context, chunk []byte) error) error {
	_, err := s.llm.GenerateContent(
		ctx,
		messages,
		llms.WithStreamingFunc(callback),
	)
	return err
}

// GetProvider returns the provider type
func (s *LLMService) GetProvider() LLMProvider {
	return s.provider
}

// GetModel returns the model name
func (s *LLMService) GetModel() string {
	return s.model
}

// GetSupportedProviders returns a list of supported providers
func GetSupportedProviders() []LLMProvider {
	return []LLMProvider{
		ProviderOllama,
		ProviderOpenAI,
		ProviderAnthropic,
		ProviderAzureOpenAI,
	}
}
