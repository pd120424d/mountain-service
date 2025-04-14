FROM golang:1.24 as builder

ARG VERSION
ARG GIT_SHA

WORKDIR /app
COPY . .

RUN go build -ldflags "-X main.Version=${VERSION} -X main.GitSHA=${GIT_SHA}" -o version-service

FROM debian:bookworm-slim

COPY --from=builder /app/version-service /version-service

EXPOSE 8090

CMD ["/version-service"]
