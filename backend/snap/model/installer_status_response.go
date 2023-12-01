package model

type InstallerStatus struct {
	IsRunning bool      `json:"is_running"`
	Progress  *Progress `json:"progress"`
}

type Progress struct {
	Summary    string `json:"summary"`
	Type       string `json:"type"`
	Percentage int    `json:"percentage"`
}
