# Dockerfile
FROM golang:1.23.0 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

ENV PORT=8443

EXPOSE 8443

CMD ["./main"]