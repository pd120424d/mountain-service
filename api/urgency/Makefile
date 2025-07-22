# Project variables
SERVICE_NAME := urgency-service
BINARY_NAME := urgency-service

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
	set GOFLAGS=-mod=mod && go generate ./internal/repositories
	@echo "Mocks generated successfully."

# Run tests
test:
	@echo "Running tests for $(SERVICE_NAME)..."
	@go test -v ./...
	@echo "Tests completed."

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@echo "Clean completed."

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/main.go -o cmd/docs
	@echo "Swagger documentation generated."

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatting completed."

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run
	@echo "Linting completed."

# Run all checks (format, lint, test)
check: fmt lint test
	@echo "All checks completed."
