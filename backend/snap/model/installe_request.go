package model

type InstallRequest struct {
	Action string `json:"action"`
	Purge  bool   `json:"purge,omitempty"`
}
