# TMDB CLI Enhanced Makefile
# =========================
#
# This Makefile provides a comprehensive set of targets for developing,
# testing, building, and releasing the TMDB CLI application.
#
# GOLANGCI-LINT VERSION: v1.64.8 (Last stable V1 version)
# ========================================================
#
# We use golangci-lint v1.64.8, which is the last stable V1 version.
# V2 behaves weirdly with GitHub Actions.
#
# WORKFLOWS OVERVIEW:
# ==================
#
# 📈 DEVELOPMENT WORKFLOW:
#   1. make deps        - Install and update dependencies
#   2. make dev         - Complete development cycle (fmt, vet, test, build)
#   3. make run ARGS='' - Test your changes locally
#
# 🔍 QUALITY ASSURANCE WORKFLOW:
#   1. make check       - Run all quality checks (fmt, vet, lint, test, security)
#   2. make coverage    - Generate and view detailed coverage reports
#   3. make benchmark   - Performance testing and optimization
#
# 🚀 RELEASE WORKFLOW:
#   1. make pre-release    - Complete pre-release validation
#   2. git tag v1.0.0      - Create release tag
#   3. git push --tags     - Trigger automated release via GitHub Actions
#   4. make release        - Manual release (alternative to GitHub Actions)
#
# 🧪 TESTING WORKFLOWS:
#   - Unit Tests:        make test
#   - Integration Tests: make test-integration (requires TMDB_API_KEY)
#   - E2E Tests:         make test-e2e (requires TMDB_API_KEY)
#   - All Tests:         make check
#
# 🔧 BUILD WORKFLOWS:
#   - Single Platform:   make build
#   - Multi-Platform:    make build-all
#   - Release Build:     make release-snapshot (test release without publish)
#
# GITHUB ACTIONS INTEGRATION:
# ===========================
#
# This Makefile is designed to work seamlessly with GitHub Actions:
#
# CI Pipeline (.github/workflows/ci.yml):
#   - Triggered on: push to main/dev, pull requests
#   - Uses: make tools, make check, make test-integration
#   - Purpose: Continuous validation of code quality
#
# Release Pipeline (.github/workflows/release.yml):
#   - Triggered on: git tag push (v*.*.*)
#   - Uses: make tools, make release-check, make release
#   - Purpose: Automated releases with GoReleaser
#
# LOCAL VS CI CONSISTENCY:
# ========================
#
# Commands work identically in both environments:
#   Local:  make check        CI:  make check
#   Local:  make test         CI:  make test
#   Local:  make build        CI:  make build
#
# This ensures "works on my machine" problems are minimized.
#
# DEPENDENCIES:
# ============
#
# Required:
#   - Go 1.23.5+ (golang.org)
#   - Git (for version information)
#
# Auto-installed via 'make tools':
#   - GoReleaser (goreleaser.com)
#   - golangci-lint v1.64.8 (golangci-lint.run)
#   - gosec (securecodewarrior.github.io/gosec)
#
# Optional (for specific workflows):
#   - TMDB_API_KEY (for integration/e2e tests)
#   - bc command (for coverage threshold checking)
#
# QUICK START:
# ===========
#
#   make deps     # Install dependencies
#   make dev      # Development build
#   make check    # Quality validation
#   make tools    # Install release tools
#
# For detailed help: make help
#
# =============================================================================

.PHONY: all build test clean install run fmt vet lint help
.PHONY: release release-snapshot release-test tools deps check
.PHONY: coverage benchmark security

# Variables
BINARY_NAME=tmdb-cli
CMD_DIR=./cmd
BUILD_DIR=./bin
DIST_DIR=./dist
COVERAGE_DIR=./coverage

# Tool versions
GOLANGCI_LINT_VERSION=v1.64.8

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go build flags
LDFLAGS=-ldflags="-s -w -X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE) -X main.gitCommit=$(GIT_COMMIT)"
BUILD_FLAGS=-trimpath $(LDFLAGS)

# Go environment
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

# Default target
all: deps fmt vet lint test build

## Development targets

# Install dependencies
deps:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	@go mod download
	@go mod tidy

# Format code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	@gofmt -s -w .
	@go mod tidy

# Vet code
vet:
	@echo "$(BLUE)Vetting code...$(NC)"
	@go vet ./...

# Lint code (requires golangci-lint v1.64.8)
lint:
	@echo "$(BLUE)Linting code with golangci-lint $(GOLANGCI_LINT_VERSION)...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint --version | grep -q "$(GOLANGCI_LINT_VERSION)" || \
		(echo "$(YELLOW)Warning: Different golangci-lint version detected. Expected $(GOLANGCI_LINT_VERSION)$(NC)"); \
		golangci-lint run --config .golangci.yml; \
	else \
		echo "$(YELLOW)golangci-lint not found. Install with 'make tools'$(NC)"; \
		exit 1; \
	fi

# Run tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out ./internal/... ./cmd/...

# Run tests with coverage report
coverage: test
	@echo "$(BLUE)Generating coverage report...$(NC)"
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | grep total:
	@echo "$(GREEN)Coverage report generated: $(COVERAGE_DIR)/coverage.html$(NC)"

# Generate coverage report excluding test packages
coverage-clean:
	@echo "$(BLUE)Running tests and generating clean coverage report...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out ./internal/... ./cmd/...
	@echo "$(BLUE)Filtering out test files from coverage...$(NC)"
	@grep -v "_test.go\|tests/" $(COVERAGE_DIR)/coverage.out > $(COVERAGE_DIR)/coverage-clean.out || true
	@go tool cover -html=$(COVERAGE_DIR)/coverage-clean.out -o $(COVERAGE_DIR)/coverage-clean.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage-clean.out | grep total:
	@echo "$(GREEN)Clean coverage report generated: $(COVERAGE_DIR)/coverage-clean.html$(NC)"

# Coverage with threshold check (excluding tests)
coverage-check: coverage-clean
	@echo "$(BLUE)Checking coverage threshold...$(NC)"
	@COVERAGE=$$(go tool cover -func=$(COVERAGE_DIR)/coverage-clean.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	echo "Coverage: $$COVERAGE%"; \
	if command -v bc >/dev/null 2>&1; then \
		if [ $$(echo "$$COVERAGE < 70" | bc -l) -eq 1 ]; then \
			echo "$(RED)Error: Coverage $$COVERAGE% is below minimum threshold of 70%$(NC)"; \
			exit 1; \
		else \
			echo "$(GREEN)Coverage $$COVERAGE% meets minimum threshold$(NC)"; \
		fi; \
	else \
		echo "$(YELLOW)bc not available, skipping coverage threshold check$(NC)"; \
	fi

# Run integration tests (requires TMDB_API_KEY)
test-integration:
	@echo "$(BLUE)Running integration tests...$(NC)"
	@if [ -z "$$TMDB_API_KEY" ]; then \
		echo "$(YELLOW)TMDB_API_KEY not set, skipping integration tests$(NC)"; \
	else \
		go test -v -tags=integration ./tests/integration/...; \
	fi

# Run end-to-end tests (requires TMDB_API_KEY)
test-e2e: build
	@echo "$(BLUE)Running E2E tests...$(NC)"
	@if [ -z "$$TMDB_API_KEY" ]; then \
		echo "$(YELLOW)TMDB_API_KEY not set, skipping E2E tests$(NC)"; \
	else \
		go test -v -tags=e2e ./tests/e2e/...; \
	fi

# Run benchmarks
benchmark:
	@echo "$(BLUE)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...

# Security scan (requires gosec)
security:
	@echo "$(BLUE)Running security scan...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "$(YELLOW)gosec not found. Install with 'make tools'$(NC)"; \
	fi

## Build targets

# Build the binary
build:
	@echo "$(BLUE)Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Build for multiple platforms
build-all:
	@echo "$(BLUE)Building for multiple platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ] && [ "$$arch" = "arm64" ]; then continue; fi; \
			ext=""; \
			if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
			echo "Building for $$os/$$arch..."; \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch \
				go build $(BUILD_FLAGS) \
				-o $(BUILD_DIR)/$(BINARY_NAME)-$$os-$$arch$$ext $(CMD_DIR); \
		done; \
	done
	@echo "$(GREEN)Multi-platform build complete$(NC)"

## Release targets (GoReleaser)

# Install tools with specific versions
tools:
	@echo "$(BLUE)Installing development tools...$(NC)"
	@if ! command -v goreleaser >/dev/null 2>&1; then \
		echo "Installing GoReleaser..."; \
		go install github.com/goreleaser/goreleaser@latest; \
	fi
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
	else \
		CURRENT_VERSION=$$(golangci-lint --version | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -1); \
		if [ "$$CURRENT_VERSION" != "$(GOLANGCI_LINT_VERSION)" ]; then \
			echo "Updating golangci-lint from $$CURRENT_VERSION to $(GOLANGCI_LINT_VERSION)..."; \
			curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
		else \
			echo "golangci-lint $(GOLANGCI_LINT_VERSION) already installed"; \
		fi; \
	fi
	@if ! command -v gosec >/dev/null 2>&1; then \
		echo "Installing gosec..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
	fi
	@echo "$(GREEN)Tools installation complete$(NC)"

# Validate GoReleaser configuration
release-check:
	@echo "$(BLUE)Validating GoReleaser configuration...$(NC)"
	@goreleaser check

# Test release without publishing (snapshot)
release-snapshot: clean
	@echo "$(BLUE)Creating snapshot release...$(NC)"
	@goreleaser release --snapshot --clean

# Test release build without publishing
release-test: clean
	@echo "$(BLUE)Testing release build...$(NC)"
	@goreleaser build --snapshot --clean

# Create a full release (requires clean git state and tag)
release:
	@echo "$(BLUE)Creating release...$(NC)"
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "$(RED)Error: Working directory is not clean$(NC)"; \
		exit 1; \
	fi
	@goreleaser release --clean

## Installation and utility targets

# Install binary to GOPATH/bin
install: build
	@echo "$(BLUE)Installing $(BINARY_NAME)...$(NC)"
	@go install $(BUILD_FLAGS) $(CMD_DIR)
	@echo "$(GREEN)$(BINARY_NAME) installed to $$(go env GOPATH)/bin$(NC)"

# Run the application
run: build
	@echo "$(BLUE)Running $(BINARY_NAME)...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Run with popular movies example
run-example: build
	@echo "$(BLUE)Running example: popular movies$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME) popular 5

# Clean build artifacts
clean:
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@rm -rf $(COVERAGE_DIR)
	@go clean

# Deep clean (including module cache)
clean-all: clean
	@echo "$(BLUE)Deep cleaning...$(NC)"
	@go clean -modcache
	@go clean -cache

# Check code quality with updated linter
check: deps fmt vet lint test security
	@echo "$(GREEN)All quality checks passed with golangci-lint $(GOLANGCI_LINT_VERSION)!$(NC)"

# Prepare for release (run all checks)
pre-release: clean check test-integration release-test
	@echo "$(GREEN)Ready for release!$(NC)"

# Development workflow
dev: deps fmt vet test build
	@echo "$(GREEN)Development build complete!$(NC)"

# Show version information
version:
	@echo "Version: $(VERSION)"
	@echo "Build Date: $(BUILD_DATE)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $$(go version)"
	@echo "golangci-lint Version: $(GOLANGCI_LINT_VERSION)"

# Show help
help:
	@echo "$(BLUE)TMDB CLI Makefile - golangci-lint $(GOLANGCI_LINT_VERSION)$(NC)"
	@echo ""
	@echo "$(YELLOW)Development:$(NC)"
	@echo "  deps              Install dependencies"
	@echo "  fmt               Format code"
	@echo "  vet               Vet code"
	@echo "  lint              Lint code with golangci-lint $(GOLANGCI_LINT_VERSION)"
	@echo "  test              Run unit tests"
	@echo "  test-integration  Run integration tests (requires TMDB_API_KEY)"
	@echo "  test-e2e          Run end-to-end tests"
	@echo "  coverage          Generate test coverage report"
	@echo "  coverage-clean    Generate coverage report excluding test files"
	@echo "  coverage-check    Check coverage threshold (excluding tests)"
	@echo "  benchmark         Run benchmarks"
	@echo "  security          Run security scan (requires gosec)"
	@echo "  check             Run all quality checks"
	@echo "  dev               Development workflow (deps, fmt, vet, test, build)"
	@echo ""
	@echo "$(YELLOW)Building:$(NC)"
	@echo "  build             Build binary for current platform"
	@echo "  build-all         Build for multiple platforms"
	@echo "  install           Install binary to GOPATH/bin"
	@echo "  run ARGS='...'    Build and run with arguments"
	@echo "  run-example       Run with example command"
	@echo ""
	@echo "$(YELLOW)Release:$(NC)"
	@echo "  tools             Install development tools (golangci-lint $(GOLANGCI_LINT_VERSION))"
	@echo "  release-check     Validate GoReleaser configuration"
	@echo "  release-snapshot  Create snapshot release (no publish)"
	@echo "  release-test      Test release build"
	@echo "  release           Create full release (requires tag)"
	@echo "  pre-release       Run all checks before release"
	@echo ""
	@echo "$(YELLOW)Utility:$(NC)"
	@echo "  clean             Clean build artifacts"
	@echo "  clean-all         Deep clean (including caches)"
	@echo "  version           Show version information"
	@echo "  help              Show this help"
	@echo ""
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  make dev                    # Development workflow"
	@echo "  make run ARGS='popular 10'  # Run with arguments"
	@echo "  make test-integration       # Integration tests"
	@echo "  make release-snapshot       # Test release"
	@echo ""
	@echo "$(YELLOW)golangci-lint Info:$(NC)"
	@echo "  Version: $(GOLANGCI_LINT_VERSION) (Last stable V1 version)"
	@echo "  Config:  .golangci.yml"
	@echo "  Action:  golangci/golangci-lint-action@v6"
