package storage

import (
	"github.com/syncloud/platform/storage/kernel"
	"github.com/syncloud/platform/storage/model"
	"go.uber.org/zap"
)

type HealthSmart interface {
	Health(device string) model.SmartHealth
}

type HealthKernel interface {
	Events() (*kernel.Events, error)
}

type HealthPorts interface {
	Port(device string) string
}

type HealthScrub interface {
	Status() (*model.ScrubStatus, error)
}

type HealthBtrfs interface {
	Mounted() bool
}

type HealthSchedule interface {
	StorageCheckSchedule() model.CheckSchedule
}

type HealthGuard interface {
	Check() error
}

type Health struct {
	lsblk    DisksLsblk
	smart    HealthSmart
	kernel   HealthKernel
	ports    HealthPorts
	scrub    HealthScrub
	btrfs    HealthBtrfs
	schedule HealthSchedule
	guard    HealthGuard
	logger   *zap.Logger
}

func NewHealth(lsblk DisksLsblk, smart HealthSmart, kernel HealthKernel, ports HealthPorts, scrub HealthScrub, btrfs HealthBtrfs, schedule HealthSchedule, guard HealthGuard, logger *zap.Logger) *Health {
	return &Health{
		lsblk:    lsblk,
		smart:    smart,
		kernel:   kernel,
		ports:    ports,
		scrub:    scrub,
		btrfs:    btrfs,
		schedule: schedule,
		guard:    guard,
		logger:   logger,
	}
}

func (h *Health) Get() (*model.StorageHealth, error) {
	disks, err := h.lsblk.AvailableDisks()
	if err != nil {
		return nil, err
	}
	events, err := h.kernel.Events()
	if err != nil {
		return nil, err
	}
	result := &model.StorageHealth{
		Disks:    []model.DiskHealth{},
		Schedule: h.schedule.StorageCheckSchedule(),
		Busy:     h.guard.Check() != nil,
	}
	for _, disk := range disks {
		port := h.ports.Port(disk.Device)
		result.Disks = append(result.Disks, model.DiskHealth{
			Device: disk.Device,
			Port:   port,
			Smart:  h.smart.Health(disk.Device),
			Link:   events.Link(port),
			Io:     events.Io(disk.Device),
		})
	}
	if h.btrfs.Mounted() {
		scrub, err := h.scrub.Status()
		if err != nil {
			return nil, err
		}
		result.Scrub = scrub
	}
	return result, nil
}
