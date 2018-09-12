.PHONY: all testdeps lint test gotest build run checkswagger swagger runswagger

HTTP_PORT ?= 8080
LOG_LEVEL ?= debug
CI_TAG ?= $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
TAG_SUFFIX  := $(shell echo $(CI_TAG) | tail -c 3)

all: test

testdeps:
	go get github.com/go-swagger/go-swagger/cmd/swagger

lint: testdeps
	go get github.com/alecthomas/gometalinter
	gometalinter --install --vendored-linters
	go build -i
	go list -f '{{.TestImports}}' ./... | sed -e 's/\[\(.*\)\]/\1/' | tr ' ' '\n' | grep '^.*\..*/.*$$' | xargs go install
	gometalinter -j 2 --fast --disable=gosec --disable=gotype --disable=gas --disable=gocyclo --deadline=10m --tests --vendor ./...

gotest: testdeps
	go test ./...

test: lint checkswagger gotest

coverage: lint checkswagger
	@rm -f coverage.txt; for p in $$(go list ./...); do \
		go test -coverprofile=profile.out -covermode=atomic $$p || export status=2; \
		if [ -f profile.out ]; then cat profile.out >> coverage.txt; rm profile.out; fi \
		done; \
		exit $${status:-0}

build:
	go build

run: build
	HTTP_PORT=$(HTTP_PORT) APP_LOG_LEVEL=$(LOG_LEVEL) ./video-transcoding-api

swagger:
	go get github.com/go-swagger/go-swagger/cmd/swagger
	swagger generate spec -o swagger.json

checkswagger:
	swagger validate swagger.json
