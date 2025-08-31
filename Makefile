.PHONY: build run clean proto test sdk-build sdk-publish sdk-dev

# Variables
BINARY_NAME=watchdog-server
BINARY_PATH=./bin/$(BINARY_NAME)
PROTO_DIR=proto
API_DIR=api
SDK_JS_DIR=sdk/javascript

# Build the server
build:
	go build -o $(BINARY_PATH) ./cmd/main.go

# Run the server
run: build
	$(BINARY_PATH)

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Generate protobuf code (requires protoc to be installed)
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/watchdog.proto

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Setup environment file
env-setup:
	@if [ ! -f .env ]; then \
		cp .env.default .env; \
		echo "Created .env file from .env.default"; \
		echo "Please edit .env with your configuration"; \
	else \
		echo ".env file already exists"; \
	fi

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/main.go

# Database operations (Ent-based)
db-migrate-ent:
	@echo "Running Ent-based database migration..."
	go run scripts/migrate-ent.go

db-migrate-ent-dry:
	@echo "Running Ent migration in dry-run mode..."
	go run scripts/migrate-ent.go -dry-run

# Generate Ent code
ent-generate:
	@echo "Generating Ent code from schema..."
	go run ent/generate.go

# Deployment
deploy-dev:
	./scripts/deploy.sh development

deploy-prod:
	./scripts/deploy.sh production

# JavaScript SDK operations
sdk-install:
	@echo "Installing JavaScript SDK dependencies..."
	cd $(SDK_JS_DIR) && npm install

sdk-build: sdk-install
	@echo "Building JavaScript SDK..."
	cd $(SDK_JS_DIR) && npm run build

sdk-test: sdk-build
	@echo "Testing JavaScript SDK..."
	cd $(SDK_JS_DIR) && npm test

sdk-lint: sdk-install
	@echo "Linting JavaScript SDK..."
	cd $(SDK_JS_DIR) && npm run lint

sdk-publish: sdk-build sdk-test
	@echo "Publishing JavaScript SDK to npm registry..."
	cd $(SDK_JS_DIR) && npm publish

sdk-publish-beta: sdk-build sdk-test
	@echo "Publishing JavaScript SDK as beta version..."
	cd $(SDK_JS_DIR) && npm publish --tag beta

sdk-version-patch: sdk-install
	@echo "Bumping patch version..."
	cd $(SDK_JS_DIR) && npm version patch

sdk-version-minor: sdk-install
	@echo "Bumping minor version..."
	cd $(SDK_JS_DIR) && npm version minor

sdk-version-major: sdk-install
	@echo "Bumping major version..."
	cd $(SDK_JS_DIR) && npm version major

sdk-release: sdk-build sdk-test
	@echo "Running automated SDK release..."
	./scripts/sdk-release.sh

sdk-release-patch: sdk-build sdk-test
	@echo "Releasing patch version..."
	./scripts/sdk-release.sh patch

sdk-release-minor: sdk-build sdk-test
	@echo "Releasing minor version..."
	./scripts/sdk-release.sh minor

sdk-release-major: sdk-build sdk-test
	@echo "Releasing major version..."
	./scripts/sdk-release.sh major

sdk-clean:
	@echo "Cleaning JavaScript SDK build artifacts..."
	cd $(SDK_JS_DIR) && rm -rf dist/ node_modules/ src/generated/

# Help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Server Commands:"
	@echo "  build           - Build the server binary"
	@echo "  run             - Build and run the server"
	@echo "  clean           - Clean build artifacts"
	@echo "  proto           - Generate protobuf code"
	@echo "  test            - Run tests"
	@echo "  deps            - Install dependencies"
	@echo "  env-setup       - Create .env file from template"
	@echo "  fmt             - Format code"
	@echo "  lint            - Run linter"
	@echo "  build-all       - Build for multiple platforms"
	@echo ""
	@echo "Database Commands:"
	@echo "  db-migrate-ent     - Run Ent-based database migration (recommended)"
	@echo "  db-migrate-ent-dry - Show what migration would do without executing"
	@echo "  ent-generate       - Generate Ent code from schema"
	@echo ""
	@echo "Deployment Commands:"
	@echo "  deploy-dev      - Deploy for development"
	@echo "  deploy-prod     - Deploy for production"
	@echo ""
	@echo "JavaScript SDK Commands:"
	@echo "  sdk-install        - Install SDK dependencies"
	@echo "  sdk-build          - Build the JavaScript SDK"
	@echo "  sdk-test           - Run SDK tests"
	@echo "  sdk-lint           - Lint SDK code"
	@echo "  sdk-clean          - Clean SDK build artifacts"
	@echo ""
	@echo "SDK Publishing Commands:"
	@echo "  sdk-publish        - Publish SDK to npm registry"
	@echo "  sdk-publish-beta   - Publish SDK as beta version"
	@echo "  sdk-release        - Automated release (patch)"
	@echo "  sdk-release-patch  - Release patch version (1.0.x)"
	@echo "  sdk-release-minor  - Release minor version (1.x.0)"
	@echo "  sdk-release-major  - Release major version (x.0.0)"