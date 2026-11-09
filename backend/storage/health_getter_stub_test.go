package storage

import "github.com/syncloud/platform/storage/model"

type HealthGetterStub struct {
	health model.StorageHealth
}

func (h *HealthGetterStub) Get() (*model.StorageHealth, error) {
	return &h.health, nil
}
