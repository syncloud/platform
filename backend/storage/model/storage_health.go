package model

type StorageHealth struct {
	Disks    []DiskHealth  `json:"disks"`
	Scrub    *ScrubStatus  `json:"scrub"`
	Schedule CheckSchedule `json:"schedule"`
	Busy     bool          `json:"busy"`
}
