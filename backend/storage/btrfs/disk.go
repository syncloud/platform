package btrfs

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/syncloud/platform/cli"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

const MKFS = "/snap/platform/current/btrfs/bin/mkfs.sh"
const BTRFS = "/snap/platform/current/btrfs/bin/btrfs.sh"

type Config interface {
	ExternalDiskDir() string
}

type Systemd interface {
	AddMount(device string) error
}

type Disks struct {
	config   Config
	executor cli.Executor
	systemd  Systemd
	logger   *zap.Logger
}

func NewDisks(
	config Config,
	executor cli.Executor,
	systemd Systemd,
	logger *zap.Logger) *Disks {

	return &Disks{
		config:   config,
		executor: executor,
		systemd:  systemd,
		logger:   logger,
	}
}

func (d *Disks) Update(existingDevices []string, newDevices []string, existingUuid string, format bool) (string, error) {
	newUuid := existingUuid
	if format {
		newUuid = uuid.New().String()
	}
	err := d.apply(existingDevices, newDevices, newUuid, format)
	if err != nil {
		return "", err
	}
	return newUuid, nil
}

func (d *Disks) apply(before []string, after []string, newUuid string, format bool) error {
	removed := Diff(before, after)
	added := Diff(after, before)

	//if len(removed) == len(added) {
	//	for i, _ := range removed {
	//		changes = append(changes, NewChange("replace", removed[i], added[i]))
	//	}
	//	return changes, nil
	//}

	mode := "single"
	if len(after) >= 2 {
		mode = "raid1"
	}
	if len(after) >= 4 {
		mode = "raid10"
	}

	if format {
		args := []string{"-U", newUuid, "-f", "-m", mode, "-d", mode}
		args = append(args, after...)
		output, err := d.executor.CombinedOutput(MKFS, args...)
		if err != nil {
			d.logger.Info("error", zap.String("output", string(output)))
			return err
		}
		err = d.systemd.AddMount(fmt.Sprintf("/dev/disk/by-uuid/%s", newUuid))
		if err != nil {
			return err
		}

		return nil
	}

	if len(after) > 0 {
		err := d.systemd.AddMount(fmt.Sprintf("/dev/disk/by-uuid/%s", newUuid))
		if err != nil {
			return err
		}
	}

	if len(before) > 0 {
		if len(added) > 0 {
			args := []string{"device", "add", "--enqueue", "-f"}
			args = append(args, added...)
			args = append(args, d.config.ExternalDiskDir())
			output, err := d.executor.CombinedOutput(BTRFS, args...)
			if err != nil {
				d.logger.Info("error", zap.String("output", string(output)))
				return err
			}

			args = []string{"balance", "start", "--enqueue", fmt.Sprintf("-dconvert=%s", mode), fmt.Sprintf("-mconvert=%s", mode), d.config.ExternalDiskDir()}
			output, err = d.executor.CombinedOutput(BTRFS, args...)
			if err != nil {
				d.logger.Info("error", zap.String("output", string(output)))
				return err
			}
		}
		if len(after) == 1 {
			args := []string{"balance", "start", "--enqueue", "-f", fmt.Sprintf("-dconvert=%s", mode), fmt.Sprintf("-mconvert=%s", mode), d.config.ExternalDiskDir()}
			output, err := d.executor.CombinedOutput(BTRFS, args...)
			if err != nil {
				d.logger.Info("error", zap.String("output", string(output)))
				return err
			}
		}
	}

	if len(after) > 0 {
		if len(removed) > 0 {
			args := []string{"device", "delete", "--enqueue"}
			args = append(args, removed...)
			args = append(args, d.config.ExternalDiskDir())
			output, err := d.executor.CombinedOutput(BTRFS, args...)
			if err != nil {
				d.logger.Info("error", zap.String("output", string(output)))
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
