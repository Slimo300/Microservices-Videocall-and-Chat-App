# syntax=docker/dockerfile:1

FROM golang:1.22-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY  . ./
RUN CGO_ENABLED=0 go build -o userservice

FROM gcr.io/distroless/base-debian10

WORKDIR /

USER nonroot:nonroot

COPY --from=build app/userservice /userservice

# Database address for storing user information
ENV MYSQL_ADDRESS=
# Port for HTTP traffic
ENV HTTP_PORT=8080
# Address to connect with token service
ENV TOKEN_SERVICE_ADDRESS=
# Address to connect with email service
ENV EMAIL_SERVICE_ADDRESS=
# Origin for CORS
ENV ORIGIN=http://localhost:3000
# Domain for cookies
ENV DOMAIN=localhost
# Kafka Address
ENV BROKER_ADDRESS=
# Storage config
ENV STORAGE_BUCKET=
ENV STORAGE_KEY_ID=
ENV STORAGE_KEY_SECRET=
ENV STORAGE_USE_DO=


EXPOSE 8080

ENTRYPOINT ["./userservice"]