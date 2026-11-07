package model

type KernelEvents struct {
	Count int    `json:"count"`
	Last  string `json:"last,omitempty"`
}
