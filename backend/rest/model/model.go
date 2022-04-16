package model

type Access struct {
	Ipv4        *string `json:"ipv4,omitempty"`
	Ipv4Enabled bool    `json:"ipv4_enabled"`
	Ipv4Public  bool    `json:"ipv4_public"`
	AccessPort  *int    `json:"access_port"`
	Ipv6Enabled bool    `json:"ipv6_enabled"`
}

type RedirectInfoResponse struct {
	Domain string `json:"domain"`
}

type BackupCreateRequest struct {
	App string `json:"app"`
}

type BackupRestoreRequest struct {
	File string `json:"file"`
}

type BackupRemoveRequest struct {
	File string `json:"file"`
}

type StorageFormatRequest struct {
	Device string `json:"device"`
}

type EventTriggerRequest struct {
	Event string `json:"event"`
}

type Response struct {
	Success            bool                 `json:"success"`
	Message            string               `json:"message,omitempty"`
	Data               *interface{}         `json:"data,omitempty"`
	ParametersMessages *[]ParameterMessages `json:"parameters_messages,omitempty"`
}
