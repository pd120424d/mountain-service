# Step 1: Build
FROM golang:1.23-alpine AS build

# Build arguments to bust cache when needed
ARG BUILD_TIMESTAMP
ARG GIT_COMMIT_SHA

WORKDIR /app

# Copy module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the activity-readmodel-updater service, shared code, and contracts (relative to context ./api)
COPY activity-readmodel-updater/ ./activity-readmodel-updater/
COPY shared/ ./shared/
COPY contracts/ ./contracts/

# Build the activity-readmodel-updater service
RUN go build -v -o /activity-readmodel-updater ./activity-readmodel-updater/cmd/main.go

# Step 2: Slim runtime image
FROM alpine:latest

# Install wget for health checks and ca-certificates for HTTPS
RUN apk --no-cache add wget ca-certificates

WORKDIR /root/
COPY --from=build /activity-readmodel-updater .

EXPOSE 8090
CMD ["./activity-readmodel-updater"]
