FROM golang:1.25.1 AS builder

ARG VERSION
ARG GIT_SHA
ARG GIT_TAG
ARG BUILD_TIMESTAMP
ARG GIT_COMMIT_SHA

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy version-service, employee, and shared code
COPY version-service/ ./version-service/
COPY shared/ ./shared/
COPY employee/ ./employee/

RUN GIT_SHA_FINAL="${GIT_SHA:-$GIT_COMMIT_SHA}"; \
  GIT_FULL_SHA_FINAL="${GIT_COMMIT_SHA}"; \
  go build -ldflags "-X main.Version=${VERSION} -X main.GitSHA=${GIT_SHA_FINAL} -X main.FullGitSHA=${GIT_FULL_SHA_FINAL} -X main.GitTag=${GIT_TAG}" \
  -v -o /version-service ./version-service/main.go

FROM debian:bookworm-slim

# Install wget for health checks
RUN apt-get update && apt-get install -y wget && rm -rf /var/lib/apt/lists/*

COPY --from=builder /version-service /version-service

EXPOSE 8090

CMD ["/version-service"]
