package btrfs

type DeviceStats struct {
	Header struct {
		Version string `json:"version"`
	} `json:"__header"`
	DeviceStats []struct {
		Device         string `json:"device"`
		Devid          string `json:"devid"`
		WriteIoErrs    string `json:"write_io_errs"`
		ReadIoErrs     string `json:"read_io_errs"`
		FlushIoErrs    string `json:"flush_io_errs"`
		CorruptionErrs string `json:"corruption_errs"`
		GenerationErrs string `json:"generation_errs"`
	} `json:"device-stats"`
}

func (d *DeviceStats) HasErrors(device string) bool {
	for _, stats := range d.DeviceStats {
		if stats.Device == device {
			return stats.WriteIoErrs != "0" || stats.ReadIoErrs != "0" || stats.FlushIoErrs != "0" || stats.CorruptionErrs != "0" || stats.GenerationErrs != "0"
		}
	}
	return false
}
