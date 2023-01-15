# syntax=docker/dockerfile:1

FROM golang:1.17-buster AS build

WORKDIR /app

COPY backend/token-service/go.mod ./
COPY backend/token-service/go.sum ./
RUN go mod download

COPY  backend/token-service/. ./
RUN CGO_ENABLED=0 go build -o tokenservice

FROM gcr.io/distroless/base-debian10

WORKDIR /

ARG configFile=/config/docker-conf.yaml
ARG rsaKey=private.pem

COPY --from=build app/tokenservice /tokenservice
COPY ${configFile} ./config.yaml
COPY ${rsaKey} ./private.pem

ENV CHAT_CONFIG=.

EXPOSE 9000

USER nonroot:nonroot
ENTRYPOINT ["./tokenservice"]