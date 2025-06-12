# TMDB CLI Makefile
.PHONY: all build test clean install run fmt vet

# Variables
BINARY_NAME=tmdb
CMD_DIR=./cmd
MAIN_FILE=$(CMD_DIR)/main.go
BUILD_DIR=./bin

# Default target
all: fmt vet test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean

# Install binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	go install $(CMD_DIR)

# Run the application
run: build
	@$(BUILD_DIR)/$(BINARY_NAME)

# Development workflow
dev: fmt vet test

# Show help
help:
	@echo "Available targets:"
	@echo "  all      - Format, vet, test, and build"
	@echo "  build    - Build the binary"
	@echo "  test     - Run tests"
	@echo "  fmt      - Format code"
	@echo "  vet      - Vet code"
	@echo "  clean    - Clean build artifacts"
	@echo "  install  - Install binary"
	@echo "  run      - Build and run"
	@echo "  dev      - Development workflow (fmt, vet, test)"
	@echo "  help     - Show this help"
