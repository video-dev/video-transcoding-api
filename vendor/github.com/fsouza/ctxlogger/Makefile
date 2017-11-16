# Copyright 2016 Francisco Souza. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

all: test

test: lint gotest

coverage: testdeps
	go test -race -coverprofile coverage.txt -covermode=atomic

lint: testdeps
	go get github.com/alecthomas/gometalinter
	gometalinter --install --vendored-linters
	go install
	go list -f '{{.TestImports}}' | sed -e 's/\[\(.*\)\]/\1/' | tr ' ' '\n' | grep '^.*\..*/.*$$' | xargs go install
	gometalinter -j 4 --enable-all --line-length=120 --deadline=10m --tests

gotest: testdeps
	go test -race

testdeps:
	go get -d -t
