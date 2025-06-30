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

# 定义服务列表和对应的可执行文件名
SERVICES := adservice cartservice checkoutservice currencyservice emailservice frontendservice paymentservice productcatalogservice recommendationservice shippingservice

# 定义服务目录名到可执行文件名的映射
# 定义颜色输出
BLUE := \033[34m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m # No Color

# 默认目标
.PHONY: help
help:
	@echo "$(BLUE)Online Boutique Makefile$(NC)"
	@echo "$(GREEN)Available targets:$(NC)"
	@echo "  $(YELLOW)build-all$(NC)      - 构建所有服务"
	@echo "  $(YELLOW)start-all$(NC)      - 启动所有服务（后台运行）"
	@echo "  $(YELLOW)stop-all$(NC)       - 停止所有服务"
	@echo "  $(YELLOW)clean-all$(NC)      - 清理所有构建文件"
	@echo "  $(YELLOW)tidy-all$(NC)       - 执行所有服务的 go mod tidy"
	@echo "  $(YELLOW)logs$(NC)           - 查看所有服务日志"
	@echo "  $(YELLOW)status$(NC)         - 检查所有服务状态"
	@echo ""
	@echo "$(GREEN)单个服务操作:$(NC)"
	@echo "  $(YELLOW)build-<service>$(NC)  - 构建指定服务"
	@echo "  $(YELLOW)start-<service>$(NC)  - 启动指定服务"
	@echo "  $(YELLOW)stop-<service>$(NC)   - 停止指定服务"
	@echo ""
	@echo "$(GREEN)可用服务:$(NC) $(SERVICES)"

# 构建所有服务
.PHONY: build-all
build-all:
	@echo "$(BLUE)构建所有服务...$(NC)"
	@for service in $(SERVICES); do \
		echo "$(GREEN)构建 $$service...$(NC)"; \
		$(MAKE) -C src/$$service build; \
	done
	@echo "$(GREEN)所有服务构建完成！$(NC)"

# 启动所有服务
.PHONY: start-all
start-all: build-all setup
	@echo "$(BLUE)启动所有服务...$(NC)"
	@for service in $(SERVICES); do \
		echo "$(GREEN)启动 $$service...$(NC)"; \
		(cd src/$$service; \
		 nohup ./$$service > ../../logs/$$service.log 2>&1 & echo $$! > ../../pids/$$service.pid); \
		sleep 2; \
	done
	@echo "$(GREEN)所有服务已启动！$(NC)"

# 停止所有服务
.PHONY: stop-all
stop-all:
	@echo "$(BLUE)停止所有服务...$(NC)"
	@for service in $(SERVICES); do \
		if [ -f pids/$$service.pid ]; then \
			echo "$(YELLOW)停止 $$service...$(NC)"; \
			kill -TERM $$(cat pids/$$service.pid) 2>/dev/null || true; \
			rm -f pids/$$service.pid; \
		fi; \
	done
	@echo "$(GREEN)所有服务已停止！$(NC)"

# 清理所有构建文件
.PHONY: clean-all
clean-all:
	@echo "$(BLUE)清理所有构建文件...$(NC)"
	@for service in $(SERVICES); do \
		echo "$(YELLOW)清理 $$service...$(NC)"; \
		rm -f src/$$service/$$service; \
	done
	@rm -rf logs pids
	@echo "$(GREEN)清理完成！$(NC)"

# 执行所有服务的 go mod tidy
.PHONY: tidy-all
tidy-all:
	@echo "$(BLUE)执行所有服务的 go mod tidy...$(NC)"
	@for service in $(SERVICES); do \
		echo "$(GREEN)tidy $$service...$(NC)"; \
		$(MAKE) -C src/$$service tidy; \
	done
	@echo "$(GREEN)所有服务 tidy 完成！$(NC)"

# 查看所有服务日志
.PHONY: logs
logs:
	@echo "$(BLUE)查看所有服务日志...$(NC)"
	@mkdir -p logs
	@for service in $(SERVICES); do \
		if [ -f logs/$$service.log ]; then \
			echo "$(GREEN)=== $$service 日志 ===$(NC)"; \
			tail -20 logs/$$service.log; \
			echo ""; \
		fi; \
	done

# 检查所有服务状态
.PHONY: status
status:
	@echo "$(BLUE)检查所有服务状态...$(NC)"
	@mkdir -p pids
	@for service in $(SERVICES); do \
		if [ -f pids/$$service.pid ]; then \
			pid=$$(cat pids/$$service.pid); \
			if ps -p $$pid > /dev/null 2>&1; then \
				echo "$(GREEN)✓ $$service (PID: $$pid) - 运行中$(NC)"; \
			else \
				echo "$(RED)✗ $$service - 已停止$(NC)"; \
				rm -f pids/$$service.pid; \
			fi; \
		else \
			echo "$(YELLOW)- $$service - 未启动$(NC)"; \
		fi; \
	done

# 创建必要的目录
.PHONY: setup
setup:
	@mkdir -p logs pids

# 单个服务操作的模板
define service-template
.PHONY: build-$(1)
build-$(1):
	@echo "$(GREEN)构建 $(1)...$(NC)"
	@$(MAKE) -C src/$(1) build

.PHONY: start-$(1)
start-$(1): build-$(1) setup
	@echo "$(GREEN)启动 $(1)...$(NC)"
	@cd src/$(1) && nohup ./$(1) > ../../logs/$(1).log 2>&1 & echo $$! > ../../pids/$(1).pid

.PHONY: stop-$(1)
stop-$(1):
	@mkdir -p pids
	@if [ -f pids/$(1).pid ]; then \
		echo "$(YELLOW)停止 $(1)...$(NC)"; \
		kill -TERM $$(cat pids/$(1).pid) 2>/dev/null || true; \
		rm -f pids/$(1).pid; \
	else \
		echo "$(YELLOW)$(1) 未在运行$(NC)"; \
	fi
endef

# 为每个服务生成模板
$(foreach service,$(SERVICES),$(eval $(call service-template,$(service))))

# 环境检查
.PHONY: check-env
check-env:
	@echo "$(BLUE)检查环境依赖...$(NC)"
	@which go > /dev/null || (echo "$(RED)错误: 未找到 Go$(NC)" && exit 1)
	@echo "$(GREEN)环境检查通过！$(NC)"

# 初始化项目
.PHONY: init
init: check-env setup tidy-all
	@echo "$(GREEN)项目初始化完成！$(NC)"
	@echo "$(BLUE)使用 'make start-all' 启动所有服务$(NC)"
	@echo "$(BLUE)使用 'make help' 查看更多命令$(NC)"
