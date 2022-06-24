package btrfs

import (
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
