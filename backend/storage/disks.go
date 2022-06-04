package storage

import (
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/storage/model"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

const ExtendableFreePercent = 10

//var supportedFilesystems []string

//func init() {
//	supportedFilesystems = []string{"ext2", "ext3", "ext4", "raid"}
//}

type Disks struct {
	config  DisksConfig
	trigger DisksEventTrigger
	lsblk   DisksLsblk
	//pathChecker PathChecker
	systemd  DisksSystemd
	executor cli.CommandExecutor
	logger   *zap.Logger
}

type DisksLsblk interface {
	AvailableDisks() (*[]model.Disk, error)
	AllDisks() (*[]model.Disk, error)
}

type DisksConfig interface {
	DiskLink() string
	InternalDiskDir() string
}

type DisksEventTrigger interface {
	RunDiskChangeEvent() error
}

type DisksSystemd interface {
	RemoveMount() error
}

func NewDisks(
	config DisksConfig,
	trigger DisksEventTrigger,
	lsblk DisksLsblk,
	//pathChecker PathChecker,
	systemd DisksSystemd,
	executor cli.CommandExecutor,
	logger *zap.Logger) *Disks {

	return &Disks{
		config:  config,
		systemd: systemd,
		trigger: trigger,
		lsblk:   lsblk,
		//pathChecker: pathChecker,
		executor: executor,
		logger:   logger,
	}
}

func (d *Disks) RootPartition() (*model.Partition, error) {
	disks, err := d.lsblk.AllDisks()
	if err != nil {
		return nil, err
	}

	for _, disk := range *disks {
		partition := disk.FindRootPartition()
		if partition != nil {
			extendable, err := d.hasUnallocatedSpaceAtTheEnd(disk.Device)
			if err != nil {
				return nil, err
			}
			partition.Extendable = extendable
			return partition, nil
		}
	}
	return &model.Partition{Size: "0", Device: "unknown", MountPoint: "/", Active: true, FsType: "unknown"}, nil

}

func (d *Disks) hasUnallocatedSpaceAtTheEnd(device string) (bool, error) {
	output, err := d.executor.CommandOutput("parted", device, "unit", "%", "print", "free", "--script", "--machine")
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(output), "\n")
	last := lines[len(lines)-1]
	if !strings.Contains(last, "free") {
		return false, nil
	}
	freeString := strings.Split(last, ":")[3]
	freeString = strings.TrimSuffix(freeString, "%")
	free, err := strconv.ParseFloat(freeString, 64)
	if err != nil {
		return false, err
	}

	return free > ExtendableFreePercent, nil
}

func (d *Disks) AvailableDisks() (*[]model.Disk, error) {
	return d.lsblk.AvailableDisks()
}

/*
func (d *Disks) activate_disk(device) {
	self.log.info('activate disk: {0}'.format(device))
	self.deactivate_disk()

	partition = self.lsblk.find_partition_by_device(device)
	if not partition:
	error_message = 'unable to find device: {0}'.format(device)
	self.log.error(error_message)
	raise
	Exception(error_message)

	fs_type = partition.fs_type
	if fs_type not
	in
supported_fs:
	error_message = 'Filesystem type is not supported: {0}' \
	', use one of the following: {1}'.format(fs_type, supported_fs)
	self.log.error(error_message)
	raise
	ServiceException(error_message)

	self.systemctl.add_mount(device)

	relink_disk(
		self.platform_config.get_disk_link(),
		self.platform_config.get_external_disk_dir())
	self.event_trigger.trigger_app_event_disk()
}
*/

func (d *Disks) DeactivateDisk() error {
	d.logger.Info("deactivate disk")
	err := d.RelinkDisk(d.config.DiskLink(), d.config.InternalDiskDir())
	if err != nil {
		return err
	}
	err = d.trigger.RunDiskChangeEvent()
	if err != nil {
		return err
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

func (d *Disks) RelinkDisk(link string, target string) error {

	err := os.Chmod(target, 0o755)
	if err != nil {
		return err
	}

	fi, err := os.Lstat(link)
	if err != nil {
		return err
	}
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		err = os.Remove(link)
		if err != nil {
			return err
		}
	}
	err = os.Symlink(target, link)
	return err
}
