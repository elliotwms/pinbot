# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

ARG APP_VERSION="v0.0.0+unknown"
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X github.com/elliotwms/pinbot/internal/build.Version=${APP_VERSION}" -o /pinbot ./cmd/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /pinbot /pinbot

EXPOSE 8080

ENTRYPOINT ["/pinbot"]