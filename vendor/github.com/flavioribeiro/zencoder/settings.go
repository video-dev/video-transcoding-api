package zencoder

type AccessControlSettings struct {
	Grantee    string   `json:"grantee,omitempty"`    // Set the grantee for fine-grained S3 access_control permissions.
	Permission []string `json:"permission,omitempty"` // Set the permission for a grantee when using fine-grained access_control.
}

type ThumbnailSettings struct {
	Label               string                   `json:"label,omitempty"`                 // A label to identify each set of thumbnail groups.
	Format              string                   `json:"format,omitempty"`                // The format of the thumbnail image.
	Number              int32                    `json:"number,omitempty"`                // A number of thumbnails, evenly-spaced.
	StartAtFirstFrame   bool                     `json:"start_at_first_frame,omitempty"`  // Start generating the thumbnails starting at the first frame.
	Interval            int32                    `json:"interval,omitempty"`              // Take thumbnails at an even interval, in seconds.
	IntervalInFrames    int32                    `json:"interval_in_frames,omitempty"`    // Take thumbnails at an even interval, in frames.
	Times               []int32                  `json:"times,omitempty"`                 // An array of times, in seconds, at which to grab a thumbnail.
	AspectMode          string                   `json:"aspect_mode,omitempty"`           // How to handle a thumbnail width/height that differs from the aspect ratio of the input file.
	Size                string                   `json:"size,omitempty"`                  // Thumbnail resolution as WxH.
	Width               int32                    `json:"width,omitempty"`                 // The maximum width of the thumbnail (in pixels).
	Height              int32                    `json:"height,omitempty"`                // The maximum height of the thumbnail (in pixels).
	BaseUrl             string                   `json:"base_url,omitempty"`              // A base S3, Cloud Files, GCS, FTP, FTPS, or SFTP directory URL where we'll place the thumbnails, without a filename.
	Prefix              string                   `json:"prefix,omitempty"`                // Prefix for thumbnail filenames.
	Filename            string                   `json:"filename,omitempty"`              // Interpolated thumbnail filename.
	MakePublic          bool                     `json:"public,omitempty"`                // Make the output publicly readable on S3.
	AccessControl       []*AccessControlSettings `json:"access_control,omitempty"`        // Fine-grained access control rules for files sent to S3.
	UseRRS              bool                     `json:"rrs,omitempty"`                   // Amazon S3's Reduced Redundancy Storage.
	Headers             map[string]string        `json:"headers,omitempty"`               // HTTP headers to send with your thumbnails when we upload them.
	Credentials         string                   `json:"credentials,omitempty"`           // References saved credentials by a nickname.
	ParallelUploadLimit int32                    `json:"parallel_upload_limit,omitempty"` // The maximum number of simultaneous uploads to attempt.
}

type WatermarkSettings struct {
	Url     string  `json:"url,omitempty"`     // The URL of a remote image file to use as a watermark.
	X       int32   `json:"x,omitempty"`       // Where to place a watermark, on the x axis.
	Y       int32   `json:"y,omitempty"`       // Where to place a watermark, on the y axis.
	Width   int32   `json:"width,omitempty"`   // The scaled width of a watermark.
	Height  int32   `json:"height,omitempty"`  // The scaled height of a watermark.
	Origin  string  `json:"origin,omitempty"`  // Which part of the output to base the watermark position on.
	Opacity float64 `json:"opacity,omitempty"` // Make the watermark transparent.
}

type NotificationSettings struct {
	Url     string            `json:"url,omitempty"`     // Be notified when a job or output is complete.
	Format  string            `json:"format,omitempty"`  // A format and content type for notifications.
	Headers map[string]string `json:"headers,omitempty"` // Headers to pass along on HTTP notifications.
	Event   string            `json:"event,omitempty"`   // The event that triggers a notification. Used for Instant Play.
}

type StreamSettings struct {
	Source     string `json:"source,omitempty"`     // Specifies the label of a given source
	Path       string `json:"path,omitempty"`       // Specifies the path to a stream manifest file
	Bandwidth  int32  `json:"bandwidth,omitempty"`  // Specifies the bandwidth of a playlist stream
	Resolution string `json:"resolution,omitempty"` // Specifies the resolution of a playlist stream
	Codecs     string `json:"codecs,omitempty"`     // Specifies the codecs used in a playlist stream
}

type CuePointSettings struct {
	Type string            `json:"type,omitempty"` // A cue point type.
	Time float64           `json:"time,omitempty"` // A cue point time, in seconds.
	Name string            `json:"name,omitempty"` // A cue point name.
	Data map[string]string `json:"data,omitempty"` // Cue point data.
}

type OutputSettings struct {
	// General Output Settings
	Type                string            `json:"type,omitempty"`                  // The type of file to output.
	Label               string            `json:"label,omitempty"`                 // An optional label for this output.
	Url                 string            `json:"url,omitempty"`                   // A S3, Cloud Files, GCS, FTP, FTPS, SFTP, Aspera, HTTP, or RTMP URL where we'll put the transcoded file.
	SecondaryUrl        string            `json:"secondary_url,omitempty"`         // A S3, Cloud Files, GCS, FTP, FTPS, SFTP, Aspera, or syndication URL where we'll put the transcoded file.
	BaseUrl             string            `json:"base_url,omitempty"`              // A base S3, Cloud Files, GCS, FTP, FTPS, SFTP, or Aspera directory URL where we'll put the transcoded file, without a filename.
	Filename            string            `json:"filename,omitempty"`              // The filename of a finished file.
	PackageFilename     string            `json:"package_filename,omitempty"`      // The filename of a packaged output.
	PackageFormat       string            `json:"package_format,omitempty"`        // Zip/packaging format to use for the output file(s).
	DeviceProfile       string            `json:"device_profile,omitempty"`        // A device profile to use for mobile device compatibility.
	Strict              bool              `json:"strict,omitempty"`                // Enable strict mode.
	SkipVideo           bool              `json:"skip_video,omitempty"`            // Do not output a video track.
	SkipAudio           bool              `json:"skip_audio,omitempty"`            // Do not output an audio track.
	Source              string            `json:"source,omitempty"`                // References a label on another job and uses the video created by that output for processing instead of the input file.
	Credentials         string            `json:"credentials,omitempty"`           // References saved credentials by a nickname.
	GenerateMD5Checksum bool              `json:"generate_md5_checksum,omitempty"` // Generate an MD5 checksum of the output file.
	ParallelUploadLimit int32             `json:"parallel_upload_limit,omitempty"` // The maximum number of simultaneous uploads to attempt.
	Headers             map[string]string `json:"headers,omitempty"`               // HTTP headers to send with your file when we upload it.

	// Format And Codecs
	Format     string `json:"format,omitempty"`      // The output format to use.
	VideoCodec string `json:"video_codec,omitempty"` // The video codec to use.
	AudioCodec string `json:"audio_codec,omitempty"` // The audio codec to use.

	// Resolution
	Size       string `json:"size,omitempty"`        // The resolution of the output video (WxH, in pixels).
	Width      int32  `json:"width,omitempty"`       // The maximum width of the output video (in pixels).
	Height     int32  `json:"height,omitempty"`      // The maximum height of the output video (in pixels).
	Upscale    bool   `json:"upscale,omitempty"`     // Upscale the output if the input is smaller than the target output resolution.
	AspectMode string `json:"aspect_mode,omitempty"` // What to do when aspect ratio of input file does not match the target width/height aspect ratio.

	// Rate Control
	Quality              int32 `json:"quality,omitempty"`                // Autoselect the best video bitrate to to match a target visual quality.
	VideoBitrate         int32 `json:"video_bitrate,omitempty"`          // A target video bitrate in kbps. Not necessary if you select a quality setting, unless you want to target a specific bitrate.
	AudioQuality         int32 `json:"audio_quality,omitempty"`          // Autoselect the best audio bitrate to to match a target sound quality.
	AudioBitrate         int32 `json:"audio_bitrate,omitempty"`          // A target audio bitrate in kbps. Not necessary if you select a audio_quality setting, unless you want to target a specific bitrate.
	MaxVideoBitrate      int32 `json:"max_video_bitrate,omitempty"`      // A maximum average bitrate.
	Speed                int32 `json:"speed,omitempty"`                  // A target transcoding speed. Slower encoding generally allows for more advanced compression.
	DecoderBitrateCap    int32 `json:"decoder_bitrate_cap,omitempty"`    // Max bitrate fed to decoder buffer. Typically used for video intended for streaming, or for targeting specific devices (e.g. Blu-Ray).
	DecoderBufferSize    int32 `json:"decoder_buffer_size,omitempty"`    // Size of the decoder buffer, used in conjunction with bitrate_cap.
	OnePass              bool  `json:"one_pass,omitempty"`               // Force one-pass encoding.
	AudioConstantBitrate bool  `json:"audio_constant_bitrate,omitempty"` // Enable constant bitrate mode for audio if possible.

	// Frame Rate
	FrameRate              int32   `json:"frame_rate,omitempty"`               // The frame rate to use.
	MaxFrameRate           int32   `json:"max_frame_rate,omitempty"`           // The maximum frame rate to use.
	Decimate               int32   `json:"decimate,omitempty"`                 // Reduce the input bitrate by a divisor.
	KeyframeInterval       int32   `json:"keyframe_interval,omitempty"`        // The maximum number of frames between each keyframe.
	KeyframeRate           int32   `json:"keyframe_rate,omitempty"`            // The number of keyframes per second.
	FixedKeyframeInterval  bool    `json:"fixed_keyframe_interval,omitempty"`  // Enable fixed keyframe interval mode (VP6 and H.264 only).
	ForcedKeyframeInterval int32   `json:"forced_keyframe_interval,omitempty"` // Force keyframes at the specified interval (H.264 only).
	ForcedKeyframeRate     float64 `json:"forced_keyframe_rate,omitempty"`     // Specify the number of keyframes per-second, taking frame rate into account (H.264 only).

	// Audio
	AudioSampleRate    int32 `json:"audio_sample_rate,omitempty"`     // The audio sample rate, in Hz.
	MaxAudioSampleRate int32 `json:"max_audio_sample_rate,omitempty"` // The max audio sample rate, in Hz.
	AudioChannels      int32 `json:"audio_channels,omitempty"`        // The number of audio channels: 1 or 2.

	// Thumbnails
	Thumbnails []*ThumbnailSettings `json:"thumbnails,omitempty"` // Capture thumbnails for a given video.

	// Watermarks
	Watermarks []*WatermarkSettings `json:"watermarks,omitempty"` // Add one or more watermarks to an output video.

	// Captions
	CaptionUrl   string `json:"caption_url,omitempty"`   // URL to an SCC, DFXP, or SAMI caption file to include in the output.
	SkipCaptions bool   `json:"skip_captions,omitempty"` // Don't add or pass through captions to the output file.

	// Live Streaming
	LiveStream                bool  `json:"live_stream,omitempty"`                  // Create a live_stream job or output that is ready for playback within seconds.
	ReconnectTime             int32 `json:"reconnect_time,omitempty"`               // The time, in seconds, to wait for a stream to reconnect.
	EventLength               int32 `json:"event_length,omitempty"`                 // The minimum time, in seconds, to keep a live stream available.
	LiveSlidingWindowDuration int32 `json:"live_sliding_window_duration,omitempty"` // The time, in seconds, to keep in the HLS playlist.

	// Video Processing
	Rotate      int32  `json:"rotate,omitempty"`      // Rotate a video.
	Deinterlace string `json:"deinterlace,omitempty"` // Deinterlace input video.
	Sharpen     bool   `json:"sharpen,omitempty"`     // Apply a sharpen filter.
	Denoise     string `json:"denoise,omitempty"`     // Apply denoise filter.
	Autolevel   bool   `json:"autolevel,omitempty"`   // Apply a color auto-level filter.
	Deblock     bool   `json:"deblock,omitempty"`     // Apply deblock filter.

	// Audio Processing
	AudioGain                 int32   `json:"audio_gain,omitempty"`                  // Apply a gain amount to the audio, in dB.
	AudioNormalize            bool    `json:"audio_normalize,omitempty"`             // Normalize audio to 0dB.
	AudioPreNormalize         bool    `json:"audio_pre_normalize,omitempty"`         // Normalize the audio before applying expansion or compression effects.
	AudioPostNormalize        bool    `json:"audio_post_normalize,omitempty"`        // Normalize the audio after applying expansion or compression effects.
	AudioBass                 int32   `json:"audio_bass,omitempty"`                  // Increase or decrease the amount of bass in the audio.
	AudioTreble               int32   `json:"audio_treble,omitempty"`                // Increase or decrease the amount of treble in the audio.
	AudioHighpass             int32   `json:"audio_highpass,omitempty"`              // Apply a high-pass filter to the audio.
	AudioLowpass              int32   `json:"audio_lowpass,omitempty"`               // Apply a low-pass filter to the audio.
	AudioCompressionRatio     float64 `json:"audio_compression_ratio,omitempty"`     // Compress the dynamic range of the audio.
	AudioCompressionThreshold int32   `json:"audio_compression_threshold,omitempty"` // Compress the dynamic range of the audio.
	AudioExpansionRatio       float64 `json:"audio_expansion_ratio,omitempty"`       // Expand the dynamic range of the audio.
	AudioExpansionThreshold   int32   `json:"audio_expansion_threshold,omitempty"`   // Expand the dynamic range of the audio.
	AudioFade                 float64 `json:"audio_fade,omitempty"`                  // Apply fade-in and fade-out effects to the audio.
	AudioFadeIn               float64 `json:"audio_fade_in,omitempty"`               // Apply a fade-in effect to the audio.
	AudioFadeOut              float64 `json:"audio_fade_out,omitempty"`              // Apply a fade-out effect to the audio.
	AudioKaraokeMode          bool    `json:"audio_karaoke_mode,omitempty"`          // Apply a karaoke effect to the audio.

	// Clips
	StartClip  string `json:"start_clip,omitempty"`  // Encode only a portion of the input file by setting a custom start time.
	ClipLength string `json:"clip_length,omitempty"` // Encode only a portion of the input file by setting a custom clip length.

	// S3 Settings
	MakePublic    bool                     `json:"public,omitempty"`         // Make the output publicly readable on S3.
	UseRRS        bool                     `json:"rrs,omitempty"`            // Amazon S3's Reduced Redundancy Storage.
	AccessControl []*AccessControlSettings `json:"access_control,omitempty"` // Fine-grained access control rules for files sent to S3.

	// Notifications
	Notifications []*NotificationSettings `json:"notifications,omitempty"` // Be notified when a job or output is complete.

	// Conditional Outputs
	MinSize     string `json:"min_size,omitempty"`     // Skip output if the input file is smaller than the given dimensions.
	MaxSize     string `json:"max_size,omitempty"`     // Skip output if the input file is larger than the given dimensions.
	MinDuration int32  `json:"min_duration,omitempty"` // Skip output if the input file is shorter than the given duration, in seconds.
	MaxDuration int32  `json:"max_duration,omitempty"` // Skip output if the input file is longer than the given duration, in seconds.

	// Segmented Streaming
	SegmentSeconds        int32             `json:"segment_seconds,omitempty"`          // Sets the maximum duration of each segment a segmented output
	SegmentSize           int32             `json:"segment_size,omitempty"`             // Sets the maximum data size of each segment in a segmented output
	Streams               []*StreamSettings `json:"streams,omitempty"`                  // Provides a list of stream info to be reformatted as a playlist
	SegmentImageUrl       string            `json:"segment_image_url,omitempty"`        // An image to display on audio-only segments
	SegmentVideoSnapshots bool              `json:"segment_video_snapshots,omitempty"`  // When segmenting a video file into audio-only segments, take snapshots of the video as thumbnails for each segment.
	MaxHLSProtocolVersion int32             `json:"max_hls_protocol_version,omitempty"` // The maximum HLS protocol to use.
	HLSOptimizedTS        bool              `json:"hls_optimized_ts,omitempty"`         // Optimize TS segment files for HTTP Live Streaming on iOS.
	PrepareForSegmenting  string            `json:"prepare_for_segmenting,omitempty"`   // Include captions and keyframe timing for segmenting.
	InstantPlay           bool              `json:"instant_play,omitempty"`             // Create an instant play output that is ready for playback within seconds.
	SMILBaseUrl           string            `json:"smil_base_url,omitempty"`            // Add <meta base="smil_base_url_value"/> to the <head> section of an SMIL playlist.

	// Encryption
	EncryptionMethod            string `json:"encryption_method,omitempty"`              // Set the encryption method to use for encrypting.
	EncryptionKey               string `json:"encryption_key,omitempty"`                 // Set a single encryption key to use rather than having Zencoder generate one
	EncryptionKeyUrl            string `json:"encryption_key_url,omitempty"`             // Set a URL to a single encryption key to use rather than having Zencoder generate one
	EncryptionKeyRotationPeriod int32  `json:"encryption_key_rotation_period,omitempty"` // Rotate to a new encryption key after a number of segments
	EncryptionKeyUrlPrefix      string `json:"encryption_key_url_prefix,omitempty"`      // Prepend key URLs with the passed string
	EncryptionIV                string `json:"encryption_iv,omitempty"`                  // Set an initialization vector to use when encrypting
	EncryptionPassword          string `json:"encryption_password,omitempty"`            // Sets a password to use for generating an initialization vector

	// Decryption
	DecryptionMethod   string `json:"decryption_method,omitempty"`   // Set the decryption algorithm to use.
	DecryptionKey      string `json:"decryption_key,omitempty"`      // Set the decryption key to use.
	DecryptionKeyUrl   string `json:"decryption_key_url,omitempty"`  // The URL of a decryption key file to use.
	DecryptionPassword string `json:"decryption_password,omitempty"` // The password used in combination with the key to decrypt the input file.

	// H.264
	H264ReferenceFrames int32  `json:"h264_reference_frames,omitempty"` // A number of reference frames to use in H.264 video.
	H264Profile         string `json:"h264_profile,omitempty"`          // The H.264 profile to use.
	H264Level           string `json:"h264_level,omitempty"`            // The H.264 level to use.
	H264Bframes         int32  `json:"h264_bframes,omitempty"`          // The maximum number of consecutive B-frames.
	Tuning              string `json:"tuning,omitempty"`                // Tune the output video for a specific content type.
	Crf                 int32  `json:"crf,omitempty"`                   // Bitrate control setting.

	// FLV
	CuePoints []*CuePointSettings `json:"cue_points,omitempty"` // Add event or navigation cue points to a FLV video.

	// VP6
	VP6TemporalDownWatermark int32   `json:"vp6_temporal_down_watermark,omitempty"` // VP6 temporal down watermark percentage.
	VP6TemporalResampling    bool    `json:"vp6_temporal_resampling,omitempty"`     // Enable or disable VP6 temporal resampling.
	VP6UndershootPct         int32   `json:"vp6_undershoot_pct,omitempty"`          // Target a slightly lower datarate.
	VP6Profile               string  `json:"vp6_profile,omitempty"`                 // VP6 profile: vp6s or vp6e.
	VP6CompressionMode       string  `json:"vp6_compression_mode,omitempty"`        // VP6 compression mode: good or best.
	VP6TwoPassMinSection     int32   `json:"vp6_2pass_min_section,omitempty"`       // For two-pass VBR encoding, the lowest datarate that the encoder will allow.
	VP6TwoPassMaxSection     int32   `json:"vp6_2pass_max_section,omitempty"`       // For two-pass VBR encoding, the highest datarate that the encoder will allow.
	VP6StreamPrebuffer       int32   `json:"vp6_stream_prebuffer,omitempty"`        // Seconds of preload that are necessary before starting playback.
	VP6StreamMaxBuffer       int32   `json:"vp6_stream_max_buffer,omitempty"`       // Maximum decoder buffer size
	VP6DeinterlaceMode       string  `json:"vp6_deinterlace_mode,omitempty"`        // Deinterlace mode for VP6
	VP6DenoiseLevel          float64 `json:"vp6_denoise_level,omitempty"`           // Denoise level for VP6
	AlphaTransparency        bool    `json:"alpha_transparency,omitempty"`          // Enable alpha transparency. Currently, only supported by VP6.
	ConstantBitrate          bool    `json:"constant_bitrate,omitempty"`            // Use constant bitrate (CBR) encoding.

	// MP4
	Hint    bool  `json:"hint,omitempty"`     // Enable hinting of MP4 files for RTP/RTSP.
	MTUSize int32 `json:"mtu_size,omitempty"` // Set MTU size for MP4 hinting.

	// AAC
	MaxAacProfile   string `json:"max_aac_profile,omitempty"`   // What is the most advanced (compressed) AAC profile to allow?
	ForceAacProfile string `json:"force_aac_profile,omitempty"` // Force the use of a particular AAC profile, rather than letting Zencoder choose the best profile for the bitrate.

	// Aspera
	AsperaTransferPolicy string `json:"aspera_transfer_policy,omitempty"` // How to allocate available bandwidth for Aspera file transfers.
	TransferMinimumRate  int32  `json:"transfer_minimum_rate,omitempty"`  // A targeted rate in Kbps for data transfer minimums.
	TransferMaximumRate  int32  `json:"transfer_maximum_rate,omitempty"`  // A targeted rate in Kbps for data transfer maximums.

	// Transmuxing
	CopyVideo bool `json:"copy_video,omitempty"` // Copy the video track of the input file
	CopyAudio bool `json:"copy_audio,omitempty"` // Copy the audio track of the input file
}

type EncodingSettings struct {
	Input                string            `json:"input,omitempty"`                  // A S3, Cloud Files, GCS, FTP, FTPS, SFTP, or Aspera URL where we can download file to transcode.
	LiveStream           bool              `json:"live_stream,omitempty"`            // Create a Live streaming job.
	Outputs              []*OutputSettings `json:"outputs,omitempty"`                // An array or hash of output settings.
	Region               string            `json:"region,omitempty"`                 // The region where a file is processed: US, Europe, Asia, or Australia.
	Test                 bool              `json:"test,omitempty"`                   // Enable test mode ("Integration Mode") for a job.
	Private              bool              `json:"private,omitempty"`                // Enable privacy mode for a job.
	DownloadConnections  int32             `json:"download_connections,omitempty"`   // Utilize multiple, simultaneous connections for download acceleration (in some circumstances).
	PassThrough          string            `json:"pass_through,omitempty"`           // Optional information to store alongside this job.
	Mock                 bool              `json:"mock,omitempty"`                   // Send a mocked job request.
	Grouping             string            `json:"grouping,omitempty"`               // A report grouping for this job.
	AsperaTransferPolicy string            `json:"aspera_transfer_policy,omitempty"` // How to allocate available bandwidth for Aspera file transfers.
	TransferMinimumRate  int32             `json:"transfer_minimum_rate,omitempty"`  // A targeted rate in Kbps for data transfer minimums.
	TransferMaximumRate  int32             `json:"transfer_maximum_rate,omitempty"`  // A targeted rate in Kbps for data transfer maximums.
	ExpectedMD5Checksum  string            `json:"expected_md5_checksum,omitempty"`  // The expected checksum of the input file.
	Credentials          string            `json:"credentials,omitempty"`            // References saved credentials by a nickname.
}
