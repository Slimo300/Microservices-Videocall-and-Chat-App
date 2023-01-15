# syntax=docker/dockerfile:1

FROM golang:1.19-buster AS build

WORKDIR /app

COPY backend/ws-service/go.mod ./
COPY backend/ws-service/go.sum ./
RUN go mod download

COPY  backend/ws-service/. ./
RUN CGO_ENABLED=0 go build -o wsservice

FROM gcr.io/distroless/base-debian10

WORKDIR /

USER nonroot:nonroot

ARG configFile=/config/docker-conf.yaml
ARG certDir=/cert

COPY --from=build app/wsservice /wsservice
COPY ${configFile} ./config.yaml
COPY --chown=nonroot:nonroot ${certDir} ./cert

ENV CHAT_CONFIG=.

EXPOSE 8080
EXPOSE 8090

ENTRYPOINT ["./wsservice"]