FROM golang:alpine

RUN apk add --no-cache ca-certificates build-base git
WORKDIR /go/src/github.com/NYTimes/video-transcoding-api/
COPY . / /go/src/github.com/NYTimes/video-transcoding-api/
RUN make build

FROM alpine:latest  

WORKDIR /root/
COPY --from=0 /go/src/github.com/NYTimes/video-transcoding-api/video-transcoding-api .
CMD ["./video-transcoding-api"]  
