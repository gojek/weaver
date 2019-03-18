.PHONY: all

all: build fmt vet lint test coverage
default: build fmt vet lint test

ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
APP_EXECUTABLE="out/weaver-server"
COMMIT_HASH=$(shell git rev-parse --verify head | cut -c-1-8)
BUILD_DATE=$(shell date +%Y-%m-%dT%H:%M:%S%z)

setup:
	GO111MODULE=on go get -u golang.org/x/lint/golint
	GO111MODULE=on go get github.com/mattn/goveralls

compile:
	mkdir -p out/
	GO111MODULE=on go build -o $(APP_EXECUTABLE) -ldflags "-X main.BuildDate=$(BUILD_DATE) -X main.Commit=$(COMMIT_HASH) -s -w" ./cmd/weaver-server

build: deps compile fmt vet lint

deps:
	GO111MODULE=on go mod tidy -v

install:
	GO111MODULE=on go install ./...

fmt:
	GO111MODULE=on go fmt ./...

vet:
	GO111MODULE=on go vet ./...

lint:
	@if [[ `golint $(All_PACKAGES) | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } | wc -l | tr -d ' '` -ne 0 ]]; then \
		golint $(ALL_PACKAGES) | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; }; \
		exit 2; \
	fi;

test: copy-config
	GO111MODULE=on go test ./...

test-cover-html:
	@echo "mode: count" > coverage-all.out
	$(foreach pkg, $(ALL_PACKAGES),\
	go test -coverprofile=coverage.out -covermode=count $(pkg);\
	tail -n +2 coverage.out >> coverage-all.out;)
	GO111MODULE=on go tool cover -html=coverage-all.out -o out/coverage.html

copy-config:
	cp weaver.conf.yaml.sample weaver.conf.yaml

clean:
	go clean && rm -rf ./vendor ./build ./weaver.conf.yaml

docker-clean:
	docker-compose down

docker-spec: docker-clean
	docker-compose build
	docker-compose run --entrypoint "make test" dev_weaver

docker-server:
	docker-compose run --entrypoint "make local-server" dev_weaver

docker-up:
	docker-compose up -d

local-server: compile
	$(APP_EXECUTABLE) start

coverage:
	goveralls -service=travis-ci

