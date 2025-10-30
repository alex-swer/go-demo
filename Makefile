.PHONY: help build run test test-coverage test-race clean lint fmt install-tools

help:
	@echo "Go Demo Project - Available Commands:"
	@echo ""
	@echo "  make build         - Build the main application"
	@echo "  make run           - Run the main application"
	@echo "  make test          - Run all tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make test-race     - Run tests with race detector"
	@echo "  make lint          - Run golangci-lint"
	@echo "  make fmt           - Format code"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make install-tools - Install development tools"
	@echo ""

build:
	@echo "Building application..."
	@go build -o bin/demo cmd/demo/main.go
	@echo "Build complete: bin/demo"

run:
	@echo "Running application..."
	@go run cmd/demo/main.go

test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race:
	@echo "Running tests with race detector..."
	@go test -race ./...

lint:
	@echo "Running linter..."
	@golangci-lint run

fmt:
	@echo "Formatting code..."
	@gofmt -w .
	@goimports -w .

clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -f *.exe
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Tools installed successfully"

.DEFAULT_GOAL := help

