FROM golang:1.11.2-alpine

ENV GO111MODULE on

RUN apk --no-cache add gcc g++ make ca-certificates git

RUN mkdir /weaver
WORKDIR /weaver

COPY go.mod .
COPY go.sum .
