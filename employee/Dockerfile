# Step 1: Build the Go application
FROM golang:1.23-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY ../../go.mod ../../go.sum ./

# Download the dependencies
RUN go mod download

# Copy the entire employee directory content into the container
COPY . .

# Build the Go app
RUN go build -v -o /employee-service ./employee/cmd/main.go

# Step 2: Create a small image for the final executable
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the pre-built binary file from the build stage
COPY --from=build /employee-service .

# Expose port 8080 to the outside world
EXPOSE 8082

# Command to run the executable
CMD ["./employee-service"]
