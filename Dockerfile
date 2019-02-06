FROM golang:1.11.2-alpine

ENV GOOS linux
ENV GOARCH amd64
ENV GO111MODULE on

RUN apk add bash ca-certificates git make

RUN mkdir /weaver
WORKDIR /weaver

COPY go.mod .
COPY go.sum .

RUN go mod download
