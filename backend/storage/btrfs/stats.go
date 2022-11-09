package btrfs

import (
	"encoding/json"
	"github.com/prometheus/procfs/btrfs"
	"github.com/syncloud/platform/cli"
)

type Stats struct {
	config   Config
	executor cli.Executor
}

func NewStats(config Config, executor cli.Executor) *Stats {
	return &Stats{
		config:   config,
		executor: executor,
	}
}

func (s *Stats) Info() ([]*btrfs.Stats, error) {
	fs, err := btrfs.NewDefaultFS()
	if err != nil {
		return nil, err
	}
	return fs.Stats()
}

func (s *Stats) RaidMode(uuid string) (string, error) {
	stats, err := s.Info()
	if err != nil {
		return "", err
	}

	for _, fs := range stats {
		if fs.UUID == uuid {
			for raid := range fs.Allocation.Data.Layouts {
				return raid, nil
			}
		}
	}
	return "", nil
}

func (s *Stats) HasErrors(device string) (bool, error) {
	output, err := s.executor.CombinedOutput(BTRFS, "--format", "json", "device", "stats", s.config.ExternalDiskDir())
	if err != nil {
		return false, err
	}

	var result DeviceStats
	if err := json.Unmarshal(output, &result); err != nil {
		return false, err
	}
	return result.HasErrors(device), nil
}
