package hybrik

// CreateJob .
type CreateJob struct {
	Name              string           `json:"name"`
	Payload           CreateJobPayload `json:"payload"`
	Schema            string           `json:"schema,omitempty"`
	Expiration        int              `json:"expiration,omitempty"`
	Priority          int              `json:"priority,omitempty"`
	TaskRetryCount    int              `json:"task_retry:count,omitempty"`
	TaskRetryDelaySec int              `json:"task_retry:delay_sec,omitempty"`
	TaskTags          []string         `json:"task_tags,omitempty"`
	UserTag           string           `json:"user_tag,omitempty"`
}

// CreateJobPayload .
type CreateJobPayload struct {
	Elements    []Element    `json:"elements,omitempty"`
	Connections []Connection `json:"connections,omitempty"`
}

// Element .
type Element struct {
	UID     string              `json:"uid"`
	Kind    string              `json:"kind"`
	Task    *ElementTaskOptions `json:"task,omitempty"`
	Preset  *TranscodePreset    `json:"preset,omitempty"`
	Payload interface{}         `json:"payload"` // Can be of type ElementPayload or LocationTargetPayload
}

// ElementTaskOptions .
type ElementTaskOptions struct {
	Name string `json:"name"`
}

// ElementPayload .
type ElementPayload struct {
	Kind    string       `json:"kind,omitempty"`
	Payload AssetPayload `json:"payload"`
}

// ManifestCreatorPayload .
type ManifestCreatorPayload struct {
	Location    TranscodeLocation `json:"location"`
	FilePattern string            `json:"file_pattern"`
	Kind        string            `json:"kind"`
	UID         string            `json:"uid,omitempty"`
}

// LocationTargetPayload .
type LocationTargetPayload struct {
	Location TranscodeLocation         `json:"location"`
	Targets  []TranscodeLocationTarget `json:"targets"`
}

// TranscodePreset .
type TranscodePreset struct {
	Key string `json:"key"`
}

// TranscodeLocationTarget .
type TranscodeLocationTarget struct {
	FilePattern   string                   `json:"file_pattern"`
	ExistingFiles string                   `json:"existing_files,omitempty"`
	Container     TranscodeTargetContainer `json:"container,omitempty"`
	Location      *TranscodeLocation       `json:"location,omitempty"`
}

// TranscodeTargetContainer .
type TranscodeTargetContainer struct {
	SegmentDuration uint `json:"segment_duration,omitempty"`
}

// AssetPayload .
type AssetPayload struct {
	StorageProvider string `json:"storage_provider,omitempty"`

	URL string `json:"url,omitempty"`
}

// TranscodeLocation .
type TranscodeLocation struct {
	StorageProvider string `json:"storage_provider,omitempty"`
	Path            string `json:"path,omitempty"`
}

//TranscodeTarget .
type TranscodeTarget struct {
	FilePattern   string                   `json:"file_pattern"`
	ExistingFiles string                   `json:"existing_files"`
	Container     TranscodeContainer       `json:"container"`
	Video         map[string]interface{}   `json:"video"`
	Audio         []map[string]interface{} `json:"audio"`
}

// TranscodeContainer .
type TranscodeContainer struct {
	Kind string `json:"kind"`
}

// Connection .
type Connection struct {
	From []ConnectionFrom `json:"from,omitempty"`
	To   ConnectionTo     `json:"to,omitempty"`
}

// ConnectionFrom .
type ConnectionFrom struct {
	Element string `json:"element,omitempty"`
}

// ConnectionTo .
type ConnectionTo struct {
	Success []ToSuccess `json:"success,omitempty"`
	Error   []ToError   `json:"error,omitempty"`
}

// ToSuccess .
type ToSuccess struct {
	Element string `json:"element,omitempty"`
}

// ToError .
type ToError struct {
	Element string `json:"element,omitempty"`
}
