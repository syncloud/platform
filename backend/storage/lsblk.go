package storage

import (
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/storage/model"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

type Lsblk struct {
	systemConfig config.SystemConfig
	pathChecker Checker
	executor cli.CommandExecutor
	logger *zap.Logger
}

func NewLsblk(systemConfig config.SystemConfig, pathChecker Checker, executor cli.CommandExecutor, logger *zap.Logger) *Lsblk {
	return &Lsblk{
		systemConfig: systemConfig,
		pathChecker: pathChecker,
		executor: executor,
		logger: logger,
	}
}

func (l *Lsblk) availableDisks() {
	return [d for d in self.all_disks() if not d.is_internal() and not d.has_root_partition()]
}

func (l *Lsblk) allDisks() (*[]model.Disk, error) {
	lsblkOutputBytes, err := l.executor.CommandOutput("lsblk", "-Pp", "-o", "NAME,SIZE,TYPE,MOUNTPOINT,PARTTYPE,FSTYPE,MODEL")
	if err != nil {
		return nil, err
	}
	lsblkOutput := string(lsblkOutputBytes)
	l.logger.Info(lsblkOutput)

	var disks map[string]model.Disk

	lsblkLines := strings.Split(lsblkOutput, "\n")

	for _, rawLine := range lsblkLines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		l.logger.Info("parsing", zap.String("line", line))
		r := *regexp.MustCompile(`NAME="(.*)" SIZE="(.*)" TYPE="(.*)" MOUNTPOINT="(.*)" PARTTYPE="(.*)" FSTYPE="(.*)" MODEL="(.*)"`)
		match := r.FindStringSubmatch(line)

		lsblkEntry := model.LsblkEntry{
			match[1],
			match[2],
			match[3],
			match[4],
			match[5],
			match[6],
			strings.TrimSpace(match[7]),

		}

		if lsblk_entry.is_supported_type() and
		lsblk_entry.is_supported_fs_type():
		device = lsblk_entry.name
		disk_name = lsblk_entry.model
		self.log.info('adding disk: {0}'.format(disk_name))
		disk = Disk(disk_name, device, lsblk_entry.size,[])
		if lsblk_entry.is_single_partition_disk():
		self.log.info('adding single partition disk: {0}'.format(device))

		disk.name = lsblk_entry.
		type partition = self.create_partition
		(lsblk_entry)
		disk.add_partition(partition)

		disks[device] = disk

		elif
		lsblk_entry.
		type == 'part':
		self.log.info('adding regular partition: {0}'.format(lsblk_entry.name))
		partition = self.create_partition(lsblk_entry)
		parent_device = lsblk_entry.parent_device()
		if parent_device in
	disks:
		disk = disks[parent_device]
		disk.add_partition(partition)
	}

	return disks.values()
}

func (l *Lsblk) create_partition(self, lsblk_entry) {
	mountable = False
	mount_point = lsblk_entry.mount_point
	if not lsblk_entry.is_extended_partition():
	if not mount_point
	or
	mount_point == self.platform_config.get_external_disk_dir():
	mountable = True

	if lsblk_entry.is_boot_disk():
	mountable = False
	active = False
	if mount_point == self.platform_config.get_external_disk_dir() \
	and
	self.path_checker.external_disk_link_exists():
	active = True

	return Partition(lsblk_entry.size, lsblk_entry.name, mount_point, active, lsblk_entry.get_fs_type(), mountable)
}

func (l *Lsblk) is_external_disk_attached(self, lsblk_output=None, disk_dir=None) {
	if not disk_dir:
	disk_dir = self.platform_config.get_external_disk_dir()
	for disk
	in
	self.all_disks(lsblk_output):
	for partition
	in
	disk.partitions:
	if partition.mount_point == disk_dir:
	self.log.info('external disk is attached')
	return True
	self.log.info('external disk is detached')
	return False
}

func (l *Lsblk) find_partition_by_device(self, device, lsblk_output=None) {
	for disk
	in
	self.all_disks(lsblk_output):
	for partition
	in
	disk.partitions:
	if partition.device == device:
	self.log.info('partition found')
	return partition
	self.log.info('partition not found')
	return None
}

func (l *Lsblk) find_partition_by_dir(self, mount_dir, lsblk_output=None) {
	for disk
	in
	self.all_disks(lsblk_output):
	for partition
	in
	disk.partitions:
	if partition.mount_point == mount_dir:
	self.log.info('partition found')
	return partition
	self.log.info('partition not found')
	return None
}