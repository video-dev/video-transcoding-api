FROM golang:1.12.8-alpine AS build

ARG GOPROXY=https://proxy.golang.org

ENV  CGO_ENABLED 0
WORKDIR /code
ADD  . ./
RUN  go install

FROM alpine:3.10.1
RUN apk add --no-cache ca-certificates
COPY --from=build /go/bin/video-transcoding-api /usr/bin/video-transcoding-api
ENTRYPOINT ["/usr/bin/video-transcoding-api"]
