# Makefile for Image Metadata Viewer

# Variables
BINARY_NAME=image-metadata-viewer
MAIN_PATH=./src/cmd/server
BUILD_DIR=./build
GO=go
GOFLAGS=-v

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
.PHONY: help
help:
	@echo "Image Metadata Viewer - Makefile Commands"
	@echo "==========================================="
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## run: Run the application in development mode
.PHONY: run
run:
	@echo "Starting server..."
	@$(GO) run $(MAIN_PATH)/main.go

## build: Build the application binary
.PHONY: build
build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)/main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## build-all: Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-windows build-darwin
	@echo "All builds complete!"

## build-linux: Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)/main.go

## build-windows: Build for Windows
.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)/main.go

## build-darwin: Build for macOS
.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)/main.go
	@GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)/main.go

## test: Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@$(GO) test -v ./...

## test-coverage: Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## clean: Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

## deps: Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "Dependencies updated!"

## fmt: Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...
	@echo "Formatting complete!"

## lint: Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

## vet: Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	@$(GO) vet ./...

## install: Install the application
.PHONY: install
install:
	@echo "Installing application..."
	@$(GO) install $(MAIN_PATH)/main.go
	@echo "Installation complete!"

## docker-build: Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .
	@echo "Docker image built!"

## docker-run: Run Docker container
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(BINARY_NAME):latest

## dev: Run in development mode with auto-reload (requires air)
.PHONY: dev
dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	@air

## check: Run all checks (fmt, vet, test)
.PHONY: check
check: fmt vet test
	@echo "All checks passed!"
