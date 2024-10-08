# syntax=docker/dockerfile:1

FROM golang:1.23.1-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY  . ./
RUN CGO_ENABLED=0 go build -o emailservice

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build app/emailservice /emailservice

# GRPC port for communication with other services
ENV GRPC_PORT=9000
# API address for email links
ENV ORIGIN=http://localhost:3000
# What email address will be shown as sender
ENV EMAIL_FROM=MicroChat@mail.com
# SMTP Host on which provider accepts connections
ENV SMTP_HOST=
# SMTP Port on which provider accepts connections
ENV SMTP_PORT=
# SMTP User
ENV SMTP_USER=
# SMTP Password
ENV SMTP_PASS=

EXPOSE 9000

USER nonroot:nonroot
ENTRYPOINT ["./emailservice"]