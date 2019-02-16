FROM golang:1.11.5-alpine as builder

ENV GO111MODULE on

RUN apk --no-cache add gcc g++ make ca-certificates git

RUN mkdir /weaver
WORKDIR /weaver

COPY . /weaver

RUN make setup
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /weaver/out/weaver-server .

ENTRYPOINT ["/weaver-server", "start"]
