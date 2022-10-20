package btrfs

import (
	"fmt"
	"github.com/prometheus/procfs/btrfs"
)

type Stats struct {
}

func New() *Stats {
	return &Stats{}
}

func (s *Stats) Stats() ([]*btrfs.Stats, error) {
	fs, err := btrfs.NewDefaultFS()
	if err != nil {
		return nil, err
	}
	return fs.Stats()
}

func (s *Stats) ExistingDevices(uuid string) ([]string, error) {
	stats, err := s.Stats()
	if err != nil {
		return []string{}, err
	}

	var existing []string
	for _, fs := range stats {
		if fs.UUID == uuid {
			for device := range fs.Devices {
				existing = append(existing, fmt.Sprintf("/dev/%s", device))
			}
		}
	}
	return existing, nil
}
