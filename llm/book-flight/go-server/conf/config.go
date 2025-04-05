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
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var (
	configPrompts CfgPrompts
	congifEnv     Environment
	oncePrompts   sync.Once
	onceEnv       sync.Once
)

// Config structure matches the YAML file structure
type Environment struct {
	Model   string
	Url     string
	ApiKey  string
	TimeOut int
}

// loadConfigPrompts reads and parses YAML file
func loadEnvironment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Reading environment variables
	congifEnv.Model = os.Getenv("LLM_MODEL")               //
	congifEnv.Url = os.Getenv("LLM_URL")                   // Default: http://localhost:11434
	congifEnv.ApiKey = os.Getenv("LLM_API_KEY")            //
	val, err := strconv.Atoi(os.Getenv("TIME_OUT_SECOND")) // Default: 300
	if err != nil {
		congifEnv.TimeOut = 300
	} else {
		congifEnv.TimeOut = val
	}
}

func GetEnvironment() Environment {
	onceEnv.Do(loadEnvironment)
	return congifEnv
}

// Config structure matches the YAML file structure
type CfgPrompts struct {
	ReactPrompt        string `yaml:"reactPrompt"`
	FinalPrompt        string `yaml:"finalPrompt"`
	IntentPrompt       string `yaml:"intentPrompt"`
	FormatInstructions string `yaml:"formatInstructions"`
}

// loadConfigPrompts reads and parses YAML file
func loadConfigPrompts() {
	data, err := os.ReadFile("go-server/conf/prompt_template.yml")
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
