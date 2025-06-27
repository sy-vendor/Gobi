# Gobi BI Platform Makefile
.PHONY: help build test clean run dev docker-build docker-run lint format coverage

# Variables
BINARY_NAME=gobi
BUILD_DIR=build
DOCKER_IMAGE=gobi
DOCKER_TAG=latest

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
dev: ## Run in development mode
	@echo "Starting Gobi in development mode..."
	GO_ENV=dev go run cmd/server/main.go

run: ## Run the application
	@echo "Starting Gobi..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Building
build: ## Build the application
	@echo "Building Gobi..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

build-linux: ## Build for Linux
	@echo "Building Gobi for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux ./cmd/server

build-windows: ## Build for Windows
	@echo "Building Gobi for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe ./cmd/server

build-darwin: ## Build for macOS
	@echo "Building Gobi for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin ./cmd/server

# Testing
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	go test -race -v ./...

# Code quality
lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

format: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Database
migrate: ## Run database migrations
	@echo "Running database migrations..."
	@if [ -f "migrations/001_init.up.sql" ]; then \
		echo "SQLite migrations found, running..."; \
		sqlite3 gobi.db < migrations/001_init.up.sql; \
	else \
		echo "No SQLite migrations found"; \
	fi

migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	@if [ -f "migrations/001_init.down.sql" ]; then \
		echo "SQLite migrations found, rolling back..."; \
		sqlite3 gobi.db < migrations/001_init.down.sql; \
	else \
		echo "No SQLite migrations found"; \
	fi

# Docker
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-stop: ## Stop Docker container
	@echo "Stopping Docker container..."
	docker stop $$(docker ps -q --filter ancestor=$(DOCKER_IMAGE):$(DOCKER_TAG)) 2>/dev/null || true

# Data setup
setup-data: ## Setup sample data
	@echo "Setting up sample data..."
	@for script in scripts/data/generate_*.sql; do \
		if [ -f "$$script" ]; then \
			echo "Running $$script..."; \
			sqlite3 gobi.db < "$$script"; \
		fi; \
	done

test-data: ## Test all chart data
	@echo "Testing all chart data..."
	@for script in scripts/testing/test_*.sh; do \
		if [ -f "$$script" ]; then \
			echo "Running $$script..."; \
			bash "$$script"; \
		fi; \
	done

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f gobi.db

clean-docker: ## Clean Docker images
	@echo "Cleaning Docker images..."
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) 2>/dev/null || true

# Dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Security
security-scan: ## Run security scan
	@echo "Running security scan..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "govulncheck not found, installing..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
		govulncheck ./...; \
	fi

# Development setup
setup: ## Setup development environment
	@echo "Setting up development environment..."
	$(MAKE) deps
	$(MAKE) migrate
	$(MAKE) setup-data
	@echo "Development environment setup complete!"

# Production build
prod-build: ## Production build
	@echo "Building for production..."
	$(MAKE) lint
	$(MAKE) test-coverage
	$(MAKE) security-scan
	$(MAKE) build
	@echo "Production build complete!"

# All-in-one development
dev-full: ## Full development workflow
	@echo "Running full development workflow..."
	$(MAKE) deps
	$(MAKE) format
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) build
	@echo "Development workflow complete!" 