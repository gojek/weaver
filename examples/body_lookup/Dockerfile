FROM golang:1.11.5-alpine as base

ENV GO111MODULE off

RUN mkdir /estimate
ADD . /estimate
WORKDIR /estimate

RUN go build main.go

FROM alpine:latest
COPY --from=base /estimate/main /usr/local/bin/estimator
ENTRYPOINT ["estimator"]

