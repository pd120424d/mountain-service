# /service1/Dockerfile
FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o main ./cmd/employee

EXPOSE 8082

CMD ["./main"]
