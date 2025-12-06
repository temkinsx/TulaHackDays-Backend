# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o health-map ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates curl postgresql-client bash && \
    addgroup -S app && adduser -S app -G app && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar -xz -C /usr/local/bin migrate

WORKDIR /app
COPY --from=builder /src/health-map /app/health-map
COPY --from=builder /src/migrations /app/migrations
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV POSTGRES_DSN="" \
    JWT_SECRET="" \
    SERVER_HOST=0.0.0.0 \
    SERVER_PORT=8080

EXPOSE 8080
USER app
ENTRYPOINT ["/entrypoint.sh"]
