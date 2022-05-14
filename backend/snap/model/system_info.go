package model

type SystemInfo struct {
	Result Result `json:"result"`
}

type Result struct {
	Version string `json:"version"`
}
