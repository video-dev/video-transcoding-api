package encodingcom

import (
	"time"

	"gopkg.in/check.v1"
)

func (s *S) TestGetStatusSingle(c *check.C) {
	server, requests := s.startServer(`
{
	"response": {
        "job": {
				"id": "abc123",
				"userid": "myuser",
				"sourcefile": "http://some.video/file.mp4",
				"status": "Finished",
				"notifyurl": "http://ping.me/please",
				"created": "2015-12-31 20:45:30",
				"started": "2015-12-31 20:45:34",
				"finished": "2015-12-31 21:00:03",
				"prevstatus": "Saving",
				"downloaded": "2015-12-31 20:45:32",
				"uploaded": "2015-12-31 20:59:54",
				"time_left": "0",
				"progress": "100",
				"time_left_current": "0",
				"progress_current": "100.0",
				"format": {
						"id": "f123",
						"status": "Finished",
						"description": "Something",
						"created": "2015-12-31 20:45:30",
						"started": "2015-12-31 20:45:34",
						"finished": "2015-12-31 21:00:03",
						"s3_destination": "https://s3.amazonaws.com/not-really/valid.mp4",
						"cf_destination": "https://blablabla.cloudfront.net/not-valid.mp4",
						"convertedsize": "65723",
						"destination": "s3://mynicebucket",
						"destination_status": "Saved"
					}
			}
	}
}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	status, err := client.GetStatus([]string{"abc123"}, true)
	c.Assert(err, check.IsNil)

	expectedCreateDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:30")
	expectedStartDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:34")
	expectedFinishDate, _ := time.Parse(dateTimeLayout, "2015-12-31 21:00:03")
	expectedDownloadDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:32")
	expectedUploadDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:59:54")
	expected := []StatusResponse{
		{
			MediaID:             "abc123",
			UserID:              "myuser",
			SourceFile:          "http://some.video/file.mp4",
			MediaStatus:         "Finished",
			PreviousMediaStatus: "Saving",
			NotifyURL:           "http://ping.me/please",
			CreateDate:          expectedCreateDate,
			StartDate:           expectedStartDate,
			FinishDate:          expectedFinishDate,
			DownloadDate:        expectedDownloadDate,
			UploadDate:          expectedUploadDate,
			TimeLeft:            "0",
			Progress:            100.0,
			TimeLeftCurrentJob:  "0",
			ProgressCurrentJob:  100.0,
			Formats: []FormatStatus{
				{
					ID:            "f123",
					Status:        "Finished",
					Description:   "Something",
					CreateDate:    expectedCreateDate,
					StartDate:     expectedStartDate,
					FinishDate:    expectedFinishDate,
					S3Destination: "https://s3.amazonaws.com/not-really/valid.mp4",
					CFDestination: "https://blablabla.cloudfront.net/not-valid.mp4",
					FileSize:      "65723",
					Destinations:  []DestinationStatus{{Name: "s3://mynicebucket", Status: "Saved"}},
				},
			},
		},
	}
	c.Assert(status, check.DeepEquals, expected)

	req := <-requests
	c.Assert(req.query["action"], check.Equals, "GetStatus")
	c.Assert(req.query["mediaid"], check.Equals, "abc123")
	c.Assert(req.query["extended"], check.Equals, "yes")
}

func (s *S) TestGetStatusMultiple(c *check.C) {
	server, requests := s.startServer(`
{
	"response": {
		"job": [
			{
				"id": "abc123",
				"userid": "myuser",
				"sourcefile": "http://some.video/file.mp4",
				"status": "Finished",
				"notifyurl": "http://ping.me/please",
				"created": "2015-12-31 20:45:30",
				"started": "2015-12-31 20:45:34",
				"finished": "2015-12-31 21:00:03",
				"prevstatus": "Saving",
				"downloaded": "2015-12-31 20:45:32",
				"uploaded": "2015-12-31 20:59:54",
				"time_left": "0",
				"progress": "100",
				"time_left_current": "0",
				"progress_current": "100.0",
				"format": [
					{
						"id": "f123",
						"status": "Finished",
						"created": "2015-12-31 20:45:30",
						"started": "2015-12-31 20:45:34",
						"finished": "2015-12-31 21:00:03",
						"s3_destination": "https://s3.amazonaws.com/not-really/valid.mp4",
						"cf_destination": "https://blablabla.cloudfront.net/not-valid.mp4",
						"convertedsize": "65724",
						"destination": [
							null,
							"s3://myunclebucket/file.mp4"
						],
						"destination_status": [
							null,
							"Saved"
						]
					},
					{
						"id": "f124",
						"status": "Finished",
						"created": "2015-12-31 20:45:30",
						"started": "2015-12-31 20:45:34",
						"finished": "2015-12-31 21:00:03",
						"s3_destination": "https://s3.amazonaws.com/not-really/valid.mp4",
						"cf_destination": "https://blablabla.cloudfront.net/not-valid.mp4",
						"convertedsize": "65725",
						"destination": [
							"s3://mynicebucket/file.mp4",
							"s3://myunclebucket/file.mp4"
						],
						"destination_status": null
					}
				]
			},
			{
				"id": "abc124",
				"userid": "myuser",
				"sourcefile": "http://some.video/file.mp4",
				"status": "Finished",
				"notifyurl": "http://ping.me/please",
				"created": "2015-12-31 20:45:30",
				"started": "2015-12-31 20:45:34",
				"finished": "2015-12-31 21:00:03",
				"prevstatus": "Saving",
				"downloaded": "2015-12-31 20:45:32",
				"uploaded": "2015-12-31 20:59:54",
				"time_left": "0",
				"progress": "100",
				"time_left_current": "0",
				"progress_current": "100.0",
				"format": {
						"id": "f123",
						"status": "Finished",
						"created": "2015-12-31 20:45:30",
						"started": "2015-12-31 20:45:34",
						"finished": "2015-12-31 21:00:03",
						"s3_destination": "https://s3.amazonaws.com/not-really/valid.mp4",
						"cf_destination": "https://blablabla.cloudfront.net/not-valid.mp4",
						"convertedsize": "65726",
						"destination": [
							"s3://mynicebucket/file.mp4",
							"s3://myunclebucket/file.mp4"
						],
						"destination_status": [
							"Saved",
							"Saved"
						]
					}
			}
		]
	}
}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	status, err := client.GetStatus([]string{"abc123", "abc124"}, true)
	c.Assert(err, check.IsNil)

	expectedCreateDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:30")
	expectedStartDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:34")
	expectedFinishDate, _ := time.Parse(dateTimeLayout, "2015-12-31 21:00:03")
	expectedDownloadDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:32")
	expectedUploadDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:59:54")
	expected := []StatusResponse{
		{
			MediaID:             "abc123",
			UserID:              "myuser",
			SourceFile:          "http://some.video/file.mp4",
			MediaStatus:         "Finished",
			PreviousMediaStatus: "Saving",
			NotifyURL:           "http://ping.me/please",
			CreateDate:          expectedCreateDate,
			StartDate:           expectedStartDate,
			FinishDate:          expectedFinishDate,
			DownloadDate:        expectedDownloadDate,
			UploadDate:          expectedUploadDate,
			TimeLeft:            "0",
			Progress:            100.0,
			TimeLeftCurrentJob:  "0",
			ProgressCurrentJob:  100.0,
			Formats: []FormatStatus{
				{
					ID:            "f123",
					Status:        "Finished",
					CreateDate:    expectedCreateDate,
					StartDate:     expectedStartDate,
					FinishDate:    expectedFinishDate,
					S3Destination: "https://s3.amazonaws.com/not-really/valid.mp4",
					CFDestination: "https://blablabla.cloudfront.net/not-valid.mp4",
					FileSize:      "65724",
					Destinations: []DestinationStatus{
						{Name: "", Status: ""},
						{Name: "s3://myunclebucket/file.mp4", Status: "Saved"},
					},
				},
				{
					ID:            "f124",
					Status:        "Finished",
					CreateDate:    expectedCreateDate,
					StartDate:     expectedStartDate,
					FinishDate:    expectedFinishDate,
					S3Destination: "https://s3.amazonaws.com/not-really/valid.mp4",
					CFDestination: "https://blablabla.cloudfront.net/not-valid.mp4",
					FileSize:      "65725",
					Destinations: []DestinationStatus{
						{Name: "s3://mynicebucket/file.mp4", Status: ""},
						{Name: "s3://myunclebucket/file.mp4", Status: ""},
					},
				},
			},
		},
		{
			MediaID:             "abc124",
			UserID:              "myuser",
			SourceFile:          "http://some.video/file.mp4",
			MediaStatus:         "Finished",
			PreviousMediaStatus: "Saving",
			NotifyURL:           "http://ping.me/please",
			CreateDate:          expectedCreateDate,
			StartDate:           expectedStartDate,
			FinishDate:          expectedFinishDate,
			DownloadDate:        expectedDownloadDate,
			UploadDate:          expectedUploadDate,
			TimeLeft:            "0",
			Progress:            100.0,
			TimeLeftCurrentJob:  "0",
			ProgressCurrentJob:  100.0,
			Formats: []FormatStatus{
				{
					ID:            "f123",
					Status:        "Finished",
					CreateDate:    expectedCreateDate,
					StartDate:     expectedStartDate,
					FinishDate:    expectedFinishDate,
					S3Destination: "https://s3.amazonaws.com/not-really/valid.mp4",
					CFDestination: "https://blablabla.cloudfront.net/not-valid.mp4",
					FileSize:      "65726",
					Destinations: []DestinationStatus{
						{Name: "s3://mynicebucket/file.mp4", Status: "Saved"},
						{Name: "s3://myunclebucket/file.mp4", Status: "Saved"},
					},
				},
			},
		},
	}
	c.Assert(status, check.DeepEquals, expected)

	req := <-requests
	c.Assert(req.query["action"], check.Equals, "GetStatus")
	c.Assert(req.query["mediaid"], check.Equals, "abc123,abc124")
	c.Assert(req.query["extended"], check.Equals, "yes")
}

// Some data are only available when extended=no
func (s *S) TestGetStatusNotExtended(c *check.C) {
	server, requests := s.startServer(nonExtendedStatus)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	status, err := client.GetStatus([]string{"abc123"}, false)
	c.Assert(err, check.IsNil)

	expectedCreateDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:30")
	expectedStartDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:34")
	expectedFinishDate, _ := time.Parse(dateTimeLayout, "2015-12-31 21:00:03")
	expectedDownloadDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:32")
	expected := []StatusResponse{
		{
			MediaID:      "abc123",
			UserID:       "myuser",
			SourceFile:   "http://some.video/file.mp4",
			MediaStatus:  "Finished",
			CreateDate:   expectedCreateDate,
			StartDate:    expectedStartDate,
			FinishDate:   expectedFinishDate,
			DownloadDate: expectedDownloadDate,
			TimeLeft:     "21",
			Progress:     100.0,
			Formats: []FormatStatus{
				{
					ID:           "f123",
					Status:       "Finished",
					CreateDate:   expectedCreateDate,
					StartDate:    expectedStartDate,
					FinishDate:   expectedFinishDate,
					Destinations: []DestinationStatus{{Name: "s3://mynicebucket", Status: "Saved"}},
					Size:         "0x1080",
					Bitrate:      "3500k",
					Output:       "mp4",
					VideoCodec:   "libx264",
					AudioCodec:   "dolby_aac",
					FileSize:     "78544430",
				},
			},
		},
	}
	c.Assert(status, check.DeepEquals, expected)

	req := <-requests
	c.Assert(req.query["action"], check.Equals, "GetStatus")
	c.Assert(req.query["mediaid"], check.Equals, "abc123")
	c.Assert(req.query["extended"], check.IsNil)
}

func (s *S) TestGetStatusZeroTime(c *check.C) {
	server, _ := s.startServer(`
{
	"response": {
		"job": {
			"id":"abc123",
			"userid":"myuser",
			"sourcefile":"http://some.file/wait-wat",
			"status":"Error",
			"created":"2016-01-29 19:32:32",
			"started":"2016-01-29 19:32:32",
			"finished":"0000-00-00 00:00:00",
			"downloaded":"0000-00-00 00:00:00",
			"description":"Download error:  The requested URL returned error: 403 Forbidden",
			"processor":"AMAZON",
			"region":"oak-private-clive",
			"time_left":"50",
			"progress":"50.0",
			"time_left_current":"0",
			"progress_current":"0.0",
			"format":{
				"id":"164478401",
				"status":"New",
				"created":"2016-01-29 19:32:32",
				"started":"0000-00-00 00:00:00",
				"finished":"0000-00-00 00:00:00",
				"destination":"http://s4.amazonaws.com/future",
				"destination_status":"Open",
				"convertedsize":"0",
				"queued":"0000-00-00 00:00:00",
				"converttime":"0",
				"time_left":"40",
				"progress":"0.0",
				"time_left_current":"0",
				"progress_current":"0.0"
			},
			"queue_time":"0"
		}
	}
}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	status, err := client.GetStatus([]string{"abc123"}, true)
	c.Assert(err, check.IsNil)
	c.Assert(status, check.HasLen, 1)

	expectedCreateDate, _ := time.Parse(dateTimeLayout, "2016-01-29 19:32:32")
	c.Assert(status[0].MediaID, check.Equals, "abc123")
	c.Assert(status[0].CreateDate, check.DeepEquals, expectedCreateDate)
	c.Assert(status[0].FinishDate.IsZero(), check.Equals, true)
	c.Assert(status[0].DownloadDate.IsZero(), check.Equals, true)
}

func (s *S) TestGetStatusNoMedia(c *check.C) {
	var client Client
	status, err := client.GetStatus(nil, true)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "please provide at least one media id")
	c.Assert(status, check.HasLen, 0)
}

func (s *S) TestGetStatusError(c *check.C) {
	server, _ := s.startServer(`{"response": {"message": "", "errors": {"error": "wait what?"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.GetStatus([]string{"some-media"}, true)
	c.Assert(err, check.NotNil)
	c.Assert(resp, check.IsNil)
}
