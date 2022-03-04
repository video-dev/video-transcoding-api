FROM golang:1.18rc1-alpine AS build

ENV  CGO_ENABLED 0
WORKDIR /code
ADD  . ./
RUN  go install

FROM alpine:3.15.0
RUN apk add --no-cache ca-certificates
COPY --from=build /go/bin/video-transcoding-api /usr/bin/video-transcoding-api
ENTRYPOINT ["/usr/bin/video-transcoding-api"]
