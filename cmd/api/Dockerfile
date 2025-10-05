FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/bin/main ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/main /app/bin/main

CMD ["./bin/main"]
