package model

type SmartTestResult struct {
	Running          bool   `json:"running"`
	RemainingPercent int    `json:"remaining_percent"`
	LastType         string `json:"last_type,omitempty"`
	LastStatus       string `json:"last_status,omitempty"`
	LastPassed       *bool  `json:"last_passed,omitempty"`
	LastHours        *int64 `json:"last_hours,omitempty"`
}
