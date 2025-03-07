.PHONY: all build generate lint test update help

# Версия приложения
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

all: lint test build

build: ## Build the project
	@echo "==> Building version $(VERSION)..."
	@go build -ldflags "$(LDFLAGS)" -o bin/envgen cmd/envgen/main.go

generate: ## Generate the project
	@go generate ./...

lint: ## Run linters
	@echo "==> Running linters..."
	@golangci-lint run

test: ## Run tests
	@echo "==> Running tests..."
	@go test -v ./...

update: ## Update dependencies
	@echo "==> Updating dependencies..."
	@go get -u ./...
	@go mod tidy

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
