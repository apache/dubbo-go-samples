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

PROJECT_NAME = $(shell basename "$(PWD)")
PID = /tmp/.$(PROJECT_NAME).pid
PROJECT_DIR=$(shell pwd)
BASE_DIR := $(PROJECT_DIR)/dist


SOURCES = $(wildcard $(PROJECT_DIR)/cmd/*.go)

export GO111MODULE ?= on
export GOPROXY ?= https://goproxy.io,direct
export GOSUMDB ?= sum.golang.org
export GOARCH ?= amd64

export DOCKER_HOST_IP = $(shell ifconfig en0 | grep inet | grep -v inet6 | awk '{print $$2}')

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
help: $(realpath $(firstword $(MAKEFILE_LIST)))
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## build: Build application's binaries
.PHONY: build
build: $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME) config

.PHONY: $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME)
$(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME):
	$(info   >  Buiding application binary: $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME))
	@CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GCFLAGS) -ldflags=$(LDFLAGS) -i -o $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME) $(SOURCES)


## config: Setup config files
.PHONY: config
config:
	$(info   >  Setting up config files)
	@mkdir -p $(OUT_DIR)/conf
	@-test -f $(PROJECT_DIR)/conf/log.yml && cat $(PROJECT_DIR)/conf/log.yml | sed "s#\$$HOST_IP#$(DOCKER_HOST_IP)#g" > $(OUT_DIR)/conf/log.yml && echo "  > $(OUT_DIR)/conf/log.yml"
	@-test -f $(PROJECT_DIR)/conf/server.yml && cat $(PROJECT_DIR)/conf/server.yml | sed "s#\$$HOST_IP#$(DOCKER_HOST_IP)#g" > $(OUT_DIR)/conf/server.yml && echo "  > $(OUT_DIR)/conf/server.yml"
	@-test -f $(PROJECT_DIR)/conf/client.yml && cat $(PROJECT_DIR)/conf/client.yml | sed "s#\$$HOST_IP#$(DOCKER_HOST_IP)#g" > $(OUT_DIR)/conf/client.yml && echo "  > $(OUT_DIR)/conf/client.yml"
	@-test -f $(PROJECT_DIR)/conf/router_config.yml && cat $(PROJECT_DIR)/conf/router_config.yml | sed "s#\$$HOST_IP#$(DOCKER_HOST_IP)#g" > $(OUT_DIR)/conf/router_config.yml && echo "  > $(OUT_DIR)/conf/router_config.yml"

## docker-health-check: check services health on docker
.PHONY: docker-health-check
docker-health-check:
	$(info   >  run docker health check with $(PROJECT_DIR)/docker/docker-health-check.sh)
	@-test -f $(PROJECT_DIR)/docker/docker-health-check.sh && bash -f $(PROJECT_DIR)/docker/docker-health-check.sh

## docker-up: Shutdown dependency services on docker
.PHONY: docker-up
docker-up:
	$(info   >  Starting dependency services with $(PROJECT_DIR)/docker/docker-compose.yml)
	@docker-compose -f $(PROJECT_DIR)/docker/docker-compose.yml up -d

## docker-down: Shutdown dependency services on docker
.PHONY: docker-down
docker-down:
	$(info   >  Stopping dependency services with $(PROJECT_DIR)/docker/docker-compose.yml)
	@docker-compose -f $(PROJECT_DIR)/docker/docker-compose.yml down

## clean: Clean up the output and the binary of the application
.PHONY: clean
clean: stop
	$(info   >  Cleanning up $(OUT_DIR))
	@-rm -rf $(OUT_DIR)
	@-rm $(PID)

## start: Start the application (for server)
.PHONY: start
start: export CONF_PROVIDER_FILE_PATH ?= $(OUT_DIR)/conf/server.yml
start: build
	$(info   >  Starting application $(PROJECT_NAME), output is redirected to $(LOG_FILE))
	@-$(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME) > $(LOG_FILE) 2>&1 & echo $$! > $(PID)
	@cat $(PID) | sed "/^/s/^/  \>  PID: /"

## start: print application log (for server)
.PHONY: print-server-log
print-server-log:
	$(info   >  print server log with $(LOG_FILE))
	@-cat $(LOG_FILE)

## run: Run the application (for client)
.PHONY: run
run: export CONF_CONSUMER_FILE_PATH ?= $(OUT_DIR)/conf/client.yml
run: export CONF_ROUTER_FILE_PATH ?= $(OUT_DIR)/conf/router_config.yml
run: build
	$(info   >  Running application $(PROJECT_NAME), output is redirected to $(LOG_FILE))
	@-$(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME) 2>&1 | tee $(LOG_FILE)

## stop: Stop running the application (for server)
.PHONY: stop
stop:
	$(info   >  Stopping the application $(PROJECT_NAME))
	@cat $(PID) | sed "/^/s/^/  \>  Killing PID: /"
	@-kill `cat $(PID)` 2>/dev/null || true

## integration: Run integration test for this application
.PHONY: integration
integration: export CONF_CONSUMER_FILE_PATH ?= $(OUT_DIR)/conf/client.yml
integration: export CONF_ROUTER_FILE_PATH ?= $(OUT_DIR)/conf/router_config.yml
integration:
	$(info   >  Running integration test for application $(PROJECT_NAME))
	@go clean -testcache
	@go test -tags integration -v $(PROJECT_DIR)/tests/...
