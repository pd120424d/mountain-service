# Step 1: Build
FROM golang:1.23-alpine AS build

WORKDIR /app

# Copy module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the employee service and shared code (relative to context ./api)
COPY employee/ ./employee/
COPY shared/ ./shared/

# Copy Swagger docs
COPY employee/cmd/docs/ ./docs/

# Build the employee service
RUN go build -v -o /employee-service ./employee/cmd/main.go

# Step 2: Slim runtime image
FROM alpine:latest

WORKDIR /root/
COPY --from=build /employee-service .
COPY --from=build /app/docs /docs

EXPOSE 8082
CMD ["./employee-service"]
