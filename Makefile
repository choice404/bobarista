.PHONY: build test lint examples clean install-tools build-examples help

# Default help target
help: ## Show this help message
	@echo "🧋 Bobarista Development Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build example binaries
build: ## Build all example binaries
	@echo "🔨 Building example binaries..."
	mkdir -p dist
	go build -o dist/basic ./examples/basic
	go build -o dist/advanced ./examples/advanced
	go build -o dist/wizard ./examples/wizard
	go build -o dist/registration ./examples/registration
	go build -o dist/survey ./examples/survey
	go build -o dist/themes ./examples/themes
	@echo "✅ Example binaries built in dist/"

test: ## Run all tests
	@echo "🧪 Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "📊 Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

lint: ## Run linter
	@echo "🔍 Running linter..."
	golangci-lint run

examples: ## Run example applications interactively
	@echo "🎯 Running basic example..."
	go run examples/basic/main.go
	@echo "🚀 Running advanced example..."
	go run examples/advanced/main.go

clean: ## Clean build artifacts and cache
	@echo "🧹 Cleaning up..."
	go clean ./...
	rm -rf dist/ coverage.out coverage.html
	@echo "✅ Cleaned up"

install-tools: ## Install development tools
	@echo "🛠️  Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/godoc@latest
	@echo "✅ Development tools installed"

deps: ## Download and verify dependencies
	@echo "📦 Managing dependencies..."
	go mod download
	go mod verify
	@echo "✅ Dependencies ready"

tidy: ## Tidy and verify go modules
	@echo "🧽 Tidying modules..."
	go mod tidy
	go mod verify
	@echo "✅ Modules tidied"

update: ## Update all dependencies
	@echo "⬆️  Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "✅ Dependencies updated"

release-check: test lint ## Run pre-release checks
	@echo "🚀 Running pre-release checks..."
	@echo "✅ All checks passed - ready for release!"

.DEFAULT_GOAL := help
