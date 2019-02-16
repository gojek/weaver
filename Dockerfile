FROM golang:1.11.5-alpine as base

ENV GO111MODULE on

RUN apk --no-cache add gcc g++ make ca-certificates git
RUN mkdir /weaver
ADD . /weaver
WORKDIR /weaver

RUN make setup
RUN make build

FROM golang:1.11.2-alpine

COPY --from=base /weaver/out/weaver-server /usr/local/bin/weaver

ENTRYPOINT ["weaver", "server"]
