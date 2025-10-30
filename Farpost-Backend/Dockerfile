FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/vlru-prsch/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server .

COPY ./storage ./storage

COPY ./config/local.yaml ./config.yaml

ENV CONFIG_PATH=/app/config.yaml

EXPOSE 12345

CMD ["./server"]
