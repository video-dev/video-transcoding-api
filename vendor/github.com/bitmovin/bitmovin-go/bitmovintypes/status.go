package bitmovintypes

type Status string

const (
	StatusCreated  Status = "CREATED"
	StatusQueued   Status = "QUEUED"
	StatusRunning  Status = "RUNNING"
	StatusFinished Status = "FINISHED"
	StatusError    Status = "ERROR"
)
