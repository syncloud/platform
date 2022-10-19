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

type Systemd interface {
	AddMount(device string) error
}

type Disks struct {
	config   Config
	executor cli.CommandExecutor
	stats    DiskStats
	systemd  Systemd
	logger   *zap.Logger
}

func NewDisks(
	config Config,
	executor cli.CommandExecutor,
	stats DiskStats,
	systemd Systemd,
	logger *zap.Logger) *Disks {

	return &Disks{
		config:   config,
		executor: executor,
		stats:    stats,
		systemd:  systemd,
		logger:   logger,
	}
}

func (d *Disks) Update(devices []string, existingUuid string, format bool) (string, error) {
	existingDevices, err := d.ExistingDevices(existingUuid)
	if err != nil {
		return "", err
	}
	newUuid := uuid.New().String()
	err = d.Apply(existingDevices, devices, newUuid, format)
	if err != nil {
		return "", err
	}
	return newUuid, nil
}

func (d *Disks) Apply(before []string, after []string, newUuid string, format bool) error {
	removed := Diff(before, after)
	added := Diff(after, before)

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
		if format {
			args := []string{"-U", newUuid, "-f", "-m", mode, "-d", mode}
			args = append(args, added...)
			_, err := d.executor.CommandOutput(MKFS, args...)
			if err != nil {
				return err
			}
		}
		err := d.systemd.AddMount(fmt.Sprintf("/dev/disk/by-uuid/%s", newUuid))
		if err != nil {
			return err
		}
	} else {
		if len(added) > 0 {
			args := []string{"device", "add"}
			args = append(args, added...)
			args = append(args, d.config.ExternalDiskDir())
			_, err := d.executor.CommandOutput(BTRFS, args...)
			if err != nil {
				return err
			}

			args = []string{"balance", "start", fmt.Sprintf("-dconvert=%s", mode), fmt.Sprintf("-mconvert=%s", mode), d.config.ExternalDiskDir()}
			_, err = d.executor.CommandOutput(BTRFS, args...)
			if err != nil {
				return err
			}
		}
		if len(after) == 1 {
			args := []string{"balance", "start", fmt.Sprintf("-dconvert=%s", mode), fmt.Sprintf("-mconvert=%s", mode), d.config.ExternalDiskDir()}
			_, err := d.executor.CommandOutput(BTRFS, args...)
			if err != nil {
				return err
			}
		}
	}

	if len(after) > 0 {
		if len(removed) > 0 {
			args := []string{"device", "delete"}
			args = append(args, removed...)
			args = append(args, d.config.ExternalDiskDir())
			_, err := d.executor.CommandOutput(BTRFS, args...)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
