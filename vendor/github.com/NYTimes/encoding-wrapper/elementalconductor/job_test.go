package elementalconductor

import (
	"encoding/xml"
	"net/http"

	"gopkg.in/check.v1"
)

func (s *S) TestGetJobError(c *check.C) {
	errorResponse := `<?xml version="1.0" encoding="UTF-8"?>
<errors>
  <error type="ActiveRecord::RecordNotFound">Couldn't find Job with id=1</error>
</errors>`
	server, _ := s.startServer(http.StatusNotFound, errorResponse)
	defer server.Close()
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	getJobsResponse, err := client.GetJob("1")
	c.Assert(getJobsResponse, check.IsNil)
	c.Assert(err, check.DeepEquals, &APIError{
		Status: http.StatusNotFound,
		Errors: errorResponse,
	})
}

func (s *S) TestGetJobsOnEmptyList(c *check.C) {
	server, _ := s.startServer(http.StatusOK, `<?xml version="1.0" encoding="UTF-8"?>
<job_list>
  <empty>There are currently no jobs</empty>
</job_list>`)
	defer server.Close()
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	getJobsResponse, err := client.GetJobs()
	c.Assert(err, check.IsNil)
	c.Assert(getJobsResponse, check.DeepEquals, &JobList{
		XMLName: xml.Name{
			Local: "job_list",
		},
		Empty: "There are currently no jobs",
	})
}

func (s *S) TestCreateJob(c *check.C) {
	jobResponseXML := `<job href="/jobs/1">
    <input>
        <file_input>
            <uri>http://another.non.existent/video.mp4</uri>
            <username>user</username>
            <password>pass123</password>
        </file_input>
    </input>
    <priority>50</priority>
    <output_group>
        <order>1</order>
        <file_group_settings>
            <destination>
                <uri>http://destination/video.mp4</uri>
                <username>user</username>
                <password>pass123</password>
            </destination>
        </file_group_settings>
        <apple_live_group_settings>
            <destination>
                <uri>http://destination/video.mp4</uri>
                <username>user</username>
                <password>pass123</password>
            </destination>
        </apple_live_group_settings>
        <type>apple_live_group_settings</type>
        <output>
            <stream_assembly_name>stream_1</stream_assembly_name>
            <name_modifier>_high</name_modifier>
            <order>1</order>
            <extension>.mp4</extension>
        </output>
    </output_group>
    <stream_assembly>
        <name>stream_1</name>
        <preset>17</preset>
    </stream_assembly>
</job>`
	server, _ := s.startServer(http.StatusCreated, jobResponseXML)
	defer server.Close()
	jobInput := Job{
		XMLName: xml.Name{
			Local: "job",
		},
		Href: "/jobs/1",
		Input: Input{
			FileInput: Location{
				URI:      "http://another.non.existent/video.mp4",
				Username: "user",
				Password: "pass123",
			},
		},
		Priority: 50,
		OutputGroup: []OutputGroup{
			{
				Order: 1,
				FileGroupSettings: &FileGroupSettings{
					Destination: &Location{
						URI:      "http://destination/video.mp4",
						Username: "user",
						Password: "pass123",
					},
				},
				AppleLiveGroupSettings: &AppleLiveGroupSettings{
					Destination: &Location{
						URI:      "http://destination/video.mp4",
						Username: "user",
						Password: "pass123",
					},
				},
				Type: AppleLiveOutputGroupType,
				Output: []Output{
					{
						StreamAssemblyName: "stream_1",
						NameModifier:       "_high",
						Order:              1,
						Extension:          ".mp4",
					},
				},
			},
		},
		StreamAssembly: []StreamAssembly{
			{
				Name:   "stream_1",
				Preset: "17",
			},
		},
	}
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	postJobResponse, err := client.CreateJob(&jobInput)
	c.Assert(err, check.IsNil)
	c.Assert(postJobResponse, check.NotNil)
	c.Assert(postJobResponse, check.DeepEquals, &jobInput)
}

func (s *S) TestGetJob(c *check.C) {
	jobResponseXML := `<job href="/jobs/1">
    <input>
        <file_input>
            <uri>http://another.non.existent/video.mp4</uri>
            <username>user</username>
            <password>pass123</password>
        </file_input>
        <input_info>
            <general>
                <format>MPEG-4</format>
                <format_profile>QuickTime</format_profile>
                <codec_id>qt    </codec_id>
                <file_size>185 MiB</file_size>
                <duration>1mn 19s</duration>
                <overall_bit_rate>19.4 Mbps</overall_bit_rate>
            </general>
            <video>
                <id>1</id>
                <format>AVC</format>
                <format_info>Advanced Video Codec</format_info>
                <format_profile>Main@L4.1</format_profile>
                <format_settings__cabac>No</format_settings__cabac>
                <format_settings__reframes>2 frames</format_settings__reframes>
                <format_settings__gop>M=1, N=60</format_settings__gop>
                <codec_id>avc1</codec_id>
                <codec_id_info>Advanced Video Coding</codec_id_info>
                <bit_rate>19.2 Mbps</bit_rate>
                <width>1 920 pixels</width>
                <height>1 080 pixels</height>
            </video>
            <audio>
                <id>2</id>
                <format>AAC</format>
                <format_info>Advanced Audio Codec</format_info>
                <format_profile>LC</format_profile>
            </audio>
            <other>
                <id>3</id>
                <type>Time code</type>
                <format>QuickTime TC</format>
            </other>
        </input_info>
    </input>
    <content_duration>
        <input_duration>716</input_duration>
        <clipped_input_duration>716</clipped_input_duration>
        <stream_count>1</stream_count>
        <total_stream_duration>716</total_stream_duration>
        <package_count>1</package_count>
        <total_package_duration>716</total_package_duration>
    </content_duration>
    <priority>50</priority>
    <output_group>
        <order>1</order>
        <file_group_settings>
            <destination>
                <uri>http://destination/video.mp4</uri>
                <username>user</username>
                <password>pass123</password>
            </destination>
        </file_group_settings>
        <type>file_group_settings</type>
        <output>
            <full_uri>s3://mybucket/some/dir/mynicefile.mp4</full_uri>
            <stream_assembly_name>stream_1</stream_assembly_name>
            <name_modifier>_high</name_modifier>
            <order>1</order>
            <extension>.mp4</extension>
        </output>
    </output_group>
    <stream_assembly>
        <id>1146</id>
        <name>stream_1</name>
        <preset>17</preset>
        <video_description>
            <afd_signaling>None</afd_signaling>
            <anti_alias>true</anti_alias>
            <drop_frame_timecode>true</drop_frame_timecode>
            <encoder_type>gpu</encoder_type>
            <fixed_afd nil="true"/>
            <force_cpu_encode>false</force_cpu_encode>
            <height>1080</height>
            <id>1366</id>
            <insert_color_metadata>false</insert_color_metadata>
            <respond_to_afd>None</respond_to_afd>
            <sharpness>50</sharpness>
            <stretch_to_output>false</stretch_to_output>
            <timecode_passthrough>false</timecode_passthrough>
            <vbi_passthrough>false</vbi_passthrough>
            <width nil="true"/>
            <gpu>0</gpu>
            <selected_gpu nil="true"/>
            <codec>h.264</codec>
        </video_description>
    </stream_assembly>
</job>`
	server, _ := s.startServer(http.StatusOK, jobResponseXML)
	defer server.Close()
	expectedJob := Job{
		XMLName: xml.Name{
			Local: "job",
		},
		Href: "/jobs/1",
		Input: Input{
			FileInput: Location{
				URI:      "http://another.non.existent/video.mp4",
				Username: "user",
				Password: "pass123",
			},
			InputInfo: &InputInfo{
				Video: VideoInputInfo{
					Bitrate:       "19.2 Mbps",
					Format:        "AVC",
					FormatInfo:    "Advanced Video Codec",
					FormatProfile: "Main@L4.1",
					CodecID:       "avc1",
					CodecIDInfo:   "Advanced Video Coding",
					Width:         "1 920 pixels",
					Height:        "1 080 pixels",
				},
			},
		},
		ContentDuration: &ContentDuration{InputDuration: 716},
		Priority:        50,
		OutputGroup: []OutputGroup{
			{
				Order: 1,
				FileGroupSettings: &FileGroupSettings{
					Destination: &Location{
						URI:      "http://destination/video.mp4",
						Username: "user",
						Password: "pass123",
					},
				},
				Type: FileOutputGroupType,
				Output: []Output{
					{
						FullURI:            "s3://mybucket/some/dir/mynicefile.mp4",
						StreamAssemblyName: "stream_1",
						NameModifier:       "_high",
						Order:              1,
						Extension:          ".mp4",
					},
				},
			},
		},
		StreamAssembly: []StreamAssembly{
			{
				ID:     "1146",
				Name:   "stream_1",
				Preset: "17",
				VideoDescription: &StreamVideoDescription{
					Codec:       "h.264",
					EncoderType: "gpu",
					Width:       "",
					Height:      "1080",
				},
			},
		},
	}
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	getJobsResponse, err := client.GetJob("1")
	c.Assert(err, check.IsNil)
	c.Assert(getJobsResponse, check.NotNil)
	c.Assert(*getJobsResponse, check.DeepEquals, expectedJob)
}

func (s *S) TestJobGetID(c *check.C) {
	var tests = []struct {
		href string
		id   string
	}{
		{
			"http://myelemental/jobs/123",
			"123",
		},
		{
			"job-1234",
			"job-1234",
		},
		{
			"",
			"",
		},
	}
	for _, test := range tests {
		j := Job{Href: test.href}
		id := j.GetID()
		c.Check(id, check.Equals, test.id)
	}
}

func (s *S) TestCancelJob(c *check.C) {
	jobResponseXML := `<job href="/jobs/1">
    <status>canceled</status>
    <input>
        <file_input>
            <uri>http://another.non.existent/video.mp4</uri>
            <username>user</username>
            <password>pass123</password>
        </file_input>
    </input>
    <priority>50</priority>
    <output_group>
        <order>1</order>
        <file_group_settings>
            <destination>
                <uri>http://destination/video.mp4</uri>
                <username>user</username>
                <password>pass123</password>
            </destination>
        </file_group_settings>
        <type>file_group_settings</type>
        <output>
            <stream_assembly_name>stream_1</stream_assembly_name>
            <name_modifier>_high</name_modifier>
            <order>1</order>
            <extension>.mp4</extension>
        </output>
    </output_group>
    <stream_assembly>
        <name>stream_1</name>
        <preset>17</preset>
    </stream_assembly>
</job>`
	server, reqs := s.startServer(http.StatusOK, jobResponseXML)
	defer server.Close()
	expectedJob := Job{
		XMLName: xml.Name{
			Local: "job",
		},
		Href: "/jobs/1",
		Input: Input{
			FileInput: Location{
				URI:      "http://another.non.existent/video.mp4",
				Username: "user",
				Password: "pass123",
			},
		},
		Priority: 50,
		OutputGroup: []OutputGroup{
			{
				Order: 1,
				FileGroupSettings: &FileGroupSettings{
					Destination: &Location{
						URI:      "http://destination/video.mp4",
						Username: "user",
						Password: "pass123",
					},
				},
				Type: FileOutputGroupType,
				Output: []Output{
					{
						StreamAssemblyName: "stream_1",
						NameModifier:       "_high",
						Order:              1,
						Extension:          ".mp4",
					},
				},
			},
		},
		StreamAssembly: []StreamAssembly{
			{
				Name:   "stream_1",
				Preset: "17",
			},
		},
		Status: "canceled",
	}
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	job, err := client.CancelJob("1")
	c.Assert(err, check.IsNil)
	c.Assert(job, check.NotNil)
	c.Assert(job, check.DeepEquals, &expectedJob)

	req := <-reqs
	c.Assert(req.req.Method, check.Equals, "POST")
	c.Assert(req.req.URL.Path, check.Equals, "/api/jobs/1/cancel")
	c.Assert(string(req.body), check.Equals, "<cancel></cancel>")
}

func (s *S) TestCancelJobError(c *check.C) {
	errorResponse := `<?xml version="1.0" encoding="UTF-8"?>
<errors>
  <error type="ActiveRecord::RecordNotFound">Couldn't find Job with id=1</error>
</errors>`
	server, _ := s.startServer(http.StatusNotFound, errorResponse)
	defer server.Close()
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	job, err := client.CancelJob("1")
	c.Assert(job, check.IsNil)
	c.Assert(err, check.DeepEquals, &APIError{
		Status: http.StatusNotFound,
		Errors: errorResponse,
	})
}

func (s *S) TestVideoInfoDimensions(c *check.C) {
	var tests = []struct {
		inputWidth     string
		inputHeight    string
		expectedWith   int64
		expectedHeight int64
	}{
		{
			"1 920 pixels",
			"1 080 pixels",
			1920,
			1080,
		},
		{
			"1280 pixels",
			"720 pixels",
			1280,
			720,
		},
		{
			"1 280 pixels",
			"720 pixels",
			1280,
			720,
		},
		{
			"1,280 pixels",
			"720 pixels",
			1280,
			720,
		},
		{
			"1920p",
			"1080p",
			1920,
			1080,
		},
		{
			"1280",
			"720",
			1280,
			720,
		},
		{
			"twelve eighty",
			"seven twenty",
			0,
			0,
		},
	}
	for _, t := range tests {
		job := VideoInputInfo{Width: t.inputWidth, Height: t.inputHeight}
		width := job.GetWidth()
		height := job.GetHeight()
		if width != t.expectedWith {
			c.Errorf("width=%s height=%s\nwant width=%d\ngot  width=%d", t.inputWidth, t.inputHeight, t.expectedWith, width)
		}
		if height != t.expectedHeight {
			c.Errorf("width=%s height=%s\nwant height=%d\ngot  height=%d", t.inputWidth, t.inputHeight, t.expectedHeight, height)
		}
	}
}

func (s *S) TestVideoDescriptionWidth(c *check.C) {
	var tests = []struct {
		input  string
		output int64
	}{
		{"1920", 1920},
		{"", 0},
		{"whatever", 0},
	}
	for _, t := range tests {
		desc := StreamVideoDescription{Width: t.input}
		got := desc.GetWidth()
		c.Check(got, check.Equals, t.output)
	}
}

func (s *S) TestVideoDescriptionHeight(c *check.C) {
	var tests = []struct {
		input  string
		output int64
	}{
		{"1080", 1080},
		{"", 0},
		{"whatever", 0},
	}
	for _, t := range tests {
		desc := StreamVideoDescription{Height: t.input}
		got := desc.GetHeight()
		c.Check(got, check.Equals, t.output)
	}
}
