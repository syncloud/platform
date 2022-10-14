package btrfs

import (
	"github.com/google/uuid"
	"github.com/syncloud/platform/cli"
	"go.uber.org/zap"
)

const MKFS = "/snap/platform/current/btrfs/bin/mkfs.sh"
const BTRFS = "/snap/platform/current/btrfs/bin/btrfs.sh"

type Config interface {
	ExternalDiskDir() string
}

type Disks struct {
	config   Config
	executor cli.CommandExecutor
	logger   *zap.Logger
}

func NewDisks(
	config Config,
	executor cli.CommandExecutor,
	logger *zap.Logger) *Disks {

	return &Disks{
		config:   config,
		executor: executor,
		logger:   logger,
	}
}

func (d *Disks) Create(devices []string) (string, error) {
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
