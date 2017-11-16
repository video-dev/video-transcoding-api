package elementalconductor

import (
	"encoding/xml"
	"regexp"
	"strconv"
	"strings"
)

var nonDigitRegexp = regexp.MustCompile(`[^\d]`)

// OutputGroupType is a custom type for OutputGroup type field values
type OutputGroupType string

const (
	// FileOutputGroupType is the value for the type field on OutputGroup
	// for jobs with a file output
	FileOutputGroupType = OutputGroupType("file_group_settings")
	// AppleLiveOutputGroupType is the value for the type field on OutputGroup
	// for jobs with Apple's HTTP Live Streaming (HLS) output
	AppleLiveOutputGroupType = OutputGroupType("apple_live_group_settings")
)

// Container is the Video container type for a job
type Container string

const (
	// AppleHTTPLiveStreaming is the container for HLS video files
	AppleHTTPLiveStreaming = Container("m3u8")
	// MPEG4 is the container for MPEG-4 video files
	MPEG4 = Container("mp4")
)

// GetJobs returns a list of the user's jobs
func (c *Client) GetJobs() (*JobList, error) {
	var result *JobList
	err := c.do("GET", "/jobs", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetJob returns metadata on a single job
func (c *Client) GetJob(jobID string) (*Job, error) {
	var result *Job
	err := c.do("GET", "/jobs/"+jobID, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateJob sends a single job to the current Elemental
// Cloud deployment for processing
func (c *Client) CreateJob(job *Job) (*Job, error) {
	var result *Job
	err := c.do("POST", "/jobs", *job, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CancelJob cancels the given job in the Elemental Conductor API.
func (c *Client) CancelJob(jobID string) (*Job, error) {
	var job *Job
	var payload = struct {
		XMLName xml.Name `xml:"cancel"`
	}{}
	err := c.do("POST", "/jobs/"+jobID+"/cancel", payload, &job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// GetID is a convenience function to parse the job id
// out of the Href attribute in Job
func (j *Job) GetID() string {
	if j.Href != "" {
		hrefData := strings.Split(j.Href, "/")
		return hrefData[len(hrefData)-1]
	}
	return ""
}

// JobList represents the response returned by
// a query for the list of jobs
type JobList struct {
	XMLName xml.Name `xml:"job_list"`
	Empty   string   `xml:"empty,omitempty"`
	Job     []Job    `xml:"job"`
}

// Job represents a job to be sent to Elemental Cloud
type Job struct {
	XMLName         xml.Name         `xml:"job"`
	Href            string           `xml:"href,attr,omitempty"`
	Input           Input            `xml:"input,omitempty"`
	ContentDuration *ContentDuration `xml:"content_duration,omitempty"`
	Priority        int              `xml:"priority,omitempty"`
	OutputGroup     []OutputGroup    `xml:"output_group,omitempty"`
	StreamAssembly  []StreamAssembly `xml:"stream_assembly,omitempty"`
	Status          string           `xml:"status,omitempty"`
	Submitted       DateTime         `xml:"submitted,omitempty"`
	StartTime       DateTime         `xml:"start_time,omitempty"`
	CompleteTime    DateTime         `xml:"complete_time,omitempty"`
	ErroredTime     DateTime         `xml:"errored_time,omitempty"`
	PercentComplete int              `xml:"pct_complete,omitempty"`
	ErrorMessages   []JobError       `xml:"error_messages,omitempty"`
}

// JobError represents an individual error on a job
type JobError struct {
	Code      int              `xml:"error>code,omitempty"`
	CreatedAt JobErrorDateTime `xml:"error>created_at,omitempty"`
	Message   string           `xml:"error>message,omitempty"`
}

// Input represents the spec for the job's input
type Input struct {
	FileInput Location   `xml:"file_input,omitempty"`
	InputInfo *InputInfo `xml:"input_info,omitempty"`
}

// InputInfo contains metadata related to a job input.
type InputInfo struct {
	Video VideoInputInfo `xml:"video"`
}

// VideoInputInfo contains video metadata related to a job input.
type VideoInputInfo struct {
	Format        string `xml:"format"`
	FormatInfo    string `xml:"format_info"`
	FormatProfile string `xml:"format_profile"`
	CodecID       string `xml:"codec_id"`
	CodecIDInfo   string `xml:"codec_id_info"`
	Bitrate       string `xml:"bit_rate"`
	Width         string `xml:"width"`
	Height        string `xml:"height"`
}

// GetWidth parses the underlying width returned the Elemental Conductor API
// and converts it to int64.
//
// Examples:
//  - Input: "1 920 pixels"
//    Output: 1920
//  - Input: "1920p"
//    Output: 1920
//  - Input: "1 920"
//    Output: 1920
func (v *VideoInputInfo) GetWidth() int64 {
	return v.extractNumber(v.Width)
}

// GetHeight parses the underlying height returned the Elemental Conductor API
// and converts it to int64.
//
// Examples:
//  - Input: "1 080 pixels"
//    Output: 1080
//  - Input: "1080p"
//    Output: 1080
//  - Input: "1 080"
//    Output: 1080
func (v *VideoInputInfo) GetHeight() int64 {
	return v.extractNumber(v.Height)
}

func (v *VideoInputInfo) extractNumber(input string) int64 {
	input = nonDigitRegexp.ReplaceAllString(input, "")
	n, _ := strconv.ParseInt(input, 10, 64)
	return n
}

// ContentDuration contains information about the content of the media in the
// job.
type ContentDuration struct {
	InputDuration int `xml:"input_duration"`
}

// Location defines where a file is or needs to be.
// Username and Password are required for certain
// protocols that require authentication, like S3
type Location struct {
	URI      string `xml:"uri,omitempty"`
	Username string `xml:"username,omitempty"`
	Password string `xml:"password,omitempty"`
}

// OutputGroup is a list of the indended outputs for the job
type OutputGroup struct {
	Order                  int                     `xml:"order,omitempty"`
	FileGroupSettings      *FileGroupSettings      `xml:"file_group_settings,omitempty"`
	AppleLiveGroupSettings *AppleLiveGroupSettings `xml:"apple_live_group_settings,omitempty"`
	Type                   OutputGroupType         `xml:"type,omitempty"`
	Output                 []Output                `xml:"output,omitempty"`
}

// FileGroupSettings define where the file job output should go
type FileGroupSettings struct {
	Destination *Location `xml:"destination,omitempty"`
}

// AppleLiveGroupSettings define where the HLS job output should go
type AppleLiveGroupSettings struct {
	Destination     *Location `xml:"destination,omitempty"`
	SegmentDuration uint      `xml:"segment_length,omitempty"`
	EmitSingleFile  bool      `xml:"emit_single_file,omitempty"`
}

// Output defines the different processing stream assemblies
// for the job
type Output struct {
	FullURI            string    `xml:"full_uri,omitempty"`
	StreamAssemblyName string    `xml:"stream_assembly_name,omitempty"`
	NameModifier       string    `xml:"name_modifier,omitempty"`
	Order              int       `xml:"order,omitempty"`
	Extension          string    `xml:"extension,omitempty"`
	Container          Container `xml:"container,omitempty"`
}

// StreamAssembly defines how each processing stream should behave
type StreamAssembly struct {
	ID               string                  `xml:"id,omitempty"`
	Name             string                  `xml:"name,omitempty"`
	Preset           string                  `xml:"preset,omitempty"`
	VideoDescription *StreamVideoDescription `xml:"video_description"`
}

// StreamVideoDescription contains information about the video in a given
// stream assembly.
type StreamVideoDescription struct {
	Codec       string `xml:"codec"`
	EncoderType string `xml:"encoder_type"`
	Height      string `xml:"height"`
	Width       string `xml:"width"`
}

// GetWidth returns the underlying width parsed as an int64.
func (s *StreamVideoDescription) GetWidth() int64 {
	return s.getNumber(s.Width)
}

// GetHeight returns the underlying height parsed as an int64.
func (s *StreamVideoDescription) GetHeight() int64 {
	return s.getNumber(s.Height)
}

func (s *StreamVideoDescription) getNumber(input string) int64 {
	v, _ := strconv.ParseInt(input, 10, 64)
	return v
}
