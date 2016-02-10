.PHONY: all testdeps checkfmt lint test build run vet swagger

all: test

testdeps:
	go get -d -t ./...

checkfmt: testdeps
	@export output="$$(gofmt -s -l .)" && \
		[ -n "$${output}" ] && \
		echo "Unformatted files:" && \
		echo && echo "$${output}" && \
		echo && echo "Please fix them using 'gofmt -s -w .'" && \
		export status=1; exit $${status:-0}

deadcode:
	go get github.com/reillywatson/go-misc/deadcode
	go list ./... | sed -e "s;github.com/nytm/video-transcoding-api;.;" | xargs deadcode

lint: testdeps
	go get github.com/golang/lint/golint
	@for file in $$(find . -name '*.go' | grep -v '_test'); do \
		export output="$$(golint $${file})"; \
		[ -n "$${output}" ] && echo "$${output}" && export status=1; \
	done; \
	exit $${status:-0}

test: checkfmt lint vet deadcode
	go test ./...

build:
	go build

run: build
	./video-transcoding-api -config config.json

vet: testdeps
	go get golang.org/x/tools/cmd/vet
	go vet ./...

swagger:
	go get github.com/go-swagger/go-swagger/cmd/swagger
	swagger generate spec -o swagger.json
	go build -o generate-readme ./doc
	./generate-readme -i swagger.json -o README.md
