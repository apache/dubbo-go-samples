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

package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

import (
	"github.com/joho/godotenv"
)

type Config struct {
	// LLM Provider Configuration
	LLMProvider string   // ollama, openai, anthropic, azure-openai
	LLMModels   []string // List of available models
	LLMBaseURL  string   // Base URL for LLM service
	LLMAPIKey   string   // API key for LLM service

	// Legacy Ollama fields (for backward compatibility)
	OllamaModels []string
	OllamaURL    string

	// Common Configuration
	TimeoutSeconds  int
	NacosURL        string
	MaxContextCount int
	ModelName       string
	ServerPort      int
}

var (
	config     *Config
	configOnce sync.Once
	configErr  error
)

const defaultMaxContextCount = 3 // Default to 3 for backward compatibility
const defaultTimeoutSeconds = 300

func Load(envFile string) (*Config, error) {
	configOnce.Do(func() {
		config = &Config{}
		err := godotenv.Load(envFile)
		if err != nil {
			configErr = fmt.Errorf("error loading .env file: %v", err)
			return
		}

		// Load LLM provider configuration
		llmProvider := os.Getenv("LLM_PROVIDER")
		if llmProvider == "" {
			// Fallback to legacy Ollama configuration
			llmProvider = "ollama"
		}
		config.LLMProvider = strings.ToLower(strings.TrimSpace(llmProvider))

		// Load models - try LLM_MODELS first, then fallback to OLLAMA_MODELS for backward compatibility
		modelsEnv := os.Getenv("LLM_MODELS")
		if modelsEnv == "" {
			// Backward compatibility: try OLLAMA_MODELS
			modelsEnv = os.Getenv("OLLAMA_MODELS")
			if modelsEnv == "" {
				configErr = fmt.Errorf("error: LLM_MODELS or OLLAMA_MODELS environment variable is not set")
				return
			}
		}

		modelsList := strings.Split(modelsEnv, ",")
		for i, model := range modelsList {
			modelsList[i] = strings.TrimSpace(model)
		}
		if len(modelsList) == 0 {
			configErr = fmt.Errorf("error: No models available")
			return
		}

		config.LLMModels = modelsList

		// For backward compatibility, also set OllamaModels
		if config.LLMProvider == "ollama" {
			config.OllamaModels = modelsList
		}

		modelName := os.Getenv("MODEL_NAME")
		if modelName == "" {
			configErr = fmt.Errorf("MODEL_NAME environment variable is not set")
			return
		}
		modelName = strings.TrimSpace(modelName)
		modelValid := false
		for _, m := range modelsList {
			if m == modelName {
				modelValid = true
				break
			}
		}
		if !modelValid {
			configErr = fmt.Errorf("specified model %s is not in the configured models list", modelName)
			return
		}
		config.ModelName = modelName

		portStr := os.Getenv("SERVER_PORT")
		if portStr == "" {
			configErr = fmt.Errorf("Error: SERVER_PORT environment variable is not set\n")
			return
		}
		config.ServerPort, err = strconv.Atoi(portStr)
		if err != nil {
			configErr = fmt.Errorf("Error converting SERVER_PORT to int: %v\n", err)
			return
		}

		// Load LLM base URL and API key
		llmBaseURL := os.Getenv("LLM_BASE_URL")
		llmAPIKey := os.Getenv("LLM_API_KEY")

		// For backward compatibility with Ollama
		ollamaURL := os.Getenv("OLLAMA_URL")
		if llmBaseURL == "" && ollamaURL != "" {
			// Use OLLAMA_URL as fallback for LLM_BASE_URL
			llmBaseURL = ollamaURL
		}

		// Set default URL for providers if not configured
		if llmBaseURL == "" && config.LLMProvider == "ollama" {
			llmBaseURL = "http://localhost:11434"
		}
		if llmBaseURL == "" && config.LLMProvider == "openai" {
			llmBaseURL = "https://api.openai.com/v1"
		}

		// Validate configuration based on provider
		if config.LLMProvider != "ollama" && config.LLMProvider != "openai" && llmBaseURL == "" {
			configErr = fmt.Errorf("LLM_BASE_URL is required for %s provider", config.LLMProvider)
			return
		}

		// Set URLs and API key
		config.LLMBaseURL = llmBaseURL
		config.LLMAPIKey = llmAPIKey

		// For backward compatibility, also set OllamaURL
		if config.LLMProvider == "ollama" {
			config.OllamaURL = llmBaseURL
		}

		timeoutStr := os.Getenv("TIME_OUT_SECOND")
		if timeoutStr == "" {
			config.TimeoutSeconds = defaultTimeoutSeconds
		} else {
			timeout, err := strconv.Atoi(timeoutStr)
			if err != nil {
				configErr = fmt.Errorf("invalid TIME_OUT_SECOND value: %v", err)
				return
			}
			config.TimeoutSeconds = timeout
		}

		nacosURL := os.Getenv("NACOS_URL")
		if nacosURL == "" {
			configErr = fmt.Errorf("NACOS_URL is not set")
			return
		}
		config.NacosURL = nacosURL

		maxContextStr := os.Getenv("MAX_CONTEXT_COUNT")
		if maxContextStr == "" {
			config.MaxContextCount = defaultMaxContextCount
		} else {
			maxContext, err := strconv.Atoi(maxContextStr)
			if err != nil {
				configErr = fmt.Errorf("invalid MAX_CONTEXT_COUNT value: %v", err)
				return
			}
			config.MaxContextCount = maxContext
		}
	})

	return config, configErr
}

func GetConfig() (*Config, error) {
	return Load(".env")
}

func (c *Config) DefaultModel() string {
	if len(c.OllamaModels) > 0 {
		return c.OllamaModels[0]
	}
	return ""
}
