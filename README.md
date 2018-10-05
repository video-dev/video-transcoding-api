![video-transcoding-api logo](https://cloud.githubusercontent.com/assets/244265/14191217/ae825932-f764-11e5-8eb3-d070aa8f2676.png)

# Video Transcoding API

[![Build Status](https://travis-ci.org/NYTimes/video-transcoding-api.svg?branch=master)](https://travis-ci.org/NYTimes/video-transcoding-api)
[![codecov](https://codecov.io/gh/NYTimes/video-transcoding-api/branch/master/graph/badge.svg)](https://codecov.io/gh/NYTimes/video-transcoding-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/NYTimes/video-transcoding-api)](https://goreportcard.com/report/github.com/NYTimes/video-transcoding-api)

The Video Transcoding API provides an agnostic API to transcode media assets
across different cloud services. Currently, it supports the following
providers:

- [Elemental Conductor](http://www.elementaltechnologies.com/products/elemental-conductor)
- [Encoding.com](http://encoding.com)
- [Amazon Elastic Transcoder](https://aws.amazon.com/elastictranscoder/)
- [Zencoder](http://zencoder.com)
- [Bitmovin](http://bitmovin.com)
- [Hybrik](https://www.hybrik.com)

## Setting Up

With [latest Go](https://golang.org/dl/) installed, make sure to export the follow
environment variables:

#### For [Elemental Conductor](http://www.elementaltechnologies.com/products/elemental-conductor)

```
export ELEMENTALCONDUCTOR_HOST=https://conductor-address.cloud.elementaltechnologies.com/
export ELEMENTALCONDUCTOR_USER_LOGIN=your.login
export ELEMENTALCONDUCTOR_API_KEY=your.api.key
export ELEMENTALCONDUCTOR_AUTH_EXPIRES=30
export ELEMENTALCONDUCTOR_AWS_ACCESS_KEY_ID=your.access.key.id
export ELEMENTALCONDUCTOR_AWS_SECRET_ACCESS_KEY=your.secret.access.key
export ELEMENTALCONDUCTOR_DESTINATION=s3://your-s3-bucket/
```

#### For [Encoding.com](http://encoding.com)

```
export ENCODINGCOM_USER_ID=your.user.id
export ENCODINGCOM_USER_KEY=your.user.key
export ENCODINGCOM_DESTINATION=http://access.key.id:secret.access.key@your-s3-bucket.s3.amazonaws.com/
export ENCODINGCOM_REGION="us-east-1"
```

#### For [Amazon Elastic Transcoder](https://aws.amazon.com/elastictranscoder/)

```
export AWS_ACCESS_KEY_ID=your.access.key.id
export AWS_SECRET_ACCESS_KEY=your.secret.access.key
export AWS_REGION="us-east-1"
export ELASTICTRANSCODER_PIPELINE_ID="yourpipeline-id"
```

#### For [Zencoder](http://zencoder.com)

```
export ZENCODER_API_KEY=your.api.key
export ZENCODER_DESTINATION=http://access.key.id:secret.access.key@your-s3-bucket.s3.amazonaws.com/
```

#### For [Bitmovin](http://bitmovin.com)

```
export BITMOVIN_API_KEY=your.api.key
export BITMOVIN_AWS_ACCESS_KEY_ID=your.access.key.id
export BITMOVIN_AWS_SECRET_ACCESS_KEY=your.secret.access.key
export BITMOVIN_AWS_STORAGE_REGION=your.s3.region.such.as.US_EAST_1.or.EU_WEST_1
export BITMOVIN_DESTINATION=s3://your-s3-bucket
export BITMOVIN_ENCODING_REGION=your.provider.region.such.as.AWS_US_EAST_1.or.GOOGLE_EUROPE_WEST_1
export BITMOVIN_ENCODING_VERSION=STABLE.or.BETA
```
#### For [Hybrik](https://www.hybrik.com)

```
export HYBRIK_URL=your.hybrik.api.endpoint.such.as.https://api_demo.hybrik.com/v1
export HYBRIK_COMPLIANCE_DATE=20170601
export HYBRIK_OAPI_KEY=your.hybrik.oapi.key
export HYBRIK_OAPI_SECRET=your.hybrik.oapi.secret
export HYBRIK_AUTH_KEY=your.hybrik.auth.key
export HYBRIK_AUTH_SECRET=your.hybrik.auth.secret
export HYBRIK_DESTINATION=s3://your-s3-bucket
export HYBRIK_PRESET_PATH=video-transcoding-api-presets
```

The HYBRIK_PRESET_PATH is optional and defines the folder presets will be stored in. If not specified, it will default to 'video-transcoding-api-presets'.

Please notice that for Elastic Transcoder you don't specify the destination
bucket, as it is [defined in the Elastic Transcoder
Pipeline](https://docs.aws.amazon.com/elastictranscoder/latest/developerguide/pipeline-settings.html#pipeline-settings-configure-transcoded-bucket).

In order to store preset maps and job statuses we need a Redis instance
running. Learn how to setup and run a Redis
[here](http://redis.io/topics/quickstart). With the Redis instance running, set
its configuration variables:

```
export REDIS_ADDR=192.0.2.31
export REDIS_PASSWORD=p4ssw0rd.here
```

If you are running Redis in the same host of the API and on the default port
(6379) the API will automatically find the instance and connect to it.

With all environment variables set and redis up and running, clone this
repository and run:

```
$ git clone https://github.com/NYTimes/video-transcoding-api.git
$ make run
```

## Running tests

```
$ make test
```

## Using the API

Check out on our Wiki [how
to](https://github.com/NYTimes/video-transcoding-api/wiki/Using-Video-Transcoding-API)
use this API.

## Contributing

1. Fork it
2. Create your feature branch: `git checkout -b my-awesome-new-feature`
3. Commit your changes: `git commit -m 'Add some awesome feature'`
4. Push to the branch: `git push origin my-awesome-new-feature`
5. Submit a pull request

## License

- This code is under [Apache 2.0
  license](https://github.com/NYTimes/video-transcoding-api/blob/master/LICENSE).
- The video-transcoding-api logo is a variation on the Go gopher that was
  designed by Renee French and copyrighted under the [Creative Commons
  Attribution 3.0 license](https://creativecommons.org/licenses/by/3.0/).
