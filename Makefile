.PHONY: run build test clean docker-build setup-env help

# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run:
	@echo "Running the application..."
	@cd cmd && go run main.go

build:
	@echo "Building the application..."
	@go build -o bin/fin-ai ./cmd

test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out

docker-build:
	@echo "Building Docker image..."
	@docker build -t fin-ai:latest .

setup-env:
	@echo "Setting up environment file..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file from .env.example. Please edit it with your values."; else echo ".env file already exists."; fi

# GitHub Actions local testing (requires act: https://github.com/nektos/act)
test-ci:
	@echo "Testing CI workflow locally..."
	@act -j test

lint:
	@echo "Running linter..."
	@golangci-lint run

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

help:
	@echo "Available commands:"
	@echo "  run         - Run the application"
	@echo "  build       - Build the application"
	@echo "  test        - Run tests with coverage"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-build- Build Docker image"
	@echo "  setup-env   - Create .env file from template"
	@echo "  test-ci     - Test CI workflow locally (requires act)"
	@echo "  lint        - Run linter"
	@echo "  deps        - Install and tidy dependencies"
	@echo "  help        - Show this help message"