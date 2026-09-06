package model

type StorageActivateDisksRequest struct {
	Devices []string `json:"devices"`
	Format  bool     `json:"format"`
}
