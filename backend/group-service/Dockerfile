# syntax=docker/dockerfile:1

FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY  . ./
RUN CGO_ENABLED=0 go build -o groupservice

FROM gcr.io/distroless/base-debian10

WORKDIR /

USER nonroot:nonroot

COPY --from=build app/groupservice /groupservice

# Database address for storing user information
ENV MYSQL_ADDRESS=
# Port for HTTP traffic
ENV HTTP_PORT=8080
# Port for HTTPS traffic
ENV HTTPS_PORT=8090
# Address to connect with token service
ENV TOKEN_SERVICE_ADDRESS=
# Origin for CORS
ENV ORIGIN=http://localhost:3000
# Kafka Address
ENV BROKER_ADDRESS=
# Directory on docker container in which SSL certificate and private key should be
ENV CERT_DIR=/cert
# S3 Bucket name for storing group profile pictures
ENV S3_BUCKET=



EXPOSE 8080
EXPOSE 8090

ENTRYPOINT ["./groupservice"]