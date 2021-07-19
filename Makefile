.PHONY: all testdeps lint runlint test gotest coverage gocoverage build run

HTTP_PORT ?= 8080
LOG_LEVEL ?= debug

all: test

testdeps:
	cd /tmp && GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go mod download

lint: testdeps runlint

runlint:
	golangci-lint run

gotest:
	go test -race -vet=all -mod=readonly $(GO_TEST_EXTRA_FLAGS) ./...

test: lint testdeps gotest

coverage: lint gocoverage

gocoverage:
	make gotest GO_TEST_EXTRA_FLAGS="-coverprofile=coverage.txt -covermode=atomic"

build:
	go build -mod=readonly -o video-transcoding-api

run: build
	HTTP_PORT=$(HTTP_PORT) APP_LOG_LEVEL=$(LOG_LEVEL) ./video-transcoding-api
