# Project variables
SERVICE_NAME := employee-service
BINARY_NAME := employee-service

.PHONY: all build generate mocks test clean run

# Default target that will build the service
all: generate build test

# Build the service binary
build:
	@echo "Building the $(SERVICE_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd
	@echo "$(SERVICE_NAME) built successfully."

# Run the service
run:
	@echo "Running the $(SERVICE_NAME)..."
	@./$(BINARY_NAME)

# Generate code (mocks, etc.)
generate:
	@echo "Generating code..."
	@go generate ./...
	@echo "Code generation complete."

# Generate mocks
mocks:
	@echo "Generating mocks..."
	set GOFLAGS=-mod=mod && go generate ./internal/
	@echo "Mocks generated successfully."

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests completed."

# Clean up generated files
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@find . -name 'mock_*.go' -delete
	@echo "Cleanup completed."

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod vendor
	@echo "Dependencies installed."

# Lint the code (requires golangci-lint to be installed)
lint:
	@echo "Linting the code..."
	@golangci-lint run ./...
	@echo "Linting complete."

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@if [ -n "$(SWAGGER_HOST)" ]; then \
		echo "Using SWAGGER_HOST: $(SWAGGER_HOST)"; \
		SWAGGER_HOST=$(SWAGGER_HOST) ~/go/bin/swag init -g cmd/main.go -o cmd/docs --exclude cmd/docs --pdl 3; \
	else \
		echo "No SWAGGER_HOST set, generating with empty host for local development"; \
		~/go/bin/swag init -g cmd/main.go -o cmd/docs --exclude cmd/docs --pdl 3; \
	fi
	@echo "Swagger documentation generated."

# Generate Swagger for production
swagger-prod:
	@echo "Generating Swagger documentation for production..."
	@SWAGGER_HOST=mountain-service.duckdns.org ~/go/bin/swag init -g cmd/main.go -o cmd/docs --exclude cmd/docs --pdl 3
	@echo "Production Swagger documentation generated."
