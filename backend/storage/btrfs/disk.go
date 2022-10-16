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

type Change struct {
	Cmd  string
	Args []string
}

func (d *Disks) Update(devices []string, uuid string) (string, error) {
	existingDevices, err := d.ExistingDevices(uuid)
	if err != nil {
		return "", err
	}
	changes, err := DetectChange(existingDevices, devices)
	if err != nil {
		return "", err
	}

	for _, change := range changes {
		d.logger.Info(change.Cmd)
	}

	return "", nil
}

func DetectChange(before []string, after []string) ([]Change, error) {
	removed := Diff(before, after)
	added := Diff(after, before)
	var changes []Change

	if len(removed) == len(added) {
		for i, _ := range removed {
			changes = append(changes, Change{
				Cmd:  "replace",
				Args: []string{removed[i], added[i]},
			})
		}
		return changes, nil
	}

	if len(before) == 0 {
		change := Change{Cmd: "create"}
		for _, v := range added {
			change.Args = append(change.Args, v)
		}
		changes = append(changes, change)
	} else {
		for _, v := range added {
			changes = append(changes, Change{Cmd: "add", Args: []string{v}})
		}
	}

	if len(after) == 0 {
		change := Change{Cmd: "disable"}
		for _, v := range removed {
			change.Args = append(change.Args, v)
		}
		changes = append(changes, change)
	} else {
		for _, v := range removed {
			changes = append(changes, Change{Cmd: "remove", Args: []string{v}})
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

func (d *Disks) create(devices []string) (string, error) {
	mode := "single"
	if len(devices) > 1 {
		mode = "raid1"
	}

	diskUuid := uuid.New().String()
	args := []string{"-U", diskUuid, "-f", "-m", mode, "-d", mode}
	for _, device := range devices {
		args = append(args, device)
	}

	output, err := d.executor.CommandOutput(MKFS, args...)
	if err != nil {
		d.logger.Error(string(output))
		return "", err
	}
	d.logger.Info("mkfs", zap.String("output", string(output)))

	return diskUuid, nil
}

func (d *Disks) Add(devices []string) error {

	for _, device := range devices {
		output, err := d.executor.CommandOutput(BTRFS, "device", "add", device, d.config.ExternalDiskDir())
		if err != nil {
			d.logger.Error(string(output))
			return err
		}
		d.logger.Info("btrfs", zap.String("output", string(output)))
	}

	output, err := d.executor.CommandOutput(BTRFS, "balance", "start", "-dconvert=raid1", "-mconvert=raid1", d.config.ExternalDiskDir())
	if err != nil {
		d.logger.Error(string(output))
		return err
	}
	d.logger.Info("btrfs", zap.String("output", string(output)))

	return nil
}

func (d *Disks) Remove(devices []string) error {
	output, err := d.executor.CommandOutput(BTRFS, "balance", "start", "-sconvert=single", "-dconvert=single", "-mconvert=single", d.config.ExternalDiskDir())
	if err != nil {
		d.logger.Error(string(output))
		return err
	}
	d.logger.Info("btrfs", zap.String("output", string(output)))

	for _, device := range devices {
		output, err := d.executor.CommandOutput(BTRFS, "device", "delete", device, d.config.ExternalDiskDir())
		if err != nil {
			d.logger.Error(string(output))
			return err
		}
		d.logger.Info("btrfs", zap.String("output", string(output)))
	}

	return nil
}
