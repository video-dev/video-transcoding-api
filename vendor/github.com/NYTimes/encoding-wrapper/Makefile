all: test

testdeps:
	go get -d -t ./...

checkfmt: testdeps
	[ -z "$$(gofmt -s -d . | tee /dev/stderr)" ]

lint: testdeps
	go get github.com/alecthomas/gometalinter honnef.co/go/unused/cmd/unused
	gometalinter --install --vendored-linters
	go install ./...
	go list -f '{{.TestImports}}' ./... | sed -e 's/\[\(.*\)\]/\1/' | tr ' ' '\n' | grep '^.*\..*/.*$$' | xargs go install
	gometalinter -j 4 --enable=misspell --enable=gofmt --enable=unused --disable=dupl --disable=errcheck --disable=gas --deadline=10m --tests ./...

coverage: lint
	@rm -f coverage.txt; for p in $$(go list ./...); do \
		go test -coverprofile=profile.out -covermode=atomic $$p || export status=2; \
		if [ -f profile.out ]; then cat profile.out >> coverage.txt; rm profile.out; fi \
		done; \
		exit $${status:-0}

test: lint
	go test ./...
