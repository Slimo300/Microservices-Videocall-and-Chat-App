# syntax=docker/dockerfile:1

FROM golang:1.19-buster AS build

WORKDIR /app

COPY backend/message-service/go.mod ./
COPY backend/message-service/go.sum ./
RUN go mod download

COPY  backend/message-service/. ./
RUN CGO_ENABLED=0 go build -o messageservice

FROM gcr.io/distroless/base-debian10

WORKDIR /

USER nonroot:nonroot

ARG configFile=/config/withPostgres.yaml
ARG certDir=/cert

COPY --from=build app/messageservice /messageservice
COPY ${configFile} ./config.yaml
COPY --chown=nonroot:nonroot ${certDir} ./cert

ENV CHAT_CONFIG=.

EXPOSE 8080
EXPOSE 8090

ENTRYPOINT ["./messageservice"]