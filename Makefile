.PHONY: build run clean proto test

# Variables
BINARY_NAME=watchdog-server
BINARY_PATH=./bin/$(BINARY_NAME)
PROTO_DIR=proto
API_DIR=api

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

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the server binary"
	@echo "  run          - Build and run the server"
	@echo "  clean        - Clean build artifacts"
	@echo "  proto        - Generate protobuf code"
	@echo "  test         - Run tests"
	@echo "  deps         - Install dependencies"
	@echo "  env-setup    - Create .env file from template"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  build-all       - Build for multiple platforms"
	@echo "  db-migrate-ent  - Run Ent-based database migration (recommended)"
	@echo "  db-migrate-ent-dry - Show what migration would do without executing"
	@echo "  ent-generate    - Generate Ent code from schema"
	@echo "  deploy-dev      - Deploy for development"
	@echo "  deploy-prod     - Deploy for production"