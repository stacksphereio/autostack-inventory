.PHONY: help build run test clean docker-build docker-run

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building api-inventory..."
	go build -o bin/server ./cmd/server

run: ## Run the application locally
	@echo "Running api-inventory..."
	DATA_PATH=../../data/seed go run ./cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	go test ./... -v

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t autostack-api-inventory:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8001:8001 \
		-e DATA_PATH=/app/data/seed \
		-v $(PWD)/../../data/seed:/app/data/seed:ro \
		autostack-api-inventory:latest

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
