FROM golang:1.11-alpine AS build

RUN     apk add --no-cache git
ENV     CGO_ENABLED 0
ADD     . /code
WORKDIR /code
RUN     go build -o /bin/video-transcoding-api

FROM alpine:3.8

RUN  apk add --no-cache ca-certificates
COPY --from=build /bin/video-transcoding-api /bin/video-transcoding-api
CMD ["/bin/video-transcoding-api"]
