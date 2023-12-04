package model

type InstallerStatus struct {
	IsRunning bool                         `json:"is_running"`
	Progress  map[string]InstallerProgress `json:"progress"`
}

type InstallerProgress struct {
	App           string `json:"app"`
	Summary       string `json:"summary"`
	Indeterminate bool   `json:"indeterminate"`
	Percentage    int64  `json:"percentage"`
}
