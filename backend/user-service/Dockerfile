# syntax=docker/dockerfile:1

FROM golang:1.19-buster AS build

WORKDIR /app

COPY backend/user-service/go.mod ./
COPY backend/user-service/go.sum ./
RUN go mod download

COPY  backend/user-service/. ./
RUN CGO_ENABLED=0 go build -o userservice

FROM gcr.io/distroless/base-debian10

WORKDIR /

USER nonroot:nonroot

ARG configFile=/configs/docker-conf.yaml
ARG certDir=/cert
ARG templateDir=/backend/user-service/templates

COPY --from=build app/userservice /userservice
COPY ${configFile} ./config.yaml
COPY --chown=nonroot:nonroot ${certDir} ./cert
COPY --chown=nonroot:nonroot ${templateDir} ./templates

ENV CHAT_CONFIG=.

EXPOSE 8080
EXPOSE 8090

ENTRYPOINT ["./userservice"]