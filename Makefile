.PHONY: all build generate lint test update update-templates help

# Версия приложения
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -X 'main.version=$(VERSION)'

all: lint test generate build

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

update-templates: ## Update templates
	@UPDATE_GOLDEN=1 go test ./templates_tests

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
