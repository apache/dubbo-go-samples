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

SOURCES = cmd/server.go
BINARY = user_info_server
PORT = 20000
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

ifeq ($(GOOS),windows)
	GOLANGCI_LINT := golangci-lint.exe
	export BINARY_EXT ?= .exe
else
	GOLANGCI_LINT := golangci-lint
	export BINARY_EXT ?=
endif

LOG := $(BINARY).log

CGO ?= 0

BASE_DIR := ./dist
ifeq ($(DEBUG), true)
	BUILD_TYPE := debug
	GCFLAGS := -gcflags="all=-N -l"
	LCFLAGS :=
else
	BUILD_TYPE := release
	LDFLAGS := "-s -w"
endif
OUT_DIR := $(BASE_DIR)/$(GOOS)_$(GOARCH)/$(BUILD_TYPE)

.PHONY: build
build: $(OUT_DIR)/$(BINARY)$(BINARY_EXT)

.PHONY: $(OUT_DIR)/$(BINARY)$(BINARY_EXT)
$(OUT_DIR)/$(BINARY)$(BINARY_EXT):
	$(info build $(OUT_DIR)/$(BINARY)$(BINARY_EXT))
	CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GCFLAGS) -ldflags=$(LDFLAGS) -i -o $(OUT_DIR)/$(BINARY)$(BINARY_EXT) $(SOURCES)

.PHONY: zookeeper-up
zookeeper-up:
	docker-compose -f docker/docker-compose.yml up -d

.PHONY: zookeeper-down
zookeeper-down:
	docker-compose -f docker/docker-compose.yml down

.PHONY: clean
clean:
	$(info about to clean up $(OUT_DIR))
	@rm -rf $(OUT_DIR)
	@-rm $(LOG)

.PHONY: start
start: build
	$(info start $(BINARY), output is redirected to $(LOG))
	@-$(OUT_DIR)/$(BINARY)$(BINARY_EXT) > $(LOG) 2>&1 &

.PHONY: stop
stop:
	$(eval PID := $(shell lsof -Pi :$(PORT) -sTCP:LISTEN -t))
	$(info stop $(BINARY) with PID $(PID))
	@-kill $(PID) 2>/dev/null || true

.PHONY: integration
integration:
	go clean -testcache && go test -tags integration -v ./tests/...





