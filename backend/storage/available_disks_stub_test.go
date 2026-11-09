package storage

import "github.com/syncloud/platform/storage/model"

type AvailableDisksStub struct {
	disks []model.Disk
}

func (d *AvailableDisksStub) AvailableDisks() ([]model.Disk, error) {
	return d.disks, nil
}
