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
	UserDomain string `json:"user_domain"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type Response struct {
	Success bool         `json:"success"`
	Message *string      `json:"message,omitempty"`
	Data    *interface{} `json:"data,omitempty"`
}
