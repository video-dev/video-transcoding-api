package elementalconductor

import (
	"encoding/xml"
	"net/http"

	"gopkg.in/check.v1"
)

func (s *S) TestGetPresets(c *check.C) {
	presetsResponseXML := `<?xml version="1.0" encoding="UTF-8"?>
<preset_list>
  <preset href="/presets/1" product="Elemental Conductor File + Audio Normalization Package + Audio Package" version="2.7.2vd.32545">
    <name>iPhone</name>
    <permalink>iphone</permalink>
    <description>Default output for iPhone</description>
    <preset_category href="/preset_categories/6">Devices</preset_category>
  </preset>
  <preset href="/presets/2" product="Elemental Conductor File + Audio Normalization Package + Audio Package" version="2.7.2vd.32545">
    <name>iPhone_ADAPT_HIGH</name>
    <permalink>iphone_adapt_high</permalink>
    <description>Default output for iPhone Adaptive high quality</description>
    <preset_category href="/preset_categories/6">Devices</preset_category>
  </preset>
  <next href="https://3e9n5rjaf3eb2.cloud.elementaltechnologies.com/presets?page=2&amp;amp;per_page=30"/>
</preset_list>`

	expectedPreset1 := Preset{
		XMLName:     xml.Name{Local: "preset"},
		Name:        "iPhone",
		Href:        "/presets/1",
		Permalink:   "iphone",
		Description: "Default output for iPhone",
	}

	expectedPreset2 := Preset{
		XMLName:     xml.Name{Local: "preset"},
		Name:        "iPhone_ADAPT_HIGH",
		Href:        "/presets/2",
		Permalink:   "iphone_adapt_high",
		Description: "Default output for iPhone Adaptive high quality",
	}

	var expectedOutput PresetList
	expectedOutput.Presets = make([]Preset, 2)
	expectedOutput.Presets[0] = expectedPreset1
	expectedOutput.Presets[1] = expectedPreset2

	server, _ := s.startServer(http.StatusOK, presetsResponseXML)
	defer server.Close()

	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	getPresetsResponse, _ := client.GetPresets()
	c.Assert(getPresetsResponse, check.DeepEquals, &expectedOutput)
}

func (s *S) TestGetPreset(c *check.C) {
	presetResponseXML := `<?xml version="1.0" encoding="UTF-8"?>
<preset href="/presets/1" product="Elemental Conductor File + Audio Normalization Package + Audio Package" version="2.7.2vd.32545">
  <name>iPhone</name>
  <permalink>iphone</permalink>
  <description>Default output for iPhone</description>
  <preset_category href="/preset_categories/6">Devices</preset_category>
  <container>mp4</container>
  <mp4_settings>
    <id>1</id>
    <include_cslg>false</include_cslg>
    <mp4_major_brand nil="true"/>
    <progressive_downloading>true</progressive_downloading>
  </mp4_settings>
  <log_edit_points>false</log_edit_points>
  <video_description>
    <afd_signaling>None</afd_signaling>
    <anti_alias>true</anti_alias>
    <drop_frame_timecode>true</drop_frame_timecode>
    <encoder_type nil="true"/>
    <fixed_afd nil="true"/>
    <force_cpu_encode>false</force_cpu_encode>
    <height>320</height>
    <id>1</id>
    <insert_color_metadata>false</insert_color_metadata>
    <respond_to_afd>None</respond_to_afd>
    <sharpness>50</sharpness>
    <stretch_to_output>false</stretch_to_output>
    <timecode_passthrough>false</timecode_passthrough>
    <vbi_passthrough>false</vbi_passthrough>
    <width>480</width>
    <h264_settings>
      <adaptive_quantization>medium</adaptive_quantization>
      <bitrate>960000</bitrate>
      <buf_fill_pct nil="true"/>
      <buf_size nil="true"/>
      <cabac>false</cabac>
      <flicker_reduction>off</flicker_reduction>
      <force_field_pictures>false</force_field_pictures>
      <framerate_denominator>1</framerate_denominator>
      <framerate_follow_source>false</framerate_follow_source>
      <framerate_numerator>24</framerate_numerator>
      <gop_b_reference>false</gop_b_reference>
      <gop_closed_cadence>1</gop_closed_cadence>
      <gop_markers>false</gop_markers>
      <gop_num_b_frames>0</gop_num_b_frames>
      <gop_size>80</gop_size>
      <id>1</id>
      <interpolate_frc>false</interpolate_frc>
      <look_ahead_rate_control>medium</look_ahead_rate_control>
      <max_bitrate nil="true"/>
      <max_qp nil="true"/>
      <min_i_interval>0</min_i_interval>
      <min_qp nil="true"/>
      <num_ref_frames>1</num_ref_frames>
      <par_denominator>1</par_denominator>
      <par_follow_source>false</par_follow_source>
      <par_numerator>1</par_numerator>
      <passes>1</passes>
      <qp nil="true"/>
      <qp_step nil="true"/>
      <repeat_pps>false</repeat_pps>
      <scd>true</scd>
      <sei_timecode>false</sei_timecode>
      <slices>1</slices>
      <slow_pal>false</slow_pal>
      <softness nil="true"/>
      <svq>0</svq>
      <telecine>None</telecine>
      <transition_detection>false</transition_detection>
      <level>3</level>
      <profile>Baseline</profile>
      <rate_control_mode>ABR</rate_control_mode>
      <gop_mode>fixed</gop_mode>
      <interlace_mode>progressive</interlace_mode>
    </h264_settings>
    <gpu/>
    <selected_gpu nil="true"/>
    <codec>h.264</codec>
    <video_preprocessors>
      <deinterlacer>
        <algorithm>interpolate</algorithm>
        <deinterlace_mode>Deinterlace</deinterlace_mode>
        <force>false</force>
        <id>85</id>
      </deinterlacer>
    </video_preprocessors>
  </video_description>
  <audio_description>
    <audio_type>0</audio_type>
    <follow_input_audio_type>false</follow_input_audio_type>
    <follow_input_language_code>false</follow_input_language_code>
    <id>1</id>
    <language_code nil="true"/>
    <order>1</order>
    <stream_name nil="true"/>
    <aac_settings>
      <bitrate>128000</bitrate>
      <coding_mode>2_0</coding_mode>
      <id>1</id>
      <latm_loas>false</latm_loas>
      <mpeg2>false</mpeg2>
      <sample_rate>44100</sample_rate>
      <profile>LC</profile>
      <rate_control_mode>CBR</rate_control_mode>
    </aac_settings>
    <codec>aac</codec>
  </audio_description>
</preset>`

	expectedPreset := Preset{
		XMLName:       xml.Name{Local: "preset"},
		Name:          "iPhone",
		Href:          "/presets/1",
		Permalink:     "iphone",
		Description:   "Default output for iPhone",
		Container:     "mp4",
		VideoCodec:    "h.264",
		AudioCodec:    "aac",
		Width:         "480",
		Height:        "320",
		VideoBitrate:  "960000",
		AudioBitrate:  "128000",
		GopSize:       "80",
		GopMode:       "fixed",
		Profile:       "Baseline",
		ProfileLevel:  "3",
		RateControl:   "ABR",
		InterlaceMode: "progressive",
	}

	server, _ := s.startServer(http.StatusOK, presetResponseXML)
	defer server.Close()

	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	getPresetResponse, _ := client.GetPreset("1")
	c.Assert(getPresetResponse, check.DeepEquals, &expectedPreset)
}

func (s *S) TestGetPresetForHls(c *check.C) {
	presetHLSResponseXML := `<preset href="/presets/149" product="Elemental Conductor File + Audio Normalization Package + Audio Package" version="2.7.2vd.32545">
  <name>nyt_hls_720p_high_uhd</name>
  <permalink>nyt_hls_720p_high_uhd</permalink>
  <description></description>
  <preset_category href="/preset_categories/6">Devices</preset_category>
  <container>m3u8</container>
  <m3u8_settings>
    <audio_packets_per_pes>16</audio_packets_per_pes>
    <id>200</id>
    <pat_interval>0</pat_interval>
    <pcr_every_pes>true</pcr_every_pes>
    <pcr_period nil="true"/>
    <pmt_interval>0</pmt_interval>
    <program_num>1</program_num>
    <transport_stream_id nil="true"/>
    <audio_pids>482-498</audio_pids>
    <pmt_pid>480</pmt_pid>
    <private_metadata_pid>503</private_metadata_pid>
    <scte35_pid>500</scte35_pid>
    <timed_metadata_pid>502</timed_metadata_pid>
    <video_pid>481</video_pid>
    <pcr_pid>481</pcr_pid>
  </m3u8_settings>
  <log_edit_points>false</log_edit_points>
  <video_description>
    <afd_signaling>None</afd_signaling>
    <anti_alias>true</anti_alias>
    <drop_frame_timecode>true</drop_frame_timecode>
    <encoder_type nil="true"/>
    <fixed_afd nil="true"/>
    <force_cpu_encode>false</force_cpu_encode>
    <height>720</height>
    <id>501</id>
    <insert_color_metadata>false</insert_color_metadata>
    <respond_to_afd>None</respond_to_afd>
    <sharpness>50</sharpness>
    <stretch_to_output>false</stretch_to_output>
    <timecode_passthrough>false</timecode_passthrough>
    <vbi_passthrough>false</vbi_passthrough>
    <width nil="true"/>
    <h264_settings>
      <adaptive_quantization>medium</adaptive_quantization>
      <bitrate>3800000</bitrate>
      <buf_fill_pct nil="true"/>
      <buf_size>7600000</buf_size>
      <cabac>false</cabac>
      <flicker_reduction>off</flicker_reduction>
      <force_field_pictures>false</force_field_pictures>
      <framerate_denominator nil="true"/>
      <framerate_follow_source>true</framerate_follow_source>
      <framerate_numerator nil="true"/>
      <gop_b_reference>false</gop_b_reference>
      <gop_closed_cadence>1</gop_closed_cadence>
      <gop_markers>false</gop_markers>
      <gop_num_b_frames>2</gop_num_b_frames>
      <gop_size>90</gop_size>
      <id>439</id>
      <interpolate_frc>false</interpolate_frc>
      <look_ahead_rate_control>high</look_ahead_rate_control>
      <max_bitrate>4750000</max_bitrate>
      <max_qp nil="true"/>
      <min_i_interval>0</min_i_interval>
      <min_qp nil="true"/>
      <num_ref_frames>1</num_ref_frames>
      <par_denominator>1</par_denominator>
      <par_follow_source>false</par_follow_source>
      <par_numerator>1</par_numerator>
      <passes>2</passes>
      <qp nil="true"/>
      <qp_step nil="true"/>
      <repeat_pps>false</repeat_pps>
      <scd>true</scd>
      <sei_timecode>false</sei_timecode>
      <slices>1</slices>
      <slow_pal>false</slow_pal>
      <softness nil="true"/>
      <svq>0</svq>
      <telecine>None</telecine>
      <transition_detection>false</transition_detection>
      <level>3.1</level>
      <profile>Main</profile>
      <rate_control_mode>VBR</rate_control_mode>
      <gop_mode>fixed</gop_mode>
      <interlace_mode>progressive</interlace_mode>
    </h264_settings>
    <gpu/>
    <selected_gpu nil="true"/>
    <codec>h.264</codec>
    <video_preprocessors>
      <deinterlacer>
        <algorithm>interpolate</algorithm>
        <deinterlace_mode>Deinterlace</deinterlace_mode>
        <force>false</force>
        <id>376</id>
      </deinterlacer>
    </video_preprocessors>
  </video_description>
  <audio_description>
    <audio_type>0</audio_type>
    <follow_input_audio_type>false</follow_input_audio_type>
    <follow_input_language_code>false</follow_input_language_code>
    <id>516</id>
    <language_code nil="true"/>
    <order>1</order>
    <stream_name nil="true"/>
    <aac_settings>
      <bitrate>64000</bitrate>
      <coding_mode>2_0</coding_mode>
      <id>435</id>
      <latm_loas>false</latm_loas>
      <mpeg2>false</mpeg2>
      <sample_rate>44100</sample_rate>
      <profile>LC</profile>
      <rate_control_mode>CBR</rate_control_mode>
    </aac_settings>
    <codec>aac</codec>
  </audio_description>
</preset>`

	expectedPreset := Preset{
		XMLName:       xml.Name{Local: "preset"},
		Name:          "nyt_hls_720p_high_uhd",
		Href:          "/presets/149",
		Permalink:     "nyt_hls_720p_high_uhd",
		Container:     "m3u8",
		VideoCodec:    "h.264",
		AudioCodec:    "aac",
		Height:        "720",
		VideoBitrate:  "3800000",
		AudioBitrate:  "64000",
		GopSize:       "90",
		GopMode:       "fixed",
		Profile:       "Main",
		ProfileLevel:  "3.1",
		RateControl:   "VBR",
		InterlaceMode: "progressive",
	}

	server, _ := s.startServer(http.StatusOK, presetHLSResponseXML)
	defer server.Close()

	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	getPresetResponse, _ := client.GetPreset("149")
	c.Assert(getPresetResponse, check.DeepEquals, &expectedPreset)
}

func (s *S) TestCreatePreset(c *check.C) {
	createPresetResponseXML := `<?xml version="1.0" encoding="UTF-8"?>
<preset product="Elemental Conductor File + Audio Normalization Package + Audio Package" version="2.7.2vd.32545">
  <name>TestPresetName</name>
  <permalink></permalink>
  <description>Preset test here</description>
  <container>mp4</container>
  <mp4_settings>
    <id>163</id>
    <include_cslg>false</include_cslg>
    <mp4_major_brand nil="true"/>
    <progressive_downloading>false</progressive_downloading>
  </mp4_settings>
  <log_edit_points>false</log_edit_points>
  <video_description>
    <afd_signaling>None</afd_signaling>
    <anti_alias>true</anti_alias>
    <drop_frame_timecode>true</drop_frame_timecode>
    <encoder_type nil="true"/>
    <fixed_afd nil="true"/>
    <force_cpu_encode>false</force_cpu_encode>
    <height>720</height>
    <id>602</id>
    <insert_color_metadata>false</insert_color_metadata>
    <respond_to_afd>None</respond_to_afd>
    <sharpness>50</sharpness>
    <stretch_to_output>false</stretch_to_output>
    <timecode_passthrough>false</timecode_passthrough>
    <vbi_passthrough>false</vbi_passthrough>
    <width nil="true"/>
    <h264_settings>
      <adaptive_quantization>medium</adaptive_quantization>
      <bitrate>3800000</bitrate>
      <buf_fill_pct nil="true"/>
      <buf_size nil="true"/>
      <cabac>true</cabac>
      <flicker_reduction>off</flicker_reduction>
      <force_field_pictures>false</force_field_pictures>
      <framerate_denominator nil="true"/>
      <framerate_follow_source>true</framerate_follow_source>
      <framerate_numerator nil="true"/>
      <gop_b_reference>false</gop_b_reference>
      <gop_closed_cadence>1</gop_closed_cadence>
      <gop_markers>false</gop_markers>
      <gop_num_b_frames>2</gop_num_b_frames>
      <gop_size>90</gop_size>
      <id>540</id>
      <interpolate_frc>false</interpolate_frc>
      <look_ahead_rate_control>medium</look_ahead_rate_control>
      <max_bitrate nil="true"/>
      <max_qp nil="true"/>
      <min_i_interval>0</min_i_interval>
      <min_qp nil="true"/>
      <num_ref_frames>1</num_ref_frames>
      <par_denominator nil="true"/>
      <par_follow_source>true</par_follow_source>
      <par_numerator nil="true"/>
      <passes>1</passes>
      <qp nil="true"/>
      <qp_step nil="true"/>
      <repeat_pps>false</repeat_pps>
      <scd>true</scd>
      <sei_timecode>false</sei_timecode>
      <slices>1</slices>
      <slow_pal>false</slow_pal>
      <softness nil="true"/>
      <svq>0</svq>
      <telecine>None</telecine>
      <transition_detection>false</transition_detection>
      <level>3.1</level>
      <profile>Main</profile>
      <rate_control_mode>VBR</rate_control_mode>
      <gop_mode>fixed</gop_mode>
      <interlace_mode>progressive</interlace_mode>
    </h264_settings>
    <gpu/>
    <selected_gpu nil="true"/>
    <codec>h.264</codec>
  </video_description>
  <audio_description>
    <audio_type>0</audio_type>
    <follow_input_audio_type>false</follow_input_audio_type>
    <follow_input_language_code>false</follow_input_language_code>
    <id>617</id>
    <language_code nil="true"/>
    <order>1</order>
    <stream_name nil="true"/>
    <aac_settings>
      <bitrate>64000</bitrate>
      <coding_mode>2_0</coding_mode>
      <id>536</id>
      <latm_loas>false</latm_loas>
      <mpeg2>false</mpeg2>
      <sample_rate>48000</sample_rate>
      <profile>LC</profile>
      <rate_control_mode>CBR</rate_control_mode>
    </aac_settings>
    <codec>aac</codec>
  </audio_description>
</preset>`
	server, _ := s.startServer(http.StatusOK, createPresetResponseXML)
	defer server.Close()

	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	preset := Preset{
		XMLName:       xml.Name{Local: "preset"},
		Name:          "TestPresetName",
		Description:   "Preset test here",
		Container:     "mp4",
		VideoCodec:    "h.264",
		AudioCodec:    "aac",
		Height:        "720",
		VideoBitrate:  "3800000",
		AudioBitrate:  "64000",
		GopSize:       "90",
		GopMode:       "fixed",
		Profile:       "Main",
		ProfileLevel:  "3.1",
		RateControl:   "VBR",
		InterlaceMode: "progressive",
	}

	res, _ := client.CreatePreset(&preset)
	c.Assert(res, check.DeepEquals, &preset)
}

func (s *S) TestDeletePreset(c *check.C) {
	presetsResponse := ` `
	server, _ := s.startServer(http.StatusOK, presetsResponse)
	defer server.Close()

	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	deletePresetResponse := client.DeletePreset("preset123")
	c.Assert(deletePresetResponse, check.DeepEquals, nil)
}
