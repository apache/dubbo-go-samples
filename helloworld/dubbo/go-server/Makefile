# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, softwarek
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SOURCES = $(wildcard cmd/*.go)

BASE_DIR := dist
PROJECT_NAME = $(shell basename "$(PWD)")
PID=/tmp/.$(PROJECTNAME).pid

export APP_LOG_CONF_FILE = $(shell pwd)/conf/log.yml
export CONF_PROVIDER_FILE_PATH = $(shell pwd)/conf/server.yml
export CONF_CONSUMER_FILE_PATH = $(shell pwd)/conf/client.yml

export GO111MODULE ?= on
export GOPROXY ?= https://goproxy.io,direct
export GOSUMDB ?= sum.golang.org
export GOARCH ?= amd64

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

## build: Build application's binaries
.PHONY: build
build: $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME)

.PHONY: $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME)
$(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME):
	$(info   >  Buiding application binary: $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME))
	@CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GCFLAGS) -ldflags=$(LDFLAGS) -i -o $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME) $(SOURCES)

## docker-up: Shutdown dependency services on docker
.PHONY: docker-up
docker-up:
	$(info   >  Starting dependency services with docker/docker-compose.yml)
	@docker-compose -f docker/docker-compose.yml up -d

## docker-down: Shutdown dependency services on docker
.PHONY: docker-down
docker-down:
	$(info   >  Stopping dependency services with docker/docker-compose.yml)
	@docker-compose -f docker/docker-compose.yml down

## clean: Clean up the output and the binary of the application
.PHONY: clean
clean: stop
	$(info   >  Cleanning up binaries and logs)
	@rm -rf $(OUT_DIR)

## start: Start to run the application
.PHONY: start
start: build
	$(info   >  Starting application $(PROJECT_NAME), output is redirected to $(LOG_FILE))
	@-$(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME) > $(LOG_FILE) 2>&1 & echo $$! > $(PID)
	@cat $(PID) | sed "/^/s/^/  \>  PID: /"

## stop: Stop running the application
.PHONY: stop
stop:
	$(info   >  Stopping the application $(PROJECT_NAME))
	@-kill `cat $(PID)` 2>/dev/null || true

## integration: Run integration test for this application
.PHONY: integration
integration:
	$(info   >  Running integration test for application $(PROJECT_NAME))
	@go clean -testcache
	@go test -tags integration -v ./tests/...

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo





