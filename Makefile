.PHONY: all testdeps lint test gotest build run checkswagger swagger runswagger

HTTP_PORT ?= 8080
LOG_LEVEL ?= debug
CI_TAG ?= $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
TAG_SUFFIX  := $(shell echo $(CI_TAG) | tail -c 3)

all: test

testdeps:
	cd /tmp && go get github.com/go-swagger/go-swagger/cmd/swagger

lint: testdeps
	cd /tmp && go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --fast -D errcheck -E megacheck --deadline 5m ./...

gotest: testdeps
	go test ./...

test: lint checkswagger gotest

coverage: lint checkswagger
	go test -coverprofile=coverage.txt -covermode=atomic ./...

build:
	go build

run: build
	HTTP_PORT=$(HTTP_PORT) APP_LOG_LEVEL=$(LOG_LEVEL) ./video-transcoding-api

swagger: testdeps
	swagger generate spec -o swagger.json

checkswagger: testdeps
	swagger validate swagger.json
