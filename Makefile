.PHONY: help build run test clean migrate-up migrate-down migrate-status migrate-create

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building..."
	@go build -o bin/api cmd/api/main.go

run: ## Run the application
	@echo "Running..."
	@go run cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

dev: ## Run with hot reload (requires air)
	@air

docker-build: ## Build Docker image
	@docker build -t plantpal-backend:latest .

docker-run: ## Run Docker container
	@docker run -p 8080:8080 --env-file .env plantpal-backend:latest

docker-up: ## Start all services with docker-compose
	@docker compose up -d

docker-down: ## Stop all services
	@docker compose down

docker-logs: ## View logs from all services
	@docker compose logs -f

docker-rebuild: ## Rebuild and restart services
	@docker compose up -d --build

migrate-up: ## Run all pending migrations
	@goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down: ## Rollback the last migration
	@goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-status: ## Show migration status
	@goose -dir migrations postgres "$(DATABASE_URL)" status

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@goose -dir migrations create $(NAME) sql
