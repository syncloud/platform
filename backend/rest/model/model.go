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

type Response struct {
	Success bool         `json:"success"`
	Message *string      `json:"message,omitempty"`
	Data    *interface{} `json:"data,omitempty"`
}
