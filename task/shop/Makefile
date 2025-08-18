#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Dubbo-Go Microservice Shop Project Makefile
# Provides one-click start, stop, build and other functions

.PHONY: help build clean start stop restart logs test proto check-deps install-deps

# Default target
.DEFAULT_GOAL := help

# Project information
PROJECT_NAME := dubbo-go-shop
GO_MIN_VERSION := 1.19
REGISTRY := nacos

# Service configuration
SERVICES := frontend user detail order comment
SERVICE_PORTS := 8080 20001 20002 20003 20004

# Color definitions
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

help: ## Show help information
	@echo "$(GREEN)Dubbo-Go Microservice Shop Project Makefile$(RESET)"
	@echo ""
	@echo "$(YELLOW)Available commands:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(BLUE)%-15s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

check-deps: ## Check dependencies
	@echo "$(BLUE)Checking Go version...$(RESET)"
	@go version || (echo "$(RED)Go environment not found$(RESET)" && exit 1)
	@echo "$(GREEN)Go environment check passed$(RESET)"

install-deps: check-deps ## Install project dependencies
	@echo "$(BLUE)Installing project dependencies...$(RESET)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies installation completed$(RESET)"

build: install-deps ## Build all services
	@echo "$(BLUE)Building all services...$(RESET)"
	@mkdir -p bin
	@echo "$(YELLOW)Building Frontend service...$(RESET)"
	@go build -o bin/frontend ./frontend/cmd/main.go
	@echo "$(YELLOW)Building User service...$(RESET)"
	@go build -o bin/user ./user/go-server/cmd/server.go
	@echo "$(YELLOW)Building Detail V1 service...$(RESET)"
	@go build -o bin/detail-v1 ./detail/server_v1/cmd/server.go
	@echo "$(YELLOW)Building Detail V2 service...$(RESET)"
	@go build -o bin/detail-v2 ./detail/server_v2/cmd/server.go
	@echo "$(YELLOW)Building Order V1 service...$(RESET)"
	@go build -o bin/order-v1 ./order/server_v1/cmd/server.go
	@echo "$(YELLOW)Building Order V2 service...$(RESET)"
	@go build -o bin/order-v2 ./order/server_v2/cmd/server.go
	@echo "$(YELLOW)Building Comment V1 service...$(RESET)"
	@go build -o bin/comment-v1 ./comment/server_v1/cmd/server.go
	@echo "$(YELLOW)Building Comment V2 service...$(RESET)"
	@go build -o bin/comment-v2 ./comment/server_v2/cmd/server.go
	@echo "$(GREEN)All services build completed$(RESET)"

clean: ## Clean build files
	@echo "$(BLUE)Cleaning build files...$(RESET)"
	@rm -rf bin/
	@go clean
	@echo "$(GREEN)Clean completed$(RESET)"

# Start all services (V1 version)
start: build ## Start all services (V1 version)
	@echo "$(BLUE)Starting all microservices...$(RESET)"
	@mkdir -p logs pids
	@echo "$(YELLOW)Starting Comment V1 service (port 20004)...$(RESET)"
	@nohup ./bin/comment-v1 > logs/comment-v1.log 2>&1 & echo $$! > pids/comment-v1.pid
	@sleep 2
	@echo "$(YELLOW)Starting Detail V1 service (port 20002)...$(RESET)"
	@nohup ./bin/detail-v1 > logs/detail-v1.log 2>&1 & echo $$! > pids/detail-v1.pid
	@sleep 2
	@echo "$(YELLOW)Starting Order V1 service (port 20003)...$(RESET)"
	@nohup ./bin/order-v1 > logs/order-v1.log 2>&1 & echo $$! > pids/order-v1.pid
	@sleep 2
	@echo "$(YELLOW)Starting User service (port 20001)...$(RESET)"
	@nohup ./bin/user > logs/user.log 2>&1 & echo $$! > pids/user.pid
	@sleep 2
	@echo "$(YELLOW)Starting Frontend service (port 8080)...$(RESET)"
	@nohup ./bin/frontend > logs/frontend.log 2>&1 & echo $$! > pids/frontend.pid
	@sleep 2
	@echo "$(GREEN)All services started$(RESET)"
	@echo "$(GREEN)Frontend access URL: http://localhost:8080$(RESET)"

# Start V2 version services
start-v2: build ## Start all services (V2 version)
	@echo "$(BLUE)Starting all microservices (V2 version)...$(RESET)"
	@mkdir -p logs pids
	@echo "$(YELLOW)Starting Comment V2 service...$(RESET)"
	@nohup ./bin/comment-v2 > logs/comment-v2.log 2>&1 & echo $$! > pids/comment-v2.pid
	@sleep 2
	@echo "$(YELLOW)Starting Detail V2 service...$(RESET)"
	@nohup ./bin/detail-v2 > logs/detail-v2.log 2>&1 & echo $$! > pids/detail-v2.pid
	@sleep 2
	@echo "$(YELLOW)Starting Order V2 service...$(RESET)"
	@nohup ./bin/order-v2 > logs/order-v2.log 2>&1 & echo $$! > pids/order-v2.pid
	@sleep 2
	@echo "$(YELLOW)Starting User service...$(RESET)"
	@nohup ./bin/user > logs/user.log 2>&1 & echo $$! > pids/user.pid
	@sleep 2
	@echo "$(YELLOW)Starting Frontend service...$(RESET)"
	@nohup ./bin/frontend > logs/frontend.log 2>&1 & echo $$! > pids/frontend.pid
	@sleep 2
	@echo "$(GREEN)All services (V2 version) started$(RESET)"
	@echo "$(GREEN)Frontend access URL: http://localhost:8080$(RESET)"

stop: ## Stop all services
	@echo "$(BLUE)Stopping all services...$(RESET)"
	@mkdir -p pids
	@if [ -f pids/frontend.pid ]; then \
		kill `cat pids/frontend.pid` 2>/dev/null || true; \
		rm -f pids/frontend.pid; \
		echo "$(YELLOW)Frontend service stopped$(RESET)"; \
	fi
	@if [ -f pids/user.pid ]; then \
		kill `cat pids/user.pid` 2>/dev/null || true; \
		rm -f pids/user.pid; \
		echo "$(YELLOW)User service stopped$(RESET)"; \
	fi
	@if [ -f pids/detail-v1.pid ]; then \
		kill `cat pids/detail-v1.pid` 2>/dev/null || true; \
		rm -f pids/detail-v1.pid; \
		echo "$(YELLOW)Detail V1 service stopped$(RESET)"; \
	fi
	@if [ -f pids/detail-v2.pid ]; then \
		kill `cat pids/detail-v2.pid` 2>/dev/null || true; \
		rm -f pids/detail-v2.pid; \
		echo "$(YELLOW)Detail V2 service stopped$(RESET)"; \
	fi
	@if [ -f pids/order-v1.pid ]; then \
		kill `cat pids/order-v1.pid` 2>/dev/null || true; \
		rm -f pids/order-v1.pid; \
		echo "$(YELLOW)Order V1 service stopped$(RESET)"; \
	fi
	@if [ -f pids/order-v2.pid ]; then \
		kill `cat pids/order-v2.pid` 2>/dev/null || true; \
		rm -f pids/order-v2.pid; \
		echo "$(YELLOW)Order V2 service stopped$(RESET)"; \
	fi
	@if [ -f pids/comment-v1.pid ]; then \
		kill `cat pids/comment-v1.pid` 2>/dev/null || true; \
		rm -f pids/comment-v1.pid; \
		echo "$(YELLOW)Comment V1 service stopped$(RESET)"; \
	fi
	@if [ -f pids/comment-v2.pid ]; then \
		kill `cat pids/comment-v2.pid` 2>/dev/null || true; \
		rm -f pids/comment-v2.pid; \
		echo "$(YELLOW)Comment V2 service stopped$(RESET)"; \
	fi
	@echo "$(GREEN)All services stopped$(RESET)"

restart: stop start ## Restart all services

restart-v2: stop start-v2 ## Restart all services (V2 version)

logs: ## View all service logs
	@echo "$(BLUE)Viewing service logs...$(RESET)"
	@if [ -d logs ]; then \
		echo "$(YELLOW)=== Frontend Logs ===$(RESET)"; \
		tail -n 20 logs/frontend.log 2>/dev/null || echo "No Frontend logs"; \
		echo "$(YELLOW)=== User Logs ===$(RESET)"; \
		tail -n 20 logs/user.log 2>/dev/null || echo "No User logs"; \
		echo "$(YELLOW)=== Detail Logs ===$(RESET)"; \
		tail -n 20 logs/detail-v1.log 2>/dev/null || tail -n 20 logs/detail-v2.log 2>/dev/null || echo "No Detail logs"; \
		echo "$(YELLOW)=== Order Logs ===$(RESET)"; \
		tail -n 20 logs/order-v1.log 2>/dev/null || tail -n 20 logs/order-v2.log 2>/dev/null || echo "No Order logs"; \
		echo "$(YELLOW)=== Comment Logs ===$(RESET)"; \
		tail -n 20 logs/comment-v1.log 2>/dev/null || tail -n 20 logs/comment-v2.log 2>/dev/null || echo "No Comment logs"; \
	else \
		echo "$(RED)Logs directory does not exist$(RESET)"; \
	fi

status: ## Check service status
	@echo "$(BLUE)Checking service status...$(RESET)"
	@mkdir -p pids
	@printf "$(YELLOW)Frontend:$(RESET) "
	@if [ -f pids/frontend.pid ] && kill -0 $$(cat pids/frontend.pid) 2>/dev/null; then echo "$(GREEN)Running$(RESET)"; else echo "$(RED)Stopped$(RESET)"; fi
	@printf "$(YELLOW)User:$(RESET)     "
	@if [ -f pids/user.pid ] && kill -0 $$(cat pids/user.pid) 2>/dev/null; then echo "$(GREEN)Running$(RESET)"; else echo "$(RED)Stopped$(RESET)"; fi
	@printf "$(YELLOW)Detail:$(RESET)   "
	@if [ -f pids/detail-v1.pid ] && kill -0 $$(cat pids/detail-v1.pid) 2>/dev/null; then echo "$(GREEN)V1 Running$(RESET)"; elif [ -f pids/detail-v2.pid ] && kill -0 $$(cat pids/detail-v2.pid) 2>/dev/null; then echo "$(GREEN)V2 Running$(RESET)"; else echo "$(RED)Stopped$(RESET)"; fi
	@printf "$(YELLOW)Order:$(RESET)    "
	@if [ -f pids/order-v1.pid ] && kill -0 $$(cat pids/order-v1.pid) 2>/dev/null; then echo "$(GREEN)V1 Running$(RESET)"; elif [ -f pids/order-v2.pid ] && kill -0 $$(cat pids/order-v2.pid) 2>/dev/null; then echo "$(GREEN)V2 Running$(RESET)"; else echo "$(RED)Stopped$(RESET)"; fi
	@printf "$(YELLOW)Comment:$(RESET)  "
	@if [ -f pids/comment-v1.pid ] && kill -0 $$(cat pids/comment-v1.pid) 2>/dev/null; then echo "$(GREEN)V1 Running$(RESET)"; elif [ -f pids/comment-v2.pid ] && kill -0 $$(cat pids/comment-v2.pid) 2>/dev/null; then echo "$(GREEN)V2 Running$(RESET)"; else echo "$(RED)Stopped$(RESET)"; fi

test: ## Run tests
	@echo "$(BLUE)Running tests...$(RESET)"
	@go test -v ./...
	@echo "$(GREEN)Tests completed$(RESET)"

proto: ## Generate Protocol Buffers code
	@echo "$(BLUE)Generating Protocol Buffers code...$(RESET)"
	@echo "$(YELLOW)Generating Comment API...$(RESET)"
	@cd comment/api && protoc --go_out=. --go-grpc_out=. --go-triple_out=. *.proto
	@echo "$(YELLOW)Generating Detail API...$(RESET)"
	@cd detail/api && protoc --go_out=. --go-grpc_out=. --go-triple_out=. *.proto
	@echo "$(YELLOW)Generating Order API...$(RESET)"
	@cd order/api && protoc --go_out=. --go-grpc_out=. --go-triple_out=. *.proto
	@echo "$(YELLOW)Generating User API...$(RESET)"
	@cd user/api && protoc --go_out=. --go-grpc_out=. --go-triple_out=. *.proto
	@echo "$(GREEN)Protocol Buffers code generation completed$(RESET)"

# Client test commands
test-user: ## Test user service client
	@echo "$(BLUE)Testing user service...$(RESET)"
	@go run ./user/go-client/cmd/client.go

test-detail: ## Test product detail service client
	@echo "$(BLUE)Testing product detail service...$(RESET)"
	@go run ./detail/client/cmd/client.go

test-order: ## Test order service client
	@echo "$(BLUE)Testing order service...$(RESET)"
	@go run ./order/client/cmd/client.go

test-comment: ## Test comment service client
	@echo "$(BLUE)Testing comment service...$(RESET)"
	@go run ./comment/client/cmd/client.go

# Utility commands
init: ## Initialize project directories
	@echo "$(BLUE)Initializing project directories...$(RESET)"
	@mkdir -p bin logs pids
	@echo "$(GREEN)Project directories initialization completed$(RESET)"

fmt: ## Format code
	@echo "$(BLUE)Formatting Go code...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)Code formatting completed$(RESET)"

vet: ## Static code analysis
	@echo "$(BLUE)Running static code analysis...$(RESET)"
	@go vet ./...
	@echo "$(GREEN)Code analysis completed$(RESET)"

mod-update: ## Update dependency modules
	@echo "$(BLUE)Updating dependency modules...$(RESET)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)Dependency modules update completed$(RESET)"

# Clean commands
clean-all: clean ## Deep clean (including logs and PID files)
	@echo "$(BLUE)Deep cleaning...$(RESET)"
	@rm -rf logs/ pids/
	@echo "$(GREEN)Deep clean completed$(RESET)"

# Quick start
quick-start: init install-deps start ## Quick start (init + install dependencies + start)
	@echo "$(GREEN)Quick start completed! Visit http://localhost:8080$(RESET)"