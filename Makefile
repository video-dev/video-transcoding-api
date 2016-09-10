.PHONY: all testdeps lint test build run checkswagger swagger runswagger

HTTP_ACCESS_LOG ?= access.log
HTTP_PORT ?= 8080
LOG_LEVEL ?= debug

all: test

testdeps:
	go get github.com/go-swagger/go-swagger/cmd/swagger
	go get -d -t ./...

lint: testdeps
	go get github.com/alecthomas/gometalinter honnef.co/go/unused/cmd/unused
	gometalinter --install
	go get -t ./...
	gometalinter -j 4 --enable=gofmt --enable=unused --disable=dupl --disable=errcheck --disable=gas --disable=interfacer --disable=gocyclo --deadline=10m --tests ./...

test: lint checkswagger
	go test ./...

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
