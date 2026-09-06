package model

type ActivationStatusResponse struct {
	Activated bool   `json:"activated"`
	DeviceUrl string `json:"device_url"`
}
