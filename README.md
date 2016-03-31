
![video-transcoding-api logo](https://cloud.githubusercontent.com/assets/244265/14191217/ae825932-f764-11e5-8eb3-d070aa8f2676.png)

# Video Transcoding API

The Video Transcoding API provides an agnostic API to transcode media assets across different cloud services. Currently, it supports the following providers:

- [Elemental Conductor](http://www.elementaltechnologies.com/products/elemental-conductor)
- [Encoding.com](encoding.com)
- [Amazon Elastic Transcoder](https://aws.amazon.com/elastictranscoder/)

## Setting Up

With [Go](https://golang.org/dl/) installed, make sure to export the follow environment variables:

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

#### For [Encoding.com](encoding.com)

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
```

In order to store preset maps, job statuses, and callback URLs we need a Redis instance running. Learn how to setup and run a Redis [here](http://redis.io/topics/quickstart). With the Redis instance running, set its configuration variables:

````
export REDIS_ADDR=200.221.14.140
export REDIS_PASSWORD=p4ssw0rd.here
```
If you are running Redis in the same host of the API and on the default port (6379) the API will automatically find the instance and connect to it. 

With all environment variables set and redis up and running, clone the repository and run:

```
$ git clone https://github.com/nytm/video-transcoding-api.git
$ make run
```

## Running tests

```
$ make test
```

## Using the API

### Creating a Preset

Given a JSON file called `preset.json`:

```json
{
  "providers": ["elastictranscoder", "elementalconductor", "encodingcom"],
  "preset": {
    "name": "Preset_Test",
    "description": "This is an example preset",
    "container": "m3u8",
    "height": "720",
    "videoCodec": "h264",
    "videoBitrate": "1000000",
    "gopSize": "90",
    "gopMode": "fixed",
    "profile": "Main",
    "profileLevel": "3.1",
    "rateControl": "VBR",
    "interlaceMode": "progressive",
    "audioCodec": "aac",
    "audioBitrate": "64000"
  }
}
```

The Encoding API will try to create the preset in all providers described on `providers` field. It will also create a PresetMap registry on Redis which is a map for all the PresetID on providers.

```
$ curl -X POST -d @preset.json http://api.host.com/presets
```
```json
{
  "PresetMap": "Preset_Test",
  "Results": {
    "elastictranscoder": {
      "Error": "",
      "PresetID": "1459293696042-8p8hak"
    },
    "elementalconductor": {
      "Error": "",
      "PresetID": "Preset_Test"
    },
    "encodingcom": {
      "Error": "creating preset: CreatePreset is not implemented in Encoding.com provider",
      "PresetID": ""
    }
  }
}
```

### Creating a PresetMap

Some providers like [Encoding.com](encoding.com) don't support preset creation throughout the API. The alternative is to create the preset manually on the control panel and associate it with the Encoding API by creating a PresetMap.

```
$ curl -XPOST -d '{"name":"preset_name", "providerMapping": {"encodingcom": "preset_id",
"output": {"extension": "ts"}}' http://api.host.com/presetmaps
```

### Listing PresetMaps

```
$ curl -X GET http://api.host.com/presetmaps
```
```json
{
  "nyt_preset": {
    "name": "nyt_preset",
    "output": {
      "extension": "ts"
    },
    "providerMapping": {
      "encodingcom": "preset_manual_id"
    }
  },
  "preset-1": {
    "name": "preset-1",
    "output": {
      "extension": "mp4"
    },
    "providerMapping": {
      "elastictranscoder": "1281742-93939",
      "elementalconductor": "abc123"
    }
  }
}
```

### Deleting PresetMap

```
$  curl -X DELETE http://api.host.com/presetmaps/preset-1
```

### Creating a Job

Given a `job.json`:
```json
{
  "presets": ["preset_1", "preset_2"],
  "provider": "encodingcom",
  "source": "http://nytimes.com/BickBuckBunny.mov?nocopy",
  "statusCallbackInterval": 5,
  "statusCallbackURL": "http://callback.server.com/status",
  "completionCallbackURL": "http://callback.server.com/done",
  "streamingParams": {}
}
```
```
$ curl -XPOST -H "Content-Type: application/json" -d @job.json  http://api.host.com/jobs
```

