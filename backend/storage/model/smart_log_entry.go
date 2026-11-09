package model

type SmartLogEntry struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	Passed *bool  `json:"passed,omitempty"`
	Hours  int64  `json:"hours"`
}
