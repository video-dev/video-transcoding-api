# This Dockerfile is intended to be used in the CI environment and depends on
# an existing binary built outside Docker.

FROM alpine:3.12.1 AS base

RUN apk add --no-cache upx
ADD video-transcoding-api /bin/video-transcoding-api
RUN upx --brute /bin/video-transcoding-api

FROM alpine:3.12.1

RUN apk add --no-cache ca-certificates
COPY --from=base /bin/video-transcoding-api /bin/video-transcoding-api
CMD ["/bin/video-transcoding-api"]
