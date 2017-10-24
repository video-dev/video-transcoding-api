.PHONY: all lint test gotest build run checkswagger swagger runswagger

HTTP_ACCESS_LOG ?= access.log
HTTP_PORT ?= 8080
LOG_LEVEL ?= debug
CI_TAG ?= $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
TAG_SUFFIX  := $(shell echo $(CI_TAG) | tail -c 3)

all: test

testdeps:
	go get github.com/go-swagger/go-swagger/cmd/swagger

lint:
	go get github.com/alecthomas/gometalinter
	gometalinter --install --vendored-linters
	go build -i
	go list -f '{{.TestImports}}' ./... | sed -e 's/\[\(.*\)\]/\1/' | tr ' ' '\n' | grep '^.*\..*/.*$$' | xargs go install
	gometalinter -j 2 --enable=misspell --enable=gofmt --enable=unused --disable=dupl --disable=errcheck --disable=gas --disable=interfacer --disable=gocyclo --deadline=10m --tests --vendor ./...

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

live:
	git clone --depth=1 https://$(GITHUB_TOKEN)@github.com/$(INFRA_REPO).git $(NP_PATH)
	go get github.com/$(NP_REPO)
	np deploy transcoding-api:stg#$(COMMIT)
	aws autoscaling put-notification-configuration --auto-scaling-group-name transcoding-api-stg-$(COMMIT)-web --topic-arn $(ASG_TOPIC) --notification-types "autoscaling:EC2_INSTANCE_LAUNCH" "autoscaling:EC2_INSTANCE_LAUNCH_ERROR"
ifneq ($(TAG_SUFFIX),rc)
	np deploy transcoding-api:prd#$(COMMIT)
	aws autoscaling put-notification-configuration --auto-scaling-group-name transcoding-api-prd-$(COMMIT)-web --topic-arn $(ASG_TOPIC) --notification-types "autoscaling:EC2_INSTANCE_LAUNCH" "autoscaling:EC2_INSTANCE_LAUNCH_ERROR"
endif
