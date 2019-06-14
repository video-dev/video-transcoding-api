package container

import "github.com/NYTimes/video-transcoding-api/provider"

// StatusEnricher enriches status information for output containers
type StatusEnricher interface {
	Enrich(provider.JobStatus) (provider.JobStatus, error)
}
