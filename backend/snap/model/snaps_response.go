package model

type SnapsResponse struct {
	Result []Snap `json:"result"`
	Status string `json:"status"`
}
