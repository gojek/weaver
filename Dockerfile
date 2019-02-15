FROM golang:1.11.5-alpine

ENV GO111MODULE on

RUN apk --no-cache add gcc g++ make ca-certificates git

RUN mkdir /weaver
WORKDIR /weaver
