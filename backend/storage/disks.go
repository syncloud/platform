package storage

import (
	"fmt"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/storage/model"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"strings"
)

const ExtendableFreePercent = 10

var supportedFilesystems []string

func init() {
	supportedFilesystems = []string{"ext2", "ext3", "ext4", "raid", "btrfs"}
}

type Disks struct {
	config           DisksConfig
	trigger          DisksEventTrigger
	lsblk            DisksLsblk
	systemd          DisksSystemd
	freeSpaceChecker DisksFreeSpaceChecker
	linker           DisksLinker
	executor         cli.Executor
	btrfs            BtrfsDisks
	btrfsStats       BtrfsDiskStats
	lastError        error
	logger           *zap.Logger
}

type DisksLsblk interface {
	AvailableDisks() ([]model.Disk, error)
	AllDisks() ([]model.Disk, error)
	FindPartitionByDevice(device string) (*model.Partition, error)
}

type DisksConfig interface {
	DiskLink() string
	InternalDiskDir() string
	ExternalDiskDir() string
}

type DisksEventTrigger interface {
	RunDiskChangeEvent() error
}

type DisksSystemd interface {
	RemoveMount() error
	AddMount(device string) error
}

type DisksLinker interface {
	RelinkDisk(link string, target string) error
}

type DisksFreeSpaceChecker interface {
	HasFreeSpace(device string) (bool, error)
}

type BtrfsDisks interface {
	Update(existingDevices []string, newDevices []string, uuid string, format bool) (string, error)
}

type BtrfsDiskStats interface {
	RaidMode(uuid string) (string, error)
	HasErrors(device string) (bool, error)
}

func NewDisks(
	config DisksConfig,
	trigger DisksEventTrigger,
	lsblk DisksLsblk,
	systemd DisksSystemd,
	freeSpaceChecker DisksFreeSpaceChecker,
	linker DisksLinker,
	executor cli.Executor,
	btrfs BtrfsDisks,
	btrfsStats BtrfsDiskStats,
	logger *zap.Logger) *Disks {

	return &Disks{
		config:           config,
		systemd:          systemd,
		trigger:          trigger,
		lsblk:            lsblk,
		freeSpaceChecker: freeSpaceChecker,
		linker:           linker,
		executor:         executor,
		btrfs:            btrfs,
		btrfsStats:       btrfsStats,
		logger:           logger,
	}
}

func (d *Disks) RootPartition() (*model.Partition, error) {
	disks, err := d.lsblk.AllDisks()
	if err != nil {
		return nil, err
	}

	for _, disk := range disks {
		partition := disk.FindRootPartition()
		if partition != nil {
			extendable, err := d.freeSpaceChecker.HasFreeSpace(disk.Device)
			if err != nil {
				return nil, err
			}
			partition.Extendable = extendable
			return partition, nil
		}
	}
	return &model.Partition{Size: "0", Device: "unknown", MountPoint: "/", Active: true, FsType: "unknown"}, nil

}

func (d *Disks) AvailableDisks() ([]model.Disk, error) {
	disks, err := d.lsblk.AvailableDisks()
	if err != nil {
		return nil, err
	}
	for i, disk := range disks {
		mode, err := d.btrfsStats.RaidMode(disk.Uuid)
		if err != nil {
			d.logger.Info("unable to get raid mode", zap.String("device", disk.Device), zap.String("uuid", disk.Uuid))
		} else {
			disks[i].Raid = mode
		}
		hasErrors, err := d.btrfsStats.HasErrors(disk.Device)
		if err != nil {
			d.logger.Info("unable to get errors", zap.String("device", disk.Device))
		} else {
			disks[i].HasErrors = hasErrors
		}
	}
	return disks, err
}

func (d *Disks) ActivateDisks(newDevices []string, format bool) error {
	err := d.activateDisks(newDevices, format)
	d.lastError = err
	return err
}

func (d *Disks) activateDisks(newDevices []string, format bool) error {
	d.logger.Info("activate disks", zap.Strings("disks", newDevices), zap.Bool("format", format))
	if len(newDevices) < 1 {
		return fmt.Errorf("cannot activate 0 disks")
	}
	disks, err := d.lsblk.AllDisks()
	if err != nil {
		return err
	}

	err = d.Deactivate()
	if err != nil {
		return err
	}

	existingDevices := d.activeDevices(disks)
	existingUuid := d.firstActiveUuid(newDevices, disks)
	newUuid := d.firstUuid(newDevices, disks)

	var uuid *string
	if existingUuid != nil {
		uuid = existingUuid
	} else if newUuid != nil {
		uuid = newUuid
	}

	uuidToUse := ""
	if uuid == nil {
		if !format {
			return fmt.Errorf("cannot find existing uuid to use")
		}
	} else {
		uuidToUse = *uuid
	}

	_, err = d.btrfs.Update(existingDevices, newDevices, uuidToUse, format)
	if err != nil {
		return err
	}
	return d.activateCommon()

}

func (d *Disks) activeDevices(disks []model.Disk) []string {
	var existingDevices []string
	for _, disk := range disks {
		if disk.Active {
			existingDevices = append(existingDevices, disk.Device)
		}
	}
	return existingDevices
}

func (d *Disks) firstActiveUuid(newDevices []string, disks []model.Disk) *string {
	for _, disk := range disks {
		if disk.Active {
			if slices.Contains(newDevices, disk.Device) {
				if disk.Uuid != "" {
					return &disk.Uuid
				}
			}
		}
	}
	return nil
}

func (d *Disks) firstUuid(newDevices []string, disks []model.Disk) *string {
	for _, disk := range disks {
		if slices.Contains(newDevices, disk.Device) {
			if disk.Uuid != "" {
				return &disk.Uuid
			}
		}
	}
	return nil
}

func (d *Disks) ActivatePartition(device string) error {
	d.logger.Info("activate partition", zap.String("disk", device))
	err := d.Deactivate()
	if err != nil {
		return err
	}

	partition, err := d.lsblk.FindPartitionByDevice(device)
	if err != nil {
		return err
	}
	fsType := partition.FsType
	if !slices.Contains(supportedFilesystems, fsType) {
		return fmt.Errorf("filesystem type is not supported: %s, use one of the following: %s", fsType, strings.Join(supportedFilesystems, ","))
	}
	err = d.systemd.AddMount(device)
	if err != nil {
		return err
	}

	return d.activateCommon()
}

func (d *Disks) activateCommon() error {

	err := d.linker.RelinkDisk(d.config.DiskLink(), d.config.ExternalDiskDir())
	if err != nil {
		return err
	}

	return d.trigger.RunDiskChangeEvent()
}

func (d *Disks) Deactivate() error {
	d.logger.Info("deactivate disk")
	err := d.linker.RelinkDisk(d.config.DiskLink(), d.config.InternalDiskDir())
	if err != nil {
		return err
	}
	err = d.trigger.RunDiskChangeEvent()
	if err != nil {
		d.logger.Error("some disk events produced errors", zap.Error(err))
	}
	err = d.systemd.RemoveMount()
	if err != nil {
		return err
	}
	return nil
}

func (d *Disks) GetLastError() error {
	return d.lastError
}

func (d *Disks) ClearLastError() {
	d.lastError = nil
}
