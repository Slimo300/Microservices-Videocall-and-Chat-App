# syntax=docker/dockerfile:1

FROM golang:1.23.1-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY  . ./
RUN CGO_ENABLED=0 go build -o messageservice

FROM gcr.io/distroless/base-debian10

WORKDIR /

USER nonroot:nonroot

COPY --from=build app/messageservice /messageservice

# Database address for storing user information
ENV MYSQL_ADDRESS=
# Port for HTTP traffic
ENV HTTP_PORT=8080
# Origin for CORS
ENV ORIGIN=http://localhost:3000
# Kafka Address
ENV BROKER_ADDRESS=
# Storage config
ENV STORAGE_BUCKET=
ENV STORAGE_KEY_ID=
ENV STORAGE_KEY_SECRET=
ENV STORAGE_USE_DO=
ENV STORAGE_URL=

EXPOSE 8080

ENTRYPOINT ["./messageservice"]