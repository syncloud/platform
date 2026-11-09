package model

type DiskHealth struct {
	Device string       `json:"device"`
	Port   string       `json:"port"`
	Smart  SmartHealth  `json:"smart"`
	Link   KernelEvents `json:"link"`
	Io     KernelEvents `json:"io"`
}
