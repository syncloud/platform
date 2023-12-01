package model

type InstallerStatus struct {
	IsRunning bool               `json:"is_running"`
	Progress  *InstallerProgress `json:"progress"`
}

type InstallerProgress struct {
	Summary       string `json:"summary"`
	Indeterminate bool   `json:"indeterminate"`
	Percentage    int64  `json:"percentage"`
}
