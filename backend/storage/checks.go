package storage

import (
	"errors"
	"fmt"
	"github.com/syncloud/platform/storage/model"
	"go.uber.org/zap"
)

type ChecksSmart interface {
	Health(device string) model.SmartHealth
	StartShortTest(device string) error
}

type ChecksScrub interface {
	Start() error
	Status() (*model.ScrubStatus, error)
}

type ChecksBtrfs interface {
	Mounted() bool
	ResetErrors() error
}

type ChecksConfig interface {
	StorageCheckSchedule() model.CheckSchedule
	SetStorageCheckSchedule(schedule model.CheckSchedule)
}

type ChecksGuard interface {
	Check() error
}

type Checks struct {
	lsblk  DisksLsblk
	smart  ChecksSmart
	scrub  ChecksScrub
	btrfs  ChecksBtrfs
	config ChecksConfig
	guard  ChecksGuard
	logger *zap.Logger
}

func NewChecks(lsblk DisksLsblk, smart ChecksSmart, scrub ChecksScrub, btrfs ChecksBtrfs, config ChecksConfig, guard ChecksGuard, logger *zap.Logger) *Checks {
	return &Checks{
		lsblk:  lsblk,
		smart:  smart,
		scrub:  scrub,
		btrfs:  btrfs,
		config: config,
		guard:  guard,
		logger: logger,
	}
}

func (c *Checks) ScrubAvailable() bool {
	return c.btrfs.Mounted()
}

func (c *Checks) StartScrub() error {
	if !c.btrfs.Mounted() {
		return fmt.Errorf("no active btrfs disk array")
	}
	if err := c.guard.Check(); err != nil {
		return err
	}
	status, err := c.scrub.Status()
	if err != nil {
		return err
	}
	if status.State == model.ScrubRunning {
		return fmt.Errorf("a scrub is already running")
	}
	c.logger.Info("starting read-only scrub")
	return c.scrub.Start()
}

func (c *Checks) StartSmartTest(device string) error {
	disks, err := c.lsblk.AvailableDisks()
	if err != nil {
		return err
	}
	if !containsDevice(disks, device) {
		return fmt.Errorf("unknown disk: %s", device)
	}
	if err := c.guard.Check(); err != nil {
		return err
	}
	c.logger.Info("starting smart short test", zap.String("device", device))
	return c.smart.StartShortTest(device)
}

func (c *Checks) StartSmartTests() error {
	if err := c.guard.Check(); err != nil {
		return err
	}
	disks, err := c.lsblk.AvailableDisks()
	if err != nil {
		return err
	}
	var errs []error
	for _, disk := range disks {
		health := c.smart.Health(disk.Device)
		if !health.Available || health.SelfTest.Running {
			continue
		}
		c.logger.Info("starting scheduled smart short test", zap.String("device", disk.Device))
		if err := c.smart.StartShortTest(disk.Device); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", disk.Device, err))
		}
	}
	return errors.Join(errs...)
}

func (c *Checks) ResetErrors() error {
	if !c.btrfs.Mounted() {
		return fmt.Errorf("no active btrfs disk array")
	}
	c.logger.Info("resetting btrfs device error counters")
	return c.btrfs.ResetErrors()
}

func (c *Checks) Schedule() model.CheckSchedule {
	return c.config.StorageCheckSchedule()
}

func (c *Checks) SetSchedule(schedule model.CheckSchedule) {
	c.config.SetStorageCheckSchedule(schedule)
}

func containsDevice(disks []model.Disk, device string) bool {
	for _, disk := range disks {
		if disk.Device == device {
			return true
		}
	}
	return false
}
