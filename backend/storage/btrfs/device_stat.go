package btrfs

import "github.com/syncloud/platform/storage/model"

type DeviceStat struct {
	Device         string  `json:"device"`
	Devid          Counter `json:"devid"`
	WriteIoErrs    Counter `json:"write_io_errs"`
	ReadIoErrs     Counter `json:"read_io_errs"`
	FlushIoErrs    Counter `json:"flush_io_errs"`
	CorruptionErrs Counter `json:"corruption_errs"`
	GenerationErrs Counter `json:"generation_errs"`
}

func (s DeviceStat) Errors() model.DeviceErrors {
	return model.DeviceErrors{
		Write:      s.WriteIoErrs.Int64(),
		Read:       s.ReadIoErrs.Int64(),
		Flush:      s.FlushIoErrs.Int64(),
		Corruption: s.CorruptionErrs.Int64(),
		Generation: s.GenerationErrs.Int64(),
	}
}
