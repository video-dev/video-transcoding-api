FROM golang:1.14.1-alpine AS build

ARG GOPROXY=https://proxy.golang.org,https://gocenter.io,direct

ENV  CGO_ENABLED 0
WORKDIR /code
ADD  . ./
RUN  go install

FROM alpine:3.11.2
RUN apk add --no-cache ca-certificates
COPY --from=build /go/bin/video-transcoding-api /usr/bin/video-transcoding-api
ENTRYPOINT ["/usr/bin/video-transcoding-api"]
