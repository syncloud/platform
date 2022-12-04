package model

type SnapResponse struct {
	Result Snap   `json:"result"`
	Status string `json:"status"`
}
