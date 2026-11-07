package model

const (
	ScrubNever       = "never"
	ScrubRunning     = "running"
	ScrubFinished    = "finished"
	ScrubAborted     = "aborted"
	ScrubInterrupted = "interrupted"
)

type ScrubStatus struct {
	State         string `json:"state"`
	Started       string `json:"started,omitempty"`
	Duration      string `json:"duration,omitempty"`
	Progress      string `json:"progress,omitempty"`
	Errors        string `json:"errors,omitempty"`
	Corrected     int64  `json:"corrected"`
	Uncorrectable int64  `json:"uncorrectable"`
}
