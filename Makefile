# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT_DIR ?= $(CURDIR)
PROJECT_NAME ?= $(notdir $(abspath $(PROJECT_DIR)))
PID = /tmp/.$(PROJECT_NAME).pid
BASE_DIR := $(PROJECT_DIR)/go-server/dist

SOURCES = $(wildcard $(PROJECT_DIR)/go-server/cmd/*.go)
GO ?= go

# shell
SHELL = /bin/bash

export GO111MODULE ?= on
export GOSUMDB ?= sum.golang.org
export GOARCH ?= amd64
export GONOPROXY ?= **.gitee.com**

OS := $(shell uname)
ifeq ($(OS), Linux)
	export GOOS ?= linux
else ifeq ($(OS), Darwin)
	export GOOS ?= darwin
else
	export GOOS ?= windows
endif

ifeq ($(GOOS), windows)
	export EXT_NAME ?= .exe
else
	export EXT_NAME ?=
endif

CGO ?= 0
ifeq ($(DEBUG), true)
	BUILD_TYPE := debug
	GCFLAGS := -gcflags="all=-N -l"
	LCFLAGS :=
else
	BUILD_TYPE := release
	LDFLAGS := "-s -w"
endif

OUT_DIR := $(BASE_DIR)/$(GOOS)_$(GOARCH)/$(BUILD_TYPE)
LOG_FILE := $(OUT_DIR)/$(PROJECT_NAME).log

export APP_LOG_CONF_FILE ?= $(OUT_DIR)/conf/log.yml

.PHONY: all
all: help
.DEFAULT_GOAL := help
help: $(realpath $(firstword $(MAKEFILE_LIST)))
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## build: Build application's binaries
.PHONY: build
build: $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME)

$(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME): $(SOURCES)
	$(info   >  Building application binary: $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME))
	@mkdir -p $(OUT_DIR)
	@CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build $(GCFLAGS) -ldflags=$(LDFLAGS) -o $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME) $(SOURCES)

## start: Start the application (for server)
.PHONY: start
start: export DUBBO_GO_CONFIG_PATH ?= $(PROJECT_DIR)/go-server/conf/dubbogo.yml
start: build
	$(info   >  Starting application $(PROJECT_NAME), output is redirected to $(LOG_FILE))
	@ls $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME) >/dev/null
	@-pkill -f "$(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME)" 2>/dev/null || true
	@-command -v lsof >/dev/null 2>&1 && lsof -P -sTCP:LISTEN -tiTCP:20000 -iTCP:20001 -iTCP:20002 -iTCP:20022 -iTCP:4318 -iTCP:50051 | xargs -r kill -9 || true
	@sleep 1
	@-cd $(PROJECT_DIR) && $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME) > $(LOG_FILE) 2>&1 & echo $$! > $(PID)
	@sed 's/^/  \>  PID: /' $(PID)

## stop: Stop running the application (for server)
.PHONY: stop
stop:
	$(info   >  Stopping the application $(PROJECT_NAME))
	@-test -f $(PID) && sed 's/^/  \>  Killing PID: /' $(PID) || true
	@-test -f $(PID) && kill `cat $(PID)` 2>/dev/null || true
	@-pkill -f "$(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME)" 2>/dev/null || true
	@-rm -f $(PID)
	@-command -v lsof >/dev/null 2>&1 && lsof -P -sTCP:LISTEN -tiTCP:20000 -iTCP:20001 -iTCP:20002 -iTCP:20022 -iTCP:4318 -iTCP:50051 | xargs -r kill -9 || true
