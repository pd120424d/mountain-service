# Step 1: Build
FROM golang:1.23-alpine AS build

# Build arguments to bust cache when needed
ARG BUILD_TIMESTAMP
ARG GIT_COMMIT_SHA

WORKDIR /app

# Copy module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the urgency service, shared code, and contracts (relative to context ./api)
COPY urgency/ ./urgency/
COPY shared/ ./shared/
COPY contracts/ ./contracts/

# Copy Swagger docs
COPY urgency/cmd/docs/ ./docs/

# Build the urgency service
RUN go build -v -o /urgency-service ./urgency/cmd/main.go

# Step 2: Slim runtime image
FROM alpine:latest

# Install wget for health checks
RUN apk --no-cache add wget

WORKDIR /root/
COPY --from=build /urgency-service .
COPY --from=build /app/docs /docs

EXPOSE 8083
CMD ["./urgency-service"]
