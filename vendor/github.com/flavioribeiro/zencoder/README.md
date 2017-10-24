zencoder
========

[![Build Status](https://travis-ci.org/brandscreen/zencoder.png)](https://travis-ci.org/brandscreen/zencoder) [![Coverage Status](https://coveralls.io/repos/brandscreen/zencoder/badge.png?branch=HEAD)](https://coveralls.io/r/brandscreen/zencoder?branch=HEAD)

[Go](http://golang.org) integration for [Zencoder API](http://www.zencoder.com/) video transcoding service.

# Requirements

* Go 1.1 or higher
* A Zencoder account/API key (get one at app.zencoder.com)

# Documentation

[Godoc](http://godoc.org/) documentation is available at [http://godoc.org/github.com/brandscreen/zencoder](http://godoc.org/github.com/brandscreen/zencoder).

# Installation

```bash
$ go get github.com/brandscreen/zencoder
```

# Usage

## Import Zencoder

Ensure you have imported the zencoder package at the top of your source file.

```golang
import "github.com/brandscreen/zencoder"
```

## Create a connection to Zencoder

All Zencoder methods are on the Zencoder struct.  Create a new one bound to your API key using ```zencoder.NewZencoder```.

```golang
// make sure you replace [YOUR API KEY HERE] with your API key
zc := zencoder.NewZencoder("[YOUR API KEY HERE]")
```

## [Jobs](https://app.zencoder.com/docs/api/jobs)

### [Create a Job](https://app.zencoder.com/docs/api/jobs/create)
```golang
settings := &zencoder.EncodingSettings{
    Input: "s3://zencodertesting/test.mov",
    Test:  true,
}
job, err := zc.CreateJob(settings)
```

### [List Jobs](https://app.zencoder.com/docs/api/jobs/list)
```golang
jobs, err := zc.ListJobs()
```

### [Get Job Details](https://app.zencoder.com/docs/api/jobs/show)
```golang
details, err := zc.GetJobDetails(12345)
```

### [Job Progress](https://app.zencoder.com/docs/api/jobs/progress)
```golang
progress, err := zc.GetJobProgress(12345)
```

### [Resubmit a Job](https://app.zencoder.com/docs/api/jobs/resubmit)
```golang
err := zc.ResubmitJob(12345)
```

### [Cancel a Job](https://app.zencoder.com/docs/api/jobs/cancel)
```golang
err := zc.CancelJob(12345)
```

### [Finish a Live Job](https://app.zencoder.com/docs/api/jobs/finish)
```golang
err := zc.FinishLiveJob(12345)
```

## [Inputs](https://app.zencoder.com/docs/api/inputs)

### [Get Input Details](https://app.zencoder.com/docs/api/inputs/show)
```golang
details, err := zc.GetInputDetails(12345)
```

### [Input Progress](https://app.zencoder.com/docs/api/inputs/progress)
```golang
progress, err := zc.GetInputProgress(12345)
```

## [Outputs](https://app.zencoder.com/docs/api/outputs)

### [Get Output Details](https://app.zencoder.com/docs/api/outputs/show)
```golang
details, err := zc.GetOutputDetails(12345)
```

### [Output Progress](https://app.zencoder.com/docs/api/outputs/progress)
```golang
progress, err := zc.GetOutputProgress(12345)
```

## [Accounts](https://app.zencoder.com/docs/api/accounts)

### [Get Account Details](https://app.zencoder.com/docs/api/accounts/show)
```golang
account, err := zc.GetAccount()
```

### [Set Integration Mode](https://app.zencoder.com/docs/api/accounts/integration)
```golang
err := zc.SetIntegrationMode()
```

### [Set Live Mode](https://app.zencoder.com/docs/api/accounts/integration)
```golang
err := zc.SetLiveMode()
```

## [Reports](https://app.zencoder.com/docs/api/reports)

### ReportSettings

All reporting interfaces take either ```nil``` or a ```ReportSettings``` object.

Using ```nil``` denotes to use default settings.  In this case, assume ```settings``` in the examples below is defined as:

```golang
var settings *zencoder.ReportSettings = nil
```

A ```ReportSettings``` object can either be constructed manually as in:
```golang
var start, end time.Date
settings := &zencoder.ReportSettings{
    From:     &start,
    To:       &end,
    Grouping: "key",
}
```

Or, you can use a Fluent-style interface to build a ```ReportSettings``` object, as in:

```golang
var start, end time.Date
settings := zencoder.ReportFrom(start).To(time).Grouping("key")
```

### [Get VOD Usage](https://app.zencoder.com/docs/api/reports/vod)

```golang
usage, err := zc.GetVodUsage(settings)
```

### [Get Live Usage](https://app.zencoder.com/docs/api/reports/live)

```golang
usage, err := zc.GetLiveUsage(settings)
```

### [Get Total Usage](https://app.zencoder.com/docs/api/reports/all)

```golang
usage, err := zc.GetUsage(settings)
```

## Encoding Settings

See [Zencoder API documentation](https://app.zencoder.com/docs/api/encoding) for all encoding settings available in zencoder.EncodingSettings.  All settings are currently supported, with the main difference being the casing of the options to fit with Go naming conventions.

# Contributing

Please see [CONTRIBUTING.md](https://github.com/brandscreen/zencoder/blob/master/CONTRIBUTING.md).  If you have a bugfix or new feature that you would like to contribute, please find or open an issue about it first.

# License

Licensed under the MIT License.
