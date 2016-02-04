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

lint: testdeps
	go get github.com/golang/lint/golint
	@for file in $$(find . -name '*.go' | grep -v '_test'); do \
		export output="$$(golint $${file})"; \
		[ -n "$${output}" ] && echo "$${output}" && export status=1; \
	done; \
	exit $${status:-0}

test: checkfmt lint vet
	go test ./...

build:
	go build

run: build
	./video-transcoding-api -config config.json

vet: testdeps
	go get golang.org/x/tools/cmd/vet
	go vet ./...
