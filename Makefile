.PHONY: all testdeps lint test gotest build run checkswagger swagger runswagger

HTTP_PORT ?= 8080
LOG_LEVEL ?= debug
CI_TAG ?= $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
TAG_SUFFIX  := $(shell echo $(CI_TAG) | tail -c 3)

all: test

testdeps:
	cd /tmp && go get github.com/go-swagger/go-swagger/cmd/swagger

lint: testdeps
	cd /tmp && go get golang.org/x/lint/golint
	golint $$(go list ./...)

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

swagger: testdeps
	swagger generate spec -o swagger.json

checkswagger: testdeps
	swagger validate swagger.json
