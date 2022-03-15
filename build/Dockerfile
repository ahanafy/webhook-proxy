# syntax=docker/dockerfile:experimental

##
## Build
##

FROM golang:1.17 AS builder
WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./webhook-proxy ./cmd

##
## Deploy
##
FROM golang:1.17

WORKDIR /

COPY --from=builder /app/webhook-proxy /app

EXPOSE 8080

ENV INCLUSTER=true

ENTRYPOINT ["/app"]