package storage

import "github.com/syncloud/platform/storage/model"

type BtrfsDiskStatsStub struct {
	profiles     map[string]model.RaidProfiles
	devices      []model.BtrfsDevice
	devicesCalls int
}

func (b *BtrfsDiskStatsStub) Profiles(uuid string) (model.RaidProfiles, error) {
	return b.profiles[uuid], nil
}

func (b *BtrfsDiskStatsStub) Devices() ([]model.BtrfsDevice, error) {
	b.devicesCalls++
	return b.devices, nil
}
