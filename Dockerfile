ARG GO_VERSION=1.23

FROM golang:${GO_VERSION}-alpine AS builder

RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY repositories repositories
COPY database database
COPY models models
COPY server server

COPY studentpb studentpb
COPY server-student server-student

COPY testpb testpb
COPY server-test server-test


RUN go install ./...

FROM alpine:latest

WORKDIR /usr/bin

COPY --from=builder /go/bin .