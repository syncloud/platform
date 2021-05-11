package model

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

type RedirectCheckFreeDomainRequest struct {
	Email    string `json:"redirect_email"`
	Password string `json:"redirect_password"`
	Domain   string `json:"user_domain"`
}

type Response struct {
	Success bool         `json:"success"`
	Message *string      `json:"message,omitempty"`
	Data    *interface{} `json:"data,omitempty"`
}
