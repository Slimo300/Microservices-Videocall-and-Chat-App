# syntax=docker/dockerfile:1

FROM golang:1.19-buster AS build

WORKDIR /app

COPY backend/search-service/go.mod ./
COPY backend/search-service/go.sum ./
RUN go mod download

COPY  backend/search-service/. ./
RUN CGO_ENABLED=0 go build -o searchservice

FROM gcr.io/distroless/base-debian10

WORKDIR /

USER nonroot:nonroot

ARG configFile=/configs/docker-conf.yaml
ARG certDir=/cert

COPY --from=build app/searchservice /searchservice
COPY ${configFile} ./config.yaml
COPY --chown=nonroot:nonroot ${certDir} ./cert

ENV CHAT_CONFIG=.

EXPOSE 8080
EXPOSE 8090

ENTRYPOINT ["./searchservice"]