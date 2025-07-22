FROM golang:1.24 AS builder

ARG VERSION
ARG GIT_SHA

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy version-service, employee, and shared code
COPY version-service/ ./version-service/
COPY shared/ ./shared/
COPY employee/ ./employee/

RUN go build -ldflags "-X main.Version=${VERSION} -X main.GitSHA=${GIT_SHA}" \
  -v -o /version-service ./version-service/main.go

FROM debian:bookworm-slim

# Install wget for health checks
RUN apt-get update && apt-get install -y wget && rm -rf /var/lib/apt/lists/*

COPY --from=builder /version-service /version-service

EXPOSE 8090

CMD ["/version-service"]
