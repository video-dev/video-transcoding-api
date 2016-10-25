.PHONY: all testdeps lint test gotest build run checkswagger swagger runswagger

HTTP_ACCESS_LOG ?= access.log
HTTP_PORT ?= 8080
LOG_LEVEL ?= debug

all: test

testdeps:
	go get github.com/go-swagger/go-swagger/cmd/swagger
	go get -d -t ./...

lint: testdeps
	go get github.com/alecthomas/gometalinter honnef.co/go/unused/cmd/unused
	gometalinter --install --vendored-linters
	go get -t ./...
	go list -f '{{.TestImports}}' ./... | sed -e 's/\[\(.*\)\]/\1/' | tr ' ' '\n' | grep '^.*\..*/.*$$' | xargs go install
	gometalinter -j 2 --enable=misspell --enable=gofmt --enable=unused --disable=dupl --disable=errcheck --disable=gas --disable=interfacer --disable=gocyclo --deadline=10m --tests ./...

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
	HTTP_PORT=$(HTTP_PORT) HTTP_ACCESS_LOG=$(HTTP_ACCESS_LOG) APP_LOG_LEVEL=$(LOG_LEVEL) ./video-transcoding-api

swagger:
	go get github.com/go-swagger/go-swagger/cmd/swagger
	swagger generate spec -o swagger.json

checkswagger:
	swagger validate swagger.json

runswagger:
	go run swagger-ui-server/main.go
