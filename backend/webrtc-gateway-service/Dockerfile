# syntax=docker/dockerfile:1

FROM golang:1.23.1-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY  . ./
RUN CGO_ENABLED=0 go build -o webrtc-gateway-service

FROM gcr.io/distroless/base-debian10

WORKDIR /

USER nonroot:nonroot

COPY --from=build app/webrtc-gateway-service /webrtc-gateway-service

# Database address for storing user information
ENV DB_ADDRESS=
# Database password
ENV DB_PASSWORD=
# Port for HTTP traffic
ENV HTTP_PORT=8080
# Origin for CORS
ENV ORIGIN=http://localhost:3000
# Kafka Address
ENV BROKER_ADDRESS=

EXPOSE 8080

ENTRYPOINT ["./webrtc-gateway-service"]