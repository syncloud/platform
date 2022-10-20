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
	executor         cli.CommandExecutor
	btrfs            BtrfsDisks
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
	Update(devices []string, uuid string, format bool) (string, error)
}

func NewDisks(
	config DisksConfig,
	trigger DisksEventTrigger,
	lsblk DisksLsblk,
	systemd DisksSystemd,
	freeSpaceChecker DisksFreeSpaceChecker,
	linker DisksLinker,
	executor cli.CommandExecutor,
	btrfs BtrfsDisks,
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
	return d.lsblk.AvailableDisks()
}

func (d *Disks) ActivateDisks(devices []string, format bool) error {
	d.logger.Info("activate disks", zap.Strings("disks", devices), zap.Bool("format", format))
	if len(devices) < 1 {
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

	uuid := ""
	for _, disk := range disks {
		if slices.Contains(devices, disk.Device) {
			uuid = disk.Uuid
		}
	}
	uuid, err = d.btrfs.Update(devices, uuid, format)
	return d.activateCommon()

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

/*
func (d *Disks) get_app_storage_dir(app_id) {
	app_storage_dir = join(self.platform_config.get_disk_link(), app_id)
	return app_storage_dir
}

func (d *Disks) init_app_storage(app_id, owner = None) {
	app_storage_dir = self.get_app_storage_dir(app_id)
	if not path.exists(app_storage_dir):
	os.mkdir(app_storage_dir)
	if owner:
	self.log.info('fixing permissions on {0}'.format(app_storage_dir))
	fs.chownpath(app_storage_dir, owner, recursive = True) else:
	self.log.info('not fixing permissions')
	return app_storage_dir
}
func (d *Disks) init_disk() {

	if not isdir(self.platform_config.get_disk_root()):
	os.mkdir(self.platform_config.get_disk_root())

	if not isdir(self.platform_config.get_internal_disk_dir()):
	os.mkdir(self.platform_config.get_internal_disk_dir())

	if not self.path_checker.external_disk_link_exists():
	relink_disk(self.platform_config.get_disk_link(), self.platform_config.get_internal_disk_dir())
}
*/
