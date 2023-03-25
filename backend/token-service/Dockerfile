# syntax=docker/dockerfile:1

FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY  . ./
RUN CGO_ENABLED=0 go build -o tokenservice

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build app/tokenservice /tokenservice

# GRPC port for communication with other services
ENV GRPC_PORT=9000
# Address to connect with redis
ENV REDIS_ADDRESS=redis:6379
# Redis password
ENV REDIS_PASSWORD=redis
# Secret used to sign refresh tokens
ENV REFRESH_SECRET=secret
# Refresh Token TTL
ENV REFRESH_DURATION=86400s
# Access Token TTL
ENV ACCESS_DURATION=1200s

EXPOSE 9000

USER nonroot:nonroot
ENTRYPOINT ["./tokenservice"]