# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Binary names
BINARY_NAME=yourapp
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME)_windows.exe
BINARY_DARWIN=$(BINARY_NAME)_darwin

# Directories
CMD_DIR=./cmd/server
BUILD_DIR=./build
DIST_DIR=./dist
COVERAGE_DIR=./coverage

# Version info
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"
BUILD_FLAGS=-trimpath -ldflags "-s -w"

# Default target
.PHONY: all
all: clean deps lint test build

# Help target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Clean targets
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -rf $(COVERAGE_DIR)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -f $(BINARY_DARWIN)

# Dependencies
.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOMOD) tidy
	$(GOGET) -u ./...

# Linting
.PHONY: lint
lint: ## Run linters
	@echo "Running linters..."
	$(GOLINT) run

.PHONY: lint-fix
lint-fix: ## Run linters with auto-fix
	@echo "Running linters with auto-fix..."
	$(GOLINT) run --fix

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	$(GOCMD) fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Testing
.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	$(GOTEST) -race -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -v -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"

.PHONY: test-benchmark
test-benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Building
.PHONY: build
build: ## Build the application
	@echo "Building application..."
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(CMD_DIR)

.PHONY: build-linux
build-linux: ## Build for Linux
	@echo "Building for Linux..."
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_UNIX) -v $(CMD_DIR)

.PHONY: build-windows
build-windows: ## Build for Windows
	@echo "Building for Windows..."
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_WINDOWS) -v $(CMD_DIR)

.PHONY: build-darwin
build-darwin: ## Build for macOS
	@echo "Building for macOS..."
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_DARWIN) -v $(CMD_DIR)

.PHONY: build-all
build-all: build-linux build-windows build-darwin ## Build for all platforms

# Distribution
.PHONY: dist
dist: clean build-all ## Create distribution packages
	@echo "Creating distribution packages..."
	mkdir -p $(DIST_DIR)
	cp $(BUILD_DIR)/$(BINARY_UNIX) $(DIST_DIR)/
	cp $(BUILD_DIR)/$(BINARY_WINDOWS) $(DIST_DIR)/
	cp $(BUILD_DIR)/$(BINARY_DARWIN) $(DIST_DIR)/
	cp configs/config.example.yaml $(DIST_DIR)/config.yaml
	@echo "Distribution packages created in $(DIST_DIR)/"

# Development
.PHONY: run
run: ## Run the application
	@echo "Running application..."
	$(GOCMD) run $(CMD_DIR)/main.go

.PHONY: run-dev
run-dev: ## Run the application in development mode
	@echo "Running application in development mode..."
	APP_ENV=development APP_LOG_LEVEL=debug $(GOCMD) run $(CMD_DIR)/main.go

.PHONY: run-prod
run-prod: ## Run the application in production mode
	@echo "Running application in production mode..."
	APP_ENV=production APP_LOG_LEVEL=info $(GOCMD) run $(CMD_DIR)/main.go

# Docker
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .
	docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 $(BINARY_NAME):latest

.PHONY: docker-push
docker-push: ## Push Docker image
	@echo "Pushing Docker image..."
	docker push $(BINARY_NAME):$(VERSION)
	docker push $(BINARY_NAME):latest

# Security
.PHONY: security
security: ## Run security checks
	@echo "Running security checks..."
	$(GOCMD) list -json -deps ./... | nancy sleuth

.PHONY: audit
audit: ## Audit dependencies
	@echo "Auditing dependencies..."
	$(GOCMD) list -json -deps ./... | nancy sleuth

# Code generation
.PHONY: generate
generate: ## Generate code
	@echo "Generating code..."
	$(GOCMD) generate ./...

# Install tools
.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOCMD) install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	$(GOCMD) install github.com/securecodewarrior/nancy@latest
	$(GOCMD) install github.com/air-verse/air@latest

# Pre-commit hooks
.PHONY: pre-commit
pre-commit: fmt lint test ## Run pre-commit checks
	@echo "Running pre-commit checks..."

# CI/CD
.PHONY: ci
ci: deps lint test build ## Run CI pipeline
	@echo "Running CI pipeline..."

# Release
.PHONY: release
release: clean test build-all dist ## Create release
	@echo "Creating release..."
	@echo "Version: $(VERSION)"
	@echo "Build time: $(BUILD_TIME)"
	@echo "Git commit: $(GIT_COMMIT)"

# Development server with hot reload
.PHONY: dev
dev: ## Run development server with hot reload
	@echo "Starting development server with hot reload..."
	air

# Database operations
.PHONY: db-migrate
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	$(GOCMD) run $(CMD_DIR)/main.go migrate

.PHONY: db-seed
db-seed: ## Seed database
	@echo "Seeding database..."
	$(GOCMD) run $(CMD_DIR)/main.go seed

# Monitoring
.PHONY: monitor
monitor: ## Run monitoring tools
	@echo "Starting monitoring..."
	@echo "Application metrics available at: http://localhost:8080/metrics"
	@echo "Health check available at: http://localhost:8080/health"

# Documentation
.PHONY: docs
docs: ## Generate documentation
	@echo "Generating documentation..."
	$(GOCMD) doc -all ./... > docs/api.md

.PHONY: docs-serve
docs-serve: ## Serve documentation
	@echo "Serving documentation..."
	@echo "Documentation available at: http://localhost:6060"
	godoc -http=:6060

# Cleanup
.PHONY: clean-deps
clean-deps: ## Clean dependencies
	@echo "Cleaning dependencies..."
	$(GOCMD) clean -modcache

.PHONY: clean-all
clean-all: clean clean-deps ## Clean everything
	@echo "Cleaning everything..."

# Show version info
.PHONY: version
version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Build time: $(BUILD_TIME)"
	@echo "Git commit: $(GIT_COMMIT)"
	@echo "Go version: $(shell $(GOCMD) version)"
