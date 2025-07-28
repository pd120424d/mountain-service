# Global Makefile for Mountain Service Project
# This Makefile orchestrates all services and components

# Project variables
PROJECT_NAME := mountain-service
API_SERVICES := employee urgency activity
UI_DIR := ui
VERSION_SERVICE_DIR := api/version-service

.PHONY: help all build clean test run stop swagger fmt lint check deps docker-build docker-up docker-down ui-build ui-test ui-start version-service-build version-service-run \
	employee-build employee-clean employee-test employee-run employee-generate employee-mocks employee-swagger employee-fmt employee-lint employee-check \
	urgency-build urgency-clean urgency-test urgency-run urgency-generate urgency-mocks urgency-swagger urgency-fmt urgency-lint urgency-check \
	activity-build activity-clean activity-test activity-run activity-generate activity-mocks activity-swagger activity-fmt activity-lint activity-check

# Default target
all: deps generate build test
	@echo " All services built and tested successfully"

# Help target
help:
	@echo "Mountain Service - Global Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  help              - Show this help message"
	@echo "  all               - Build and test all services (default)"
	@echo "  build             - Build all services"
	@echo "  clean             - Clean all services"
	@echo "  test              - Run tests for all services"
	@echo "  run               - Run all services locally"
	@echo "  stop              - Stop all running services"
	@echo "  generate          - Generate code for all services"
	@echo "  mocks             - Generate mocks for all services"
	@echo "  swagger           - Generate Swagger docs for all services"
	@echo "  fmt               - Format code for all services"
	@echo "  lint              - Lint code for all services"
	@echo "  check             - Run format, lint, and test for all services"
	@echo "  deps              - Install dependencies for all services"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-build      - Build Docker images for all services"
	@echo "  docker-up         - Start all services with Docker Compose"
	@echo "  docker-down       - Stop all services with Docker Compose"
	@echo ""
	@echo "UI targets:"
	@echo "  ui-build          - Build the UI"
	@echo "  ui-test           - Test the UI"
	@echo "  ui-start          - Start the UI development server"
	@echo "  ui-models         - Generate UI models from Swagger"
	@echo ""
	@echo "Individual service targets:"
	@echo "  <service>-<target> - Run specific target for specific service"
	@echo "  Available services: $(API_SERVICES) version-service ui"
	@echo "  Example: employee-build, urgency-test, activity-swagger"

# Build all services
build: $(foreach service,$(API_SERVICES),$(service)-build) version-service-build ui-build
	@echo "All services built successfully"

# Clean all services
clean: $(foreach service,$(API_SERVICES),$(service)-clean) version-service-clean ui-clean
	@echo "All services cleaned successfully"

# Test all services
test: $(foreach service,$(API_SERVICES),$(service)-test) ui-test
	@echo "All services tested successfully"

# Generate code for all services
generate: $(foreach service,$(API_SERVICES),$(service)-generate)
	@echo "Code generation completed for all services"

# Generate mocks for all services
mocks: $(foreach service,$(API_SERVICES),$(service)-mocks)
	@echo "Mocks generated for all services"

# Generate Swagger docs for all services
swagger: $(foreach service,$(API_SERVICES),$(service)-swagger)
	@echo "Swagger documentation generated for all services"

# Format code for all services
fmt: $(foreach service,$(API_SERVICES),$(service)-fmt)
	@echo "Code formatted for all services"

# Lint code for all services
lint: $(foreach service,$(API_SERVICES),$(service)-lint)
	@echo "Code linted for all services"

# Run all checks (format, lint, test)
check: $(foreach service,$(API_SERVICES),$(service)-check) ui-test
	@echo "All checks completed successfully"

# Install dependencies for all services
deps: api-deps ui-deps
	@echo "Dependencies installed for all services"

# API dependencies
api-deps:
	@echo "Installing API dependencies..."
	@cd api && go mod tidy && go mod vendor

# Employee service targets
employee-build:
	@echo "Building employee service..."
	@cd api/employee && $(MAKE) build

employee-clean:
	@echo "Cleaning employee service..."
	@cd api/employee && $(MAKE) clean

employee-test:
	@echo "Testing employee service..."
	@cd api/employee && $(MAKE) test

employee-run:
	@echo "Running employee service..."
	@cd api/employee && $(MAKE) run

employee-generate:
	@echo "Generating code for employee service..."
	@cd api/employee && $(MAKE) generate

employee-mocks:
	@echo "Generating mocks for employee service..."
	@cd api/employee && $(MAKE) mocks

employee-swagger:
	@echo "Generating Swagger docs for employee service..."
	@cd api/employee && $(MAKE) swagger

employee-fmt:
	@echo "Formatting employee service..."
	@cd api/employee && $(MAKE) fmt

employee-lint:
	@echo "Linting employee service..."
	@cd api/employee && $(MAKE) lint

employee-check:
	@echo "Running checks for employee service..."
	@cd api/employee && $(MAKE) check

# Urgency service targets
urgency-build:
	@echo "Building urgency service..."
	@cd api/urgency && $(MAKE) build

urgency-clean:
	@echo "Cleaning urgency service..."
	@cd api/urgency && $(MAKE) clean

urgency-test:
	@echo "Testing urgency service..."
	@cd api/urgency && $(MAKE) test

urgency-run:
	@echo "Running urgency service..."
	@cd api/urgency && $(MAKE) run

urgency-generate:
	@echo "Generating code for urgency service..."
	@cd api/urgency && $(MAKE) generate

urgency-mocks:
	@echo "Generating mocks for urgency service..."
	@cd api/urgency && $(MAKE) mocks

urgency-swagger:
	@echo "Generating Swagger docs for urgency service..."
	@cd api/urgency && $(MAKE) swagger

urgency-fmt:
	@echo "Formatting urgency service..."
	@cd api/urgency && $(MAKE) fmt

urgency-lint:
	@echo "Linting urgency service..."
	@cd api/urgency && $(MAKE) lint

urgency-check:
	@echo "Running checks for urgency service..."
	@cd api/urgency && $(MAKE) check

# Activity service targets
activity-build:
	@echo "Building activity service..."
	@cd api/activity && $(MAKE) build

activity-clean:
	@echo "Cleaning activity service..."
	@cd api/activity && $(MAKE) clean

activity-test:
	@echo "Testing activity service..."
	@cd api/activity && $(MAKE) test

activity-run:
	@echo "Running activity service..."
	@cd api/activity && $(MAKE) run

activity-generate:
	@echo "Generating code for activity service..."
	@cd api/activity && $(MAKE) generate

activity-mocks:
	@echo "Generating mocks for activity service..."
	@cd api/activity && $(MAKE) mocks

activity-swagger:
	@echo "Generating Swagger docs for activity service..."
	@cd api/activity && $(MAKE) swagger

activity-fmt:
	@echo "Formatting activity service..."
	@cd api/activity && $(MAKE) fmt

activity-lint:
	@echo "Linting activity service..."
	@cd api/activity && $(MAKE) lint

activity-check:
	@echo "Running checks for activity service..."
	@cd api/activity && $(MAKE) check

# Version service targets (no Makefile, so direct commands)
version-service-build:
	@echo "Building version service..."
	@cd $(VERSION_SERVICE_DIR) && go build -o version-service main.go

version-service-run:
	@echo "Running version service..."
	@cd $(VERSION_SERVICE_DIR) && ./version-service

version-service-clean:
	@echo "Cleaning version service..."
	@cd $(VERSION_SERVICE_DIR) && rm -f version-service

# UI targets
ui-build:
	@echo "Building UI..."
	@cd $(UI_DIR) && npm run build

ui-test:
	@echo "Testing UI..."
	@cd $(UI_DIR) && npm run test

ui-start:
	@echo "Starting UI development server..."
	@cd $(UI_DIR) && npm start

ui-models:
	@echo "Generating UI models from Swagger..."
	@cd $(UI_DIR) && npm run models:generate

ui-deps:
	@echo "Installing UI dependencies..."
	@cd $(UI_DIR) && npm install

ui-clean:
	@echo "Cleaning UI build artifacts..."
	@cd $(UI_DIR) && rm -rf dist node_modules/.cache

# Docker targets
docker-build:
	@echo "Building Docker images..."
	@docker compose --env-file .env.staging build

docker-up:
	@echo "Starting services with Docker Compose (staging)..."
	@docker compose --env-file .env.staging up -d

docker-down:
	@echo "Stopping services with Docker Compose..."
	@docker compose down

docker-logs:
	@echo "Showing Docker Compose logs..."
	@docker compose logs -f

# Production Docker targets
docker-prod-up:
	@echo "Starting services with Docker Compose (production)..."
	@docker compose --env-file .env.aws up -d

docker-prod-down:
	@echo "Stopping services with Docker Compose (production)..."
	@docker compose down

# Development workflow targets
dev-setup: deps generate swagger ui-models
	@echo "Development environment setup completed"

dev-start: docker-up ui-start
	@echo "Development environment started"

dev-stop: docker-down
	@echo "Development environment stopped"

# Health check
health:
	@echo "Checking service health..."
	@curl -f http://localhost:8082/api/v1/health || echo "$(RED)Employee service not healthy"
	@curl -f http://localhost:8083/api/v1/health || echo "$(RED)Urgency service not healthy"
	@curl -f http://localhost:8084/api/v1/health || echo "$(RED)Activity service not healthy"
	@curl -f http://localhost:8090/api/v1/health || echo "$(RED)Version service not healthy"

# Coverage report
coverage:
	@echo "Generating coverage reports..."
	@./backend-test-cover.sh
	@./frontend-test-cover.sh
	@echo "Coverage reports generated"
