package encodingcom

import (
	"time"

	"gopkg.in/check.v1"
)

func (s *S) TestAddMedia(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "Added", "MediaID": "1234567"}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	format := Format{
		Output:       []string{"http://another.non.existent/video.mp4"},
		VideoCodec:   "x264",
		AudioCodec:   "aac",
		Bitrate:      "900k",
		AudioBitrate: "64k",
	}
	addMediaResponse, err := client.AddMedia([]string{"http://another.non.existent/video.mov"},
		[]Format{format}, "us-east-1")

	c.Assert(err, check.IsNil)
	c.Assert(addMediaResponse, check.DeepEquals, &AddMediaResponse{
		Message: "Added",
		MediaID: "1234567",
	})
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "AddMedia")
}

func (s *S) TestAddMediaError(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "Added", "errors": {"error": "something went wrong"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	format := Format{
		Output:       []string{"http://another.non.existent/video.mp4"},
		VideoCodec:   "x264",
		AudioCodec:   "aac",
		Bitrate:      "900k",
		AudioBitrate: "64k",
	}
	addMediaResponse, err := client.AddMedia([]string{"http://another.non.existent/video.mov"},
		[]Format{format}, "us-east-1")
	c.Assert(err, check.NotNil)
	c.Assert(addMediaResponse, check.IsNil)
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "AddMedia")
}

func (s *S) TestStopMedia(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "Stopped"}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.StopMedia("some-media")
	c.Assert(err, check.IsNil)
	c.Assert(resp, check.DeepEquals, &Response{Message: "Stopped"})
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "StopMedia")
	c.Assert(req.query["mediaid"], check.Equals, "some-media")
}

func (s *S) TestStopMediaError(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "failed", "errors": {"error": "something went wrong"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.StopMedia("some-media")
	c.Assert(err, check.NotNil)
	c.Assert(resp, check.IsNil)
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "StopMedia")
}

func (s *S) TestCancelMedia(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "Canceled"}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.CancelMedia("some-media")
	c.Assert(err, check.IsNil)
	c.Assert(resp, check.DeepEquals, &Response{Message: "Canceled"})
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "CancelMedia")
	c.Assert(req.query["mediaid"], check.Equals, "some-media")
}

func (s *S) TestCancelMediaError(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "failed", "errors": {"error": "something went wrong"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.CancelMedia("some-media")
	c.Assert(err, check.NotNil)
	c.Assert(resp, check.IsNil)
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "CancelMedia")
}

func (s *S) TestRestartMedia(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "Restarted"}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.RestartMedia("some-media", false)
	c.Assert(err, check.IsNil)
	c.Assert(resp, check.DeepEquals, &Response{Message: "Restarted"})
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "RestartMedia")
	c.Assert(req.query["mediaid"], check.Equals, "some-media")
}

func (s *S) TestRestartMediaWithErrors(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "Restarted"}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.RestartMedia("some-media", true)
	c.Assert(err, check.IsNil)
	c.Assert(resp, check.DeepEquals, &Response{Message: "Restarted"})
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "RestartMediaErrors")
	c.Assert(req.query["mediaid"], check.Equals, "some-media")
}

func (s *S) TestRestartMediaError(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "failed", "errors": {"error": "something went wrong"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.RestartMedia("some-media", false)
	c.Assert(err, check.NotNil)
	c.Assert(resp, check.IsNil)
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "RestartMedia")
}

func (s *S) TestRestartMediaTask(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "Task restarted"}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.RestartMediaTask("some-media", "some-task")
	c.Assert(err, check.IsNil)
	c.Assert(resp, check.DeepEquals, &Response{Message: "Task restarted"})
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "RestartMediaTask")
	c.Assert(req.query["mediaid"], check.Equals, "some-media")
	c.Assert(req.query["taskid"], check.Equals, "some-task")
}

func (s *S) TestRestartMediaTaskError(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "Failed to restart", "errors": {"error": "something went really bad"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.RestartMediaTask("some-media", "some-task")
	c.Assert(err, check.NotNil)
	c.Assert(resp, check.IsNil)
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "RestartMediaTask")
}

func (s *S) TestListMedia(c *check.C) {
	server, requests := s.startServer(`
{
    "response":{
        "media":[
            {
                "mediafile":"http://another.non.existent/video.mp4",
                "mediaid":"1234567",
                "mediastatus":"Finished",
                "createdate":"2015-12-31 20:45:30",
                "startdate":"2015-12-31 20:45:50",
                "finishdate":"2015-12-31 20:48:54"
            }
        ]
    }
}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	listMediaResponse, err := client.ListMedia()
	c.Assert(err, check.IsNil)

	mockCreateDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:30")
	mockStartDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:45:50")
	mockFinishDate, _ := time.Parse(dateTimeLayout, "2015-12-31 20:48:54")

	c.Assert(listMediaResponse, check.DeepEquals, &ListMediaResponse{
		Media: []ListMediaResponseItem{
			{
				"http://another.non.existent/video.mp4",
				"1234567",
				"Finished",
				MediaDateTime{mockCreateDate},
				MediaDateTime{mockStartDate},
				MediaDateTime{mockFinishDate},
			},
		},
	})
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "GetMediaList")
}

func (s *S) TestListMediaError(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "", "errors": {"error": "can't list"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.ListMedia()
	c.Assert(err, check.NotNil)
	c.Assert(resp, check.IsNil)
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "GetMediaList")
}

func (s *S) TestGetMediaInfo(c *check.C) {
	server, requests := s.startServer(`
{
	"response": {
		"bitrate": "1807k",
		"duration": "6464.83",
		"audio_bitrate": "128k",
		"video_codec": "mpeg4",
		"video_bitrate": "1679k",
		"frame_rate": "23.98",
		"size": "640x352",
		"pixel_aspect_ratio": "1:1",
		"display_aspect_ratio": "20:11",
		"audio_codec": "ac3",
		"audio_sample_rate": "48000",
		"audio_channels": "2",
		"rotation":"90"
	}
}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	mediaInfo, err := client.GetMediaInfo("m-123")
	c.Assert(err, check.IsNil)

	c.Assert(mediaInfo, check.DeepEquals, &MediaInfo{
		Bitrate:            "1807k",
		Duration:           6464*time.Second + time.Duration(0.83*float64(time.Second)),
		VideoCodec:         "mpeg4",
		VideoBitrate:       "1679k",
		Framerate:          "23.98",
		Size:               "640x352",
		PixelAspectRatio:   "1:1",
		DisplayAspectRatio: "20:11",
		AudioCodec:         "ac3",
		AudioSampleRate:    uint(48000),
		AudioChannels:      "2",
		AudioBitrate:       "128k",
		Rotation:           90,
	})

	req := <-requests
	c.Assert(req.query["action"], check.Equals, "GetMediaInfo")
	c.Assert(req.query["mediaid"], check.Equals, "m-123")
}

func (s *S) TestGetMediaInfoError(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "", "errors": {"error": "wait what?"}}}`)
	defer server.Close()

	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	resp, err := client.GetMediaInfo("some-media")
	c.Assert(err, check.NotNil)
	c.Assert(resp, check.IsNil)
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "GetMediaInfo")
}
