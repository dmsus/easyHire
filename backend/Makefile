.PHONY: help build run test clean migrate docker-up docker-down lint

help: ## Show this help message
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building EasyHire Backend..."
	@go build -o bin/backend cmd/backend/main.go

run: ## Run the application locally
	@echo "Starting EasyHire Backend..."
	@go run cmd/backend/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test ./... -v

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/ coverage.out coverage.html

migrate: ## Run database migrations
	@echo "Running migrations..."
	@go run cmd/migrate/main.go

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	@docker-compose -f ../docker/docker-compose.dev.yml up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker-compose -f ../docker/docker-compose.dev.yml down

lint: ## Run linter
	@echo "Running golangci-lint..."
	@golangci-lint run

format: ## Format code
	@echo "Formatting code..."
	@gofmt -w -s .

tidy: ## Tidy go.mod
	@echo "Tidying dependencies..."
	@go mod tidy

check: ## Run all checks
	@make lint
	@make test
	@make build

dev: ## Start development environment
	@make docker-up
	@sleep 5
	@make run

version: ## Show version
	@echo "Version: $(shell git describe --tags --always --dirty)"
	@echo "Build: $(build)"

