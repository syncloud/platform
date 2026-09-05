package model

type StorageActivatePartitionRequest struct {
	Device string `json:"device"`
	Format bool   `json:"format"`
}
