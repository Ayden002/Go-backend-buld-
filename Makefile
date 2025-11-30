.PHONY: help docker-up docker-down test-e2e clean build

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

docker-up: ## Build and start all containers, wait for health checks
	@echo "Building and starting containers..."
	docker-compose up -d --build
	@echo "Waiting for services to be healthy..."
	@timeout=60; \
	while [ $$timeout -gt 0 ]; do \
		if docker-compose ps | grep -q "healthy"; then \
			echo "Services are healthy!"; \
			exit 0; \
		fi; \
		echo "Waiting for health checks... ($$timeout seconds remaining)"; \
		sleep 2; \
		timeout=$$((timeout - 2)); \
	done; \
	echo "Timeout waiting for services to be healthy"; \
	docker-compose logs; \
	exit 1

docker-down: ## Stop and remove all containers
	@echo "Stopping and removing containers..."
	docker-compose down -v
	@echo "Cleanup complete"

test-e2e: ## Run end-to-end tests
	@echo "Running E2E tests..."
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Please copy .env.example to .env and configure it."; \
		exit 1; \
	fi
	bash ./e2e-test.sh

build: ## Build the Go application locally
	@echo "Building application..."
	go build -o server ./cmd/server
	@echo "Build complete: ./server"

clean: ## Clean up build artifacts and containers
	@echo "Cleaning up..."
	docker-compose down -v
	rm -f server server.exe
	@echo "Clean complete"

.DEFAULT_GOAL := help

