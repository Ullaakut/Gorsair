# Build stage
FROM golang:alpine AS build-env

COPY . /go/src/github.com/Ullaakut/gorsair
WORKDIR /go/src/github.com/Ullaakut/gorsair

ENV GO111MODULE=on
RUN go version
RUN go build -o bin/gorsair ./cmd

# Final stage
FROM alpine

COPY . /app/gorsair
WORKDIR /app/gorsair

ENTRYPOINT ["/app/gorsair/bin/gorsair"]
