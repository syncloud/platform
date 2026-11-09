package btrfs

import (
	"encoding/json"
	"github.com/prometheus/procfs/btrfs"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/storage/model"
	"sort"
	"strings"
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
	profiles, err := s.Profiles(uuid)
	if err != nil {
		return "", err
	}
	return profiles.Data, nil
}

func (s *Stats) Profiles(uuid string) (model.RaidProfiles, error) {
	stats, err := s.Info()
	if err != nil {
		return model.RaidProfiles{}, err
	}
	for _, fs := range stats {
		if fs.UUID == uuid {
			return model.RaidProfiles{
				Data:     layouts(fs.Allocation.Data),
				Metadata: layouts(fs.Allocation.Metadata),
			}, nil
		}
	}
	return model.RaidProfiles{}, nil
}

func layouts(allocation *btrfs.AllocationStats) string {
	if allocation == nil {
		return ""
	}
	var names []string
	for name := range allocation.Layouts {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}

func (s *Stats) HasErrors(device string) (bool, error) {
	result, err := s.deviceStats()
	if err != nil {
		return false, err
	}
	return result.HasErrors(device), nil
}

func (s *Stats) Devices() ([]model.BtrfsDevice, error) {
	result, err := s.deviceStats()
	if err != nil {
		return nil, err
	}
	return result.Devices(), nil
}

func (s *Stats) deviceStats() (*DeviceStats, error) {
	output, err := s.executor.CombinedOutput(BTRFS, "--format", "json", "device", "stats", s.config.ExternalDiskDir())
	if err != nil {
		return nil, err
	}
	var result DeviceStats
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Stats) ResetErrors() error {
	_, err := s.executor.CombinedOutput(BTRFS, "device", "stats", "-z", s.config.ExternalDiskDir())
	return err
}

func (s *Stats) Mounted() bool {
	output, err := s.executor.CombinedOutput("stat", "-f", "-c", "%T", s.config.ExternalDiskDir())
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == Module
}
