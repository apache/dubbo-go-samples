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

# Define service list and corresponding executable names
SERVICES := adservice cartservice checkoutservice currencyservice emailservice frontendservice paymentservice productcatalogservice recommendationservice shippingservice

# Define mapping from service directory names to executable names
# Define color output
BLUE := \033[34m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m # No Color

# Default target
.PHONY: help
help:
	@echo "$(BLUE)Online Boutique Makefile$(NC)"
	@echo "$(GREEN)Available targets:$(NC)"
	@echo "  $(YELLOW)build-all$(NC)      - Build all services"
	@echo "  $(YELLOW)start-all$(NC)      - Start all services (run in background)"
	@echo "  $(YELLOW)stop-all$(NC)       - Stop all services"
	@echo "  $(YELLOW)clean-all$(NC)      - Clean all build files"
	@echo "  $(YELLOW)tidy-all$(NC)       - Execute go mod tidy for all services"
	@echo "  $(YELLOW)logs$(NC)           - View all service logs"
	@echo "  $(YELLOW)status$(NC)         - Check all service status"
	@echo ""
	@echo "$(GREEN)Individual service operations:$(NC)"
	@echo "  $(YELLOW)build-<service>$(NC)  - Build specified service"
	@echo "  $(YELLOW)start-<service>$(NC)  - Start specified service"
	@echo "  $(YELLOW)stop-<service>$(NC)   - Stop specified service"
	@echo ""
	@echo "$(GREEN)Available services:$(NC) $(SERVICES)"

# Build all services
.PHONY: build-all
build-all:
	@echo "$(BLUE)Building all services...$(NC)"
	@for service in $(SERVICES); do \
		echo "$(GREEN)Building $$service...$(NC)"; \
		$(MAKE) -C src/$$service build; \
	done
	@echo "$(GREEN)All services built successfully!$(NC)"

# Start all services
.PHONY: start-all
start-all: build-all setup
	@echo "$(BLUE)Starting all services...$(NC)"
	@for service in $(SERVICES); do \
		echo "$(GREEN)Starting $$service...$(NC)"; \
		(cd src/$$service; \
		 nohup ./$$service > ../../logs/$$service.log 2>&1 & echo $$! > ../../pids/$$service.pid); \
		sleep 2; \
	done
	@echo "$(GREEN)All services started!$(NC)"

# Stop all services
.PHONY: stop-all
stop-all:
	@echo "$(BLUE)Stopping all services...$(NC)"
	@for service in $(SERVICES); do \
		if [ -f pids/$$service.pid ]; then \
			echo "$(YELLOW)Stopping $$service...$(NC)"; \
			kill -TERM $$(cat pids/$$service.pid) 2>/dev/null || true; \
			rm -f pids/$$service.pid; \
		fi; \
	done
	@echo "$(GREEN)All services stopped!$(NC)"

# Clean all build files
.PHONY: clean-all
clean-all:
	@echo "$(BLUE)Cleaning all build files...$(NC)"
	@for service in $(SERVICES); do \
		echo "$(YELLOW)Cleaning $$service...$(NC)"; \
		rm -f src/$$service/$$service; \
	done
	@rm -rf logs pids
	@echo "$(GREEN)Cleanup complete!$(NC)"

# Execute go mod tidy for all services
.PHONY: tidy-all
tidy-all:
	@echo "$(BLUE)Executing go mod tidy for all services...$(NC)"
	@for service in $(SERVICES); do \
		echo "$(GREEN)Tidying $$service...$(NC)"; \
		$(MAKE) -C src/$$service tidy; \
	done
	@echo "$(GREEN)All services tidy complete!$(NC)"

# View all service logs
.PHONY: logs
logs:
	@echo "$(BLUE)Viewing all service logs...$(NC)"
	@mkdir -p logs
	@for service in $(SERVICES); do \
		if [ -f logs/$$service.log ]; then \
			echo "$(GREEN)=== $$service logs ===$(NC)"; \
			tail -20 logs/$$service.log; \
			echo ""; \
		fi; \
	done

# Check all service status
.PHONY: status
status:
	@echo "$(BLUE)Checking all service status...$(NC)"
	@mkdir -p pids
	@for service in $(SERVICES); do \
		if [ -f pids/$$service.pid ]; then \
			pid=$$(cat pids/$$service.pid); \
			if ps -p $$pid > /dev/null 2>&1; then \
				echo "$(GREEN)✓ $$service (PID: $$pid) - Running$(NC)"; \
			else \
				echo "$(RED)✗ $$service - Stopped$(NC)"; \
				rm -f pids/$$service.pid; \
			fi; \
		else \
			echo "$(YELLOW)- $$service - Not started$(NC)"; \
		fi; \
	done

# Create necessary directories
.PHONY: setup
setup:
	@mkdir -p logs pids

# Template for individual service operations
define service-template
.PHONY: build-$(1)
build-$(1):
	@echo "$(GREEN)Building $(1)...$(NC)"
	@$(MAKE) -C src/$(1) build

.PHONY: start-$(1)
start-$(1): build-$(1) setup
	@echo "$(GREEN)Starting $(1)...$(NC)"
	@cd src/$(1) && nohup ./$(1) > ../../logs/$(1).log 2>&1 & echo $$! > ../../pids/$(1).pid

.PHONY: stop-$(1)
stop-$(1):
	@mkdir -p pids
	@if [ -f pids/$(1).pid ]; then \
		echo "$(YELLOW)Stopping $(1)...$(NC)"; \
		kill -TERM $$(cat pids/$(1).pid) 2>/dev/null || true; \
		rm -f pids/$(1).pid; \
	else \
		echo "$(YELLOW)$(1) is not running$(NC)"; \
	fi
endef

# Generate template for each service
$(foreach service,$(SERVICES),$(eval $(call service-template,$(service))))

# Environment check
.PHONY: check-env
check-env:
	@echo "$(BLUE)Checking environment dependencies...$(NC)"
	@which go > /dev/null || (echo "$(RED)Error: Go not found$(NC)" && exit 1)
	@echo "$(GREEN)Environment check passed!$(NC)"

# Initialize project
.PHONY: init
init: check-env setup tidy-all
	@echo "$(GREEN)Project initialization complete!$(NC)"
	@echo "$(BLUE)Use 'make start-all' to start all services$(NC)"
	@echo "$(BLUE)Use 'make help' to view more commands$(NC)"