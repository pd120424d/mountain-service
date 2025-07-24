# Step 1: Build
FROM golang:1.23-alpine AS build

WORKDIR /app

# Copy module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the activity service, shared code, and contracts (relative to context ./api)
COPY activity/ ./activity/
COPY shared/ ./shared/
COPY contracts/ ./contracts/

# Copy Swagger docs
COPY activity/cmd/docs/ ./docs/

# Build the activity service
RUN go build -v -o /activity-service ./activity/cmd/main.go

# Step 2: Slim runtime image
FROM alpine:latest

# Install wget for health checks
RUN apk --no-cache add wget

WORKDIR /root/
COPY --from=build /activity-service .
COPY --from=build /app/docs /docs

EXPOSE 8083
CMD ["./activity-service"]
