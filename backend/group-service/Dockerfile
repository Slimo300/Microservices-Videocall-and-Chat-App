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
# Origin for CORS
ENV ORIGIN=http://localhost:3000
# Kafka Address
ENV BROKER_ADDRESS=
# S3 Bucket name for storing group profile pictures
ENV S3_BUCKET=



EXPOSE 8080

ENTRYPOINT ["./groupservice"]