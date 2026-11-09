package btrfs

import "github.com/syncloud/platform/storage/model"

type DeviceStats struct {
	Header      DeviceStatsHeader `json:"__header"`
	DeviceStats []DeviceStat      `json:"device-stats"`
}

func (d *DeviceStats) HasErrors(device string) bool {
	for _, stats := range d.DeviceStats {
		if stats.Device == device {
			return stats.Errors().Any()
		}
	}
	return false
}

func (d *DeviceStats) Devices() []model.BtrfsDevice {
	var devices []model.BtrfsDevice
	for _, stats := range d.DeviceStats {
		devices = append(devices, model.BtrfsDevice{
			Device: stats.Device,
			Devid:  string(stats.Devid),
			Errors: stats.Errors(),
		})
	}
	return devices
}
