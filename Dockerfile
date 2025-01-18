FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main

FROM alpine:latest

WORKDIR /app

COPY resources/ .

COPY --from=builder /app/main /app/main


EXPOSE 443
CMD ["./main"]
