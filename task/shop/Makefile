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

# Dubbo-Go 微服务商城项目 Makefile
# 提供一键启动、停止、构建等功能

.PHONY: help build clean start stop restart logs test proto check-deps install-deps

# 默认目标
.DEFAULT_GOAL := help

# 项目信息
PROJECT_NAME := dubbo-go-shop
GO_MIN_VERSION := 1.19
REGISTRY := nacos

# 服务配置
SERVICES := frontend user detail order comment
SERVICE_PORTS := 8080 20001 20002 20003 20004

# 颜色定义
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

help: ## 显示帮助信息
	@echo "$(GREEN)Dubbo-Go 微服务商城项目 Makefile$(RESET)"
	@echo ""
	@echo "$(YELLOW)可用命令:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(BLUE)%-15s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

check-deps: ## 检查依赖项
	@echo "$(BLUE)检查Go版本...$(RESET)"
	@go version || (echo "$(RED)未找到Go环境$(RESET)" && exit 1)
	@echo "$(GREEN)Go环境检查通过$(RESET)"

install-deps: check-deps ## 安装项目依赖
	@echo "$(BLUE)安装项目依赖...$(RESET)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)依赖安装完成$(RESET)"

build: install-deps ## 构建所有服务
	@echo "$(BLUE)构建所有服务...$(RESET)"
	@mkdir -p bin
	@echo "$(YELLOW)构建 Frontend 服务...$(RESET)"
	@go build -o bin/frontend ./frontend/cmd/main.go
	@echo "$(YELLOW)构建 User 服务...$(RESET)"
	@go build -o bin/user ./user/go-server/cmd/server.go
	@echo "$(YELLOW)构建 Detail V1 服务...$(RESET)"
	@go build -o bin/detail-v1 ./detail/server_v1/cmd/server.go
	@echo "$(YELLOW)构建 Detail V2 服务...$(RESET)"
	@go build -o bin/detail-v2 ./detail/server_v2/cmd/server.go
	@echo "$(YELLOW)构建 Order V1 服务...$(RESET)"
	@go build -o bin/order-v1 ./order/server_v1/cmd/server.go
	@echo "$(YELLOW)构建 Order V2 服务...$(RESET)"
	@go build -o bin/order-v2 ./order/server_v2/cmd/server.go
	@echo "$(YELLOW)构建 Comment V1 服务...$(RESET)"
	@go build -o bin/comment-v1 ./comment/server_v1/cmd/server.go
	@echo "$(YELLOW)构建 Comment V2 服务...$(RESET)"
	@go build -o bin/comment-v2 ./comment/server_v2/cmd/server.go
	@echo "$(GREEN)所有服务构建完成$(RESET)"

clean: ## 清理构建文件
	@echo "$(BLUE)清理构建文件...$(RESET)"
	@rm -rf bin/
	@go clean
	@echo "$(GREEN)清理完成$(RESET)"

# 启动所有服务 (V1版本)
start: build ## 启动所有服务 (V1版本)
	@echo "$(BLUE)启动所有微服务...$(RESET)"
	@mkdir -p logs
	@echo "$(YELLOW)启动 Comment V1 服务 (端口 20004)...$(RESET)"
	@nohup ./bin/comment-v1 > logs/comment-v1.log 2>&1 & echo $$! > pids/comment-v1.pid
	@sleep 2
	@echo "$(YELLOW)启动 Detail V1 服务 (端口 20002)...$(RESET)"
	@nohup ./bin/detail-v1 > logs/detail-v1.log 2>&1 & echo $$! > pids/detail-v1.pid
	@sleep 2
	@echo "$(YELLOW)启动 Order V1 服务 (端口 20003)...$(RESET)"
	@nohup ./bin/order-v1 > logs/order-v1.log 2>&1 & echo $$! > pids/order-v1.pid
	@sleep 2
	@echo "$(YELLOW)启动 User 服务 (端口 20001)...$(RESET)"
	@nohup ./bin/user > logs/user.log 2>&1 & echo $$! > pids/user.pid
	@sleep 2
	@echo "$(YELLOW)启动 Frontend 服务 (端口 8080)...$(RESET)"
	@nohup ./bin/frontend > logs/frontend.log 2>&1 & echo $$! > pids/frontend.pid
	@sleep 2
	@echo "$(GREEN)所有服务启动完成$(RESET)"
	@echo "$(GREEN)前端访问地址: http://localhost:8080$(RESET)"

# 启动 V2 版本服务
start-v2: build ## 启动所有服务 (V2版本)
	@echo "$(BLUE)启动所有微服务 (V2版本)...$(RESET)"
	@mkdir -p logs pids
	@echo "$(YELLOW)启动 Comment V2 服务...$(RESET)"
	@nohup ./bin/comment-v2 > logs/comment-v2.log 2>&1 & echo $$! > pids/comment-v2.pid
	@sleep 2
	@echo "$(YELLOW)启动 Detail V2 服务...$(RESET)"
	@nohup ./bin/detail-v2 > logs/detail-v2.log 2>&1 & echo $$! > pids/detail-v2.pid
	@sleep 2
	@echo "$(YELLOW)启动 Order V2 服务...$(RESET)"
	@nohup ./bin/order-v2 > logs/order-v2.log 2>&1 & echo $$! > pids/order-v2.pid
	@sleep 2
	@echo "$(YELLOW)启动 User 服务...$(RESET)"
	@nohup ./bin/user > logs/user.log 2>&1 & echo $$! > pids/user.pid
	@sleep 2
	@echo "$(YELLOW)启动 Frontend 服务...$(RESET)"
	@nohup ./bin/frontend > logs/frontend.log 2>&1 & echo $$! > pids/frontend.pid
	@sleep 2
	@echo "$(GREEN)所有服务 (V2版本) 启动完成$(RESET)"
	@echo "$(GREEN)前端访问地址: http://localhost:8080$(RESET)"

stop: ## 停止所有服务
	@echo "$(BLUE)停止所有服务...$(RESET)"
	@mkdir -p pids
	@if [ -f pids/frontend.pid ]; then \
		kill `cat pids/frontend.pid` 2>/dev/null || true; \
		rm -f pids/frontend.pid; \
		echo "$(YELLOW)Frontend 服务已停止$(RESET)"; \
	fi
	@if [ -f pids/user.pid ]; then \
		kill `cat pids/user.pid` 2>/dev/null || true; \
		rm -f pids/user.pid; \
		echo "$(YELLOW)User 服务已停止$(RESET)"; \
	fi
	@if [ -f pids/detail-v1.pid ]; then \
		kill `cat pids/detail-v1.pid` 2>/dev/null || true; \
		rm -f pids/detail-v1.pid; \
		echo "$(YELLOW)Detail V1 服务已停止$(RESET)"; \
	fi
	@if [ -f pids/detail-v2.pid ]; then \
		kill `cat pids/detail-v2.pid` 2>/dev/null || true; \
		rm -f pids/detail-v2.pid; \
		echo "$(YELLOW)Detail V2 服务已停止$(RESET)"; \
	fi
	@if [ -f pids/order-v1.pid ]; then \
		kill `cat pids/order-v1.pid` 2>/dev/null || true; \
		rm -f pids/order-v1.pid; \
		echo "$(YELLOW)Order V1 服务已停止$(RESET)"; \
	fi
	@if [ -f pids/order-v2.pid ]; then \
		kill `cat pids/order-v2.pid` 2>/dev/null || true; \
		rm -f pids/order-v2.pid; \
		echo "$(YELLOW)Order V2 服务已停止$(RESET)"; \
	fi
	@if [ -f pids/comment-v1.pid ]; then \
		kill `cat pids/comment-v1.pid` 2>/dev/null || true; \
		rm -f pids/comment-v1.pid; \
		echo "$(YELLOW)Comment V1 服务已停止$(RESET)"; \
	fi
	@if [ -f pids/comment-v2.pid ]; then \
		kill `cat pids/comment-v2.pid` 2>/dev/null || true; \
		rm -f pids/comment-v2.pid; \
		echo "$(YELLOW)Comment V2 服务已停止$(RESET)"; \
	fi
	@echo "$(GREEN)所有服务已停止$(RESET)"

restart: stop start ## 重启所有服务

restart-v2: stop start-v2 ## 重启所有服务 (V2版本)

logs: ## 查看所有服务日志
	@echo "$(BLUE)查看服务日志...$(RESET)"
	@if [ -d logs ]; then \
		echo "$(YELLOW)=== Frontend 日志 ===$(RESET)"; \
		tail -n 20 logs/frontend.log 2>/dev/null || echo "无 Frontend 日志"; \
		echo "$(YELLOW)=== User 日志 ===$(RESET)"; \
		tail -n 20 logs/user.log 2>/dev/null || echo "无 User 日志"; \
		echo "$(YELLOW)=== Detail 日志 ===$(RESET)"; \
		tail -n 20 logs/detail-v1.log 2>/dev/null || tail -n 20 logs/detail-v2.log 2>/dev/null || echo "无 Detail 日志"; \
		echo "$(YELLOW)=== Order 日志 ===$(RESET)"; \
		tail -n 20 logs/order-v1.log 2>/dev/null || tail -n 20 logs/order-v2.log 2>/dev/null || echo "无 Order 日志"; \
		echo "$(YELLOW)=== Comment 日志 ===$(RESET)"; \
		tail -n 20 logs/comment-v1.log 2>/dev/null || tail -n 20 logs/comment-v2.log 2>/dev/null || echo "无 Comment 日志"; \
	else \
		echo "$(RED)日志目录不存在$(RESET)"; \
	fi

status: ## 查看服务状态
	@echo "$(BLUE)查看服务状态...$(RESET)"
	@mkdir -p pids
	@echo "$(YELLOW)Frontend:$(RESET) $(shell if [ -f pids/frontend.pid ] && kill -0 `cat pids/frontend.pid` 2>/dev/null; then echo "$(GREEN)运行中$(RESET)"; else echo "$(RED)已停止$(RESET)"; fi)"
	@echo "$(YELLOW)User:$(RESET)     $(shell if [ -f pids/user.pid ] && kill -0 `cat pids/user.pid` 2>/dev/null; then echo "$(GREEN)运行中$(RESET)"; else echo "$(RED)已停止$(RESET)"; fi)"
	@echo "$(YELLOW)Detail:$(RESET)   $(shell if [ -f pids/detail-v1.pid ] && kill -0 `cat pids/detail-v1.pid` 2>/dev/null; then echo "$(GREEN)V1运行中$(RESET)"; elif [ -f pids/detail-v2.pid ] && kill -0 `cat pids/detail-v2.pid` 2>/dev/null; then echo "$(GREEN)V2运行中$(RESET)"; else echo "$(RED)已停止$(RESET)"; fi)"
	@echo "$(YELLOW)Order:$(RESET)    $(shell if [ -f pids/order-v1.pid ] && kill -0 `cat pids/order-v1.pid` 2>/dev/null; then echo "$(GREEN)V1运行中$(RESET)"; elif [ -f pids/order-v2.pid ] && kill -0 `cat pids/order-v2.pid` 2>/dev/null; then echo "$(GREEN)V2运行中$(RESET)"; else echo "$(RED)已停止$(RESET)"; fi)"
	@echo "$(YELLOW)Comment:$(RESET)  $(shell if [ -f pids/comment-v1.pid ] && kill -0 `cat pids/comment-v1.pid` 2>/dev/null; then echo "$(GREEN)V1运行中$(RESET)"; elif [ -f pids/comment-v2.pid ] && kill -0 `cat pids/comment-v2.pid` 2>/dev/null; then echo "$(GREEN)V2运行中$(RESET)"; else echo "$(RED)已停止$(RESET)"; fi)"

test: ## 运行测试
	@echo "$(BLUE)运行测试...$(RESET)"
	@go test -v ./...
	@echo "$(GREEN)测试完成$(RESET)"

proto: ## 生成 Protocol Buffers 代码
	@echo "$(BLUE)生成 Protocol Buffers 代码...$(RESET)"
	@echo "$(YELLOW)生成 Comment API...$(RESET)"
	@cd comment/api && protoc --go_out=. --go-grpc_out=. --go-triple_out=. *.proto
	@echo "$(YELLOW)生成 Detail API...$(RESET)"
	@cd detail/api && protoc --go_out=. --go-grpc_out=. --go-triple_out=. *.proto
	@echo "$(YELLOW)生成 Order API...$(RESET)"
	@cd order/api && protoc --go_out=. --go-grpc_out=. --go-triple_out=. *.proto
	@echo "$(YELLOW)生成 User API...$(RESET)"
	@cd user/api && protoc --go_out=. --go-grpc_out=. --go-triple_out=. *.proto
	@echo "$(GREEN)Protocol Buffers 代码生成完成$(RESET)"

# 客户端测试命令
test-user: ## 测试用户服务客户端
	@echo "$(BLUE)测试用户服务...$(RESET)"
	@go run ./user/go-client/cmd/client.go

test-detail: ## 测试商品详情服务客户端
	@echo "$(BLUE)测试商品详情服务...$(RESET)"
	@go run ./detail/client/cmd/client.go

test-order: ## 测试订单服务客户端
	@echo "$(BLUE)测试订单服务...$(RESET)"
	@go run ./order/client/cmd/client.go

test-comment: ## 测试评论服务客户端
	@echo "$(BLUE)测试评论服务...$(RESET)"
	@go run ./comment/client/cmd/client.go

# 清理命令
clean-all: clean ## 深度清理 (包括日志和PID文件)
	@echo "$(BLUE)深度清理...$(RESET)"
	@rm -rf logs/ pids/
	@echo "$(GREEN)深度清理完成$(RESET)"

# 快速开始
quick-start: init install-deps start ## 快速开始 (初始化+安装依赖+启动)
	@echo "$(GREEN)快速启动完成! 访问 http://localhost:8080$(RESET)"
