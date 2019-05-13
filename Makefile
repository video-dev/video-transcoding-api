.PHONY: all testdeps lint test gotest build run

HTTP_PORT ?= 8080
LOG_LEVEL ?= debug
CI_TAG ?= $(shell git describe --tags $(shell git rev-list --tags --max-count=1))

all: test

testdeps:
	GO111MODULE=off go get github.com/golangci/golangci-lint/cmd/golangci-lint
	go mod download

lint: testdeps
	golangci-lint run \
		--enable-all \
		-D errcheck \
		-D lll \
		-D gochecknoglobals \
		-D goconst \
		-D gocyclo \
		-D dupl \
		-D gocritic \
		-D gochecknoinits \
		-D unparam \
		--deadline 5m ./...

gotest: testdeps
	go test ./...

test: lint gotest

coverage: lint
	go test -coverprofile=coverage.txt -covermode=atomic ./...

build:
	go build

run: build
	HTTP_PORT=$(HTTP_PORT) APP_LOG_LEVEL=$(LOG_LEVEL) ./video-transcoding-api
