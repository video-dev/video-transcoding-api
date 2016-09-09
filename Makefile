.PHONY: all testdeps checkfmt lint test build run vet checkswagger swagger runswagger

HTTP_ACCESS_LOG ?= access.log
HTTP_PORT ?= 8080
LOG_LEVEL ?= debug

all: test

testdeps:
	go get github.com/go-swagger/go-swagger/cmd/swagger
	go get -d -t ./...

checkfmt: testdeps
	[ -z "$$(gofmt -s -d . | tee /dev/stderr)" ]

deadcode:
	go get github.com/remyoudompheng/go-misc/deadcode
	go list ./... | sed -e "s;github.com/nytm/video-transcoding-api;.;" | xargs deadcode

lint: testdeps
	go get github.com/golang/lint/golint
	@for file in $$(git ls-files '*.go'); do \
		export output="$$(golint $${file})"; \
		[ -n "$${output}" ] && echo "$${output}" && export status=1; \
	done; \
	exit $${status:-0}

test: checkfmt lint vet deadcode checkswagger
	go test ./...

build:
	go build

run: build
	HTTP_PORT=$(HTTP_PORT) HTTP_ACCESS_LOG=$(HTTP_ACCESS_LOG) LOG_LEVEL=$(LOG_LEVEL) ./video-transcoding-api

vet: testdeps
	go vet ./...

swagger:
	go get github.com/go-swagger/go-swagger/cmd/swagger
	swagger generate spec -o swagger.json

checkswagger:
	swagger validate swagger.json

runswagger:
	go run swagger-ui-server/main.go
