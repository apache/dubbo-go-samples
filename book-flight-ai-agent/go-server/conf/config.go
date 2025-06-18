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

package conf

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

import (
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var (
	configPrompts CfgPrompts
	configEnv     Environment
	oncePrompts   sync.Once
	onceEnv       sync.Once
)

// Config structure matches the environment file structure
type Environment struct {
	Model      string `env:"LLM_MODEL"`
	Url        string `env:"LLM_URL"`
	ApiKey     string `env:"LLM_API_KEY"`
	HostClient string `env:"CLIENT_HOST"`
	PortClient int    `env:"CLIENT_PORT"`
	UrlClient  string `env:"_"`
	PortWeb    int    `env:"WEB_PORT"`
	TimeOut    int    `env:"TIMEOUT_SECONDS"`
}

// loadConfigPrompts reads and parses environment file
func loadEnvironment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Reading environment variables
	configEnv.Model = os.Getenv("LLM_MODEL")
	configEnv.Url = os.Getenv("LLM_URL")
	configEnv.ApiKey = os.Getenv("LLM_API_KEY")
	configEnv.HostClient = os.Getenv("CLIENT_HOST")
	configEnv.PortClient = AtoiWithDefault("CLIENT_PORT", 20000)
	configEnv.UrlClient = fmt.Sprintf("%s:%d", configEnv.HostClient, configEnv.PortClient)
	configEnv.PortWeb = AtoiWithDefault("WEB_PORT", 8080)
	configEnv.TimeOut = AtoiWithDefault("TIMEOUT_SECONDS", 300)
}

func GetEnvironment() Environment {
	onceEnv.Do(loadEnvironment)
	return configEnv
}

func AtoiWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}
	}
	return defaultValue
}

// Config structure matches the YAML file structure
type CfgPrompts struct {
	ReactPrompt        string `yaml:"reactPrompt"`
	FinalPrompt        string `yaml:"finalPrompt"`
	IntentPrompt       string `yaml:"intentPrompt"`
	InputPrompt        string `yaml:"inputPrompt"`
	UnrelatedPrompt    string `yaml:"unrelatedPrompt"`
	FormatInstructions string `yaml:"formatInstructions"`
}

// loadConfigPrompts reads and parses YAML file
func loadConfigPrompts() {
	data, err := os.ReadFile("go-server/conf/bookflight_prompt.yml")
	if err != nil {
		return
	}

	if err := yaml.Unmarshal(data, &configPrompts); err != nil {
		return
	}
}

func GetConfigPrompts() CfgPrompts {
	oncePrompts.Do(loadConfigPrompts)
	return configPrompts
}
