# Step 1: Build
FROM golang:1.23-alpine AS build

# Build arguments to bust cache when needed
ARG BUILD_TIMESTAMP
ARG GIT_COMMIT_SHA
ARG SWAGGER_HOST=localhost:8082

WORKDIR /app

# Install swag for swagger generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the employee service, shared code, and contracts (relative to context ./api)
COPY employee/ ./employee/
COPY shared/ ./shared/
COPY contracts/ ./contracts/

# Generate swagger documentation with correct host
WORKDIR /app/employee
ENV SWAGGER_HOST=${SWAGGER_HOST}
RUN swag init -g cmd/main.go -o cmd/docs --exclude cmd/docs --pdl 3

# Build the employee service
WORKDIR /app
RUN go build -v -o /employee-service ./employee/cmd/main.go

# Step 2: Slim runtime image
FROM alpine:latest

# Install wget for health checks
RUN apk --no-cache add wget

WORKDIR /root/
COPY --from=build /employee-service .
COPY --from=build /app/employee/cmd/docs /docs

EXPOSE 8082
CMD ["./employee-service"]
