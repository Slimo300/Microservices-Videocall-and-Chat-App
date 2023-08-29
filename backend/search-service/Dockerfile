# syntax=docker/dockerfile:1

FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY  . ./
RUN CGO_ENABLED=0 go build -o searchservice

FROM gcr.io/distroless/base-debian10

WORKDIR /

USER nonroot:nonroot

COPY --from=build app/searchservice /searchservice

# Database address for storing user information
ENV ES_ADDRESS=
# Database user for authentication
ENV ES_USER=
# Database user password for authentication
ENV ES_PASSWORD=
# Port for HTTP traffic
ENV HTTP_PORT=8080
# Origin for CORS
ENV ORIGIN=http://localhost:3000
# Kafka Address
ENV BROKER_ADDRESS=

EXPOSE 8080

ENTRYPOINT ["./searchservice"]