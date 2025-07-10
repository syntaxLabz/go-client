# Makefile for httpclient

.PHONY: test build examples clean fmt lint deps check help install-tools

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build parameters
BINARY_DIR=bin
EXAMPLES_DIR=examples

# Test parameters
TEST_TIMEOUT=30s
COVERAGE_FILE=coverage.out

# Default target
all: check

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) -u golang.org/x/tools/cmd/goimports
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) -u github.com/securecodewarrior/sast-scan/cmd/gosec@latest

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...
	goimports -w .

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Run security scan
security:
	@echo "Running security scan..."
	gosec ./...

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) -race ./...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Build examples
build-examples:
	@echo "Building examples..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/basic ./$(EXAMPLES_DIR)/basic/
	$(GOBUILD) -o $(BINARY_DIR)/advanced ./$(EXAMPLES_DIR)/advanced/
	$(GOBUILD) -o $(BINARY_DIR)/microservice ./$(EXAMPLES_DIR)/microservice/

# Run basic example
run-basic: build-examples
	@echo "Running basic example..."
	./$(BINARY_DIR)/basic

# Run advanced example
run-advanced: build-examples
	@echo "Running advanced example..."
	./$(BINARY_DIR)/advanced

# Run microservice example
run-microservice: build-examples
	@echo "Running microservice example..."
	./$(BINARY_DIR)/microservice

# Run all examples
run-examples: run-basic run-advanced run-microservice

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)/
	rm -f $(COVERAGE_FILE)
	rm -f coverage.html

# Run all checks (format, lint, test)
check: fmt lint test

# Run full CI pipeline
ci: deps fmt lint security test-race test-coverage

# Generate documentation
docs:
	@echo "Generating documentation..."
	$(GOCMD) doc -all . > API.md

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Verify dependencies
verify-deps:
	@echo "Verifying dependencies..."
	$(GOMOD) verify

# Show help
help:
	@echo "Available commands:"
	@echo "  install-tools    - Install development tools"
	@echo "  deps            - Install dependencies"
	@echo "  fmt             - Format code"
	@echo "  lint            - Run linter"
	@echo "  security        - Run security scan"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  test-race       - Run tests with race detection"
	@echo "  bench           - Run benchmarks"
	@echo "  build-examples  - Build all examples"
	@echo "  run-basic       - Run basic example"
	@echo "  run-advanced    - Run advanced example"
	@echo "  run-microservice- Run microservice example"
	@echo "  run-examples    - Run all examples"
	@echo "  clean           - Clean build artifacts"
	@echo "  check           - Run format, lint, and test"
	@echo "  ci              - Run full CI pipeline"
	@echo "  docs            - Generate documentation"
	@echo "  update-deps     - Update dependencies"
	@echo "  verify-deps     - Verify dependencies"
	@echo "  help            - Show this help"

# Development workflow targets
dev-setup: install-tools deps
	@echo "Development environment setup complete!"

dev-test: fmt test
	@echo "Development test complete!"

dev-check: fmt lint test
	@echo "Development check complete!"