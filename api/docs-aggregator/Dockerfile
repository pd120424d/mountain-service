# Build stage
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/docs-aggregator
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/docs-aggregator ./main.go

# Runtime stage
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /
COPY --from=build /out/docs-aggregator /docs-aggregator
USER nonroot:nonroot
EXPOSE 8080
ENV EXTERNAL_SCHEME=https
ENV EXTERNAL_HOST=mountain-service.duckdns.org
ENV EXTERNAL_BASE_PATH=/api/v1
ENTRYPOINT ["/docs-aggregator"]

