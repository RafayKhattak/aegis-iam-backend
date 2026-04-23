FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o aegis-api ./cmd/api

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -g '' khattak

WORKDIR /app

COPY --from=builder /app/aegis-api ./aegis-api
COPY .env.example .env

EXPOSE 8080

USER khattak

ENTRYPOINT ["./aegis-api"]
