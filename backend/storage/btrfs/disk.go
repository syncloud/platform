package btrfs

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/prometheus/procfs/btrfs"
	"github.com/syncloud/platform/cli"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

const MKFS = "/snap/platform/current/btrfs/bin/mkfs.sh"
const BTRFS = "/snap/platform/current/btrfs/bin/btrfs.sh"

type Config interface {
	ExternalDiskDir() string
}

type DiskStats interface {
	Stats() ([]*btrfs.Stats, error)
}

type Disks struct {
	config   Config
	executor cli.CommandExecutor
	stats    DiskStats
	logger   *zap.Logger
}

func NewDisks(
	config Config,
	executor cli.CommandExecutor,
	stats DiskStats,
	logger *zap.Logger) *Disks {

	return &Disks{
		config:   config,
		executor: executor,
		stats:    stats,
		logger:   logger,
	}
}

func (d *Disks) Update(devices []string, existingUuid string) (string, error) {
	existingDevices, err := d.ExistingDevices(existingUuid)
	if err != nil {
		return "", err
	}
	newUuid := uuid.New().String()
	changes, err := d.Apply(existingDevices, devices, newUuid)
	if err != nil {
		return "", err
	}

	for _, change := range changes {
		output, err := d.executor.CommandOutput(change.Cmd, change.Args...)
		if err != nil {
			d.logger.Error(string(output))
			return "", err
		}
		d.logger.Info("btrfs", zap.String("output", string(output)))
	}

	return newUuid, nil
}

func (d *Disks) Apply(before []string, after []string, newUuid string) ([]Change, error) {
	removed := Diff(before, after)
	added := Diff(after, before)
	var changes []Change

	//if len(removed) == len(added) {
	//	for i, _ := range removed {
	//		changes = append(changes, NewChange("replace", removed[i], added[i]))
	//	}
	//	return changes, nil
	//}

	mode := "single"
	if len(after) == 2 {
		mode = "raid1"
	}
	if len(after) > 2 {
		mode = "raid10"
	}

	if len(before) == 0 {
		change := NewChange(MKFS, "-U", newUuid, "-f", "-m", mode, "-d", mode)
		change.Append(added...)
		changes = append(changes, change)
	} else {
		if len(added) > 0 {
			change := NewChange(BTRFS, "device", "add")
			change.Append(added...)
			change.Append(d.config.ExternalDiskDir())
			changes = append(changes, change)
			change = NewChange(BTRFS, "balance", "start", fmt.Sprintf("-dconvert=%s", mode), fmt.Sprintf("-mconvert=%s", mode), d.config.ExternalDiskDir())
			changes = append(changes, change)
		}
		if len(after) == 1 {
			change := NewChange(BTRFS, "balance", "start", fmt.Sprintf("-dconvert=%s", mode), fmt.Sprintf("-mconvert=%s", mode), d.config.ExternalDiskDir())
			changes = append(changes, change)
		}
	}

	if len(after) > 0 {
		if len(removed) > 0 {
			change := NewChange(BTRFS, "device", "delete")
			change.Append(removed...)
			change.Append(d.config.ExternalDiskDir())
			changes = append(changes, change)
		}
	}
	return changes, nil
}

func Diff(from []string, to []string) []string {
	var diff []string
	for _, v := range from {
		if !slices.Contains(to, v) {
			diff = append(diff, v)
		}
	}
	return diff
}

func (d *Disks) ExistingDevices(uuid string) ([]string, error) {
	stats, err := d.stats.Stats()
	if err != nil {
		return []string{}, err
	}

	var existing []string
	for _, fs := range stats {
		if fs.UUID == uuid {
			for device, _ := range fs.Devices {
				existing = append(existing, fmt.Sprintf("/dev/%s", device))
			}
		}
	}
	return existing, nil
}
