# Step 1: Build
FROM golang:1.25.1-alpine AS build

# Build arguments to bust cache when needed
ARG BUILD_TIMESTAMP
ARG GIT_COMMIT_SHA
ARG SWAGGER_HOST=localhost:8084

WORKDIR /app

# Install swag for swagger generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the activity service, shared code, and contracts (relative to context ./api)
COPY activity/ ./activity/
COPY shared/ ./shared/
COPY contracts/ ./contracts/

# Generate swagger documentation with correct host
WORKDIR /app/activity
ENV SWAGGER_HOST=${SWAGGER_HOST}
RUN swag init -g cmd/main.go -o cmd/docs --exclude cmd/docs --pdl 3

# Build the activity service
WORKDIR /app
RUN go build -v -o /activity-service ./activity/cmd/main.go

# Step 2: Slim runtime image
FROM alpine:latest

# Install wget for health checks
RUN apk --no-cache add wget

WORKDIR /root/
COPY --from=build /activity-service .
COPY --from=build /app/activity/cmd/docs /docs

EXPOSE 8084
CMD ["./activity-service"]
