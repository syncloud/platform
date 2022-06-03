package storage

import (
	"errors"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/storage/model"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

var ErrNotFound = errors.New("partition not found")

type Lsblk struct {
	systemConfig Config
	pathChecker  Checker
	executor     cli.CommandExecutor
	logger       *zap.Logger
}

type Config interface {
	ExternalDiskDir() string
}

func NewLsblk(config Config, pathChecker Checker, executor cli.CommandExecutor, logger *zap.Logger) *Lsblk {
	return &Lsblk{
		systemConfig: config,
		pathChecker:  pathChecker,
		executor:     executor,
		logger:       logger,
	}
}

func (l *Lsblk) AvailableDisks() (*[]model.Disk, error) {
	var disks []model.Disk

	allDisks, err := l.AllDisks()
	if err != nil {
		return nil, err
	}
	for _, disk := range *allDisks {
		if !disk.IsInternal() && !disk.HasRootPartition() {
			disks = append(disks, disk)
		}
	}
	return &disks, nil
}

func (l *Lsblk) AllDisks() (*[]model.Disk, error) {
	lsblkOutputBytes, err := l.executor.CommandOutput("lsblk", "-Pp", "-o", "NAME,SIZE,TYPE,MOUNTPOINT,PARTTYPE,FSTYPE,MODEL")
	if err != nil {
		return nil, err
	}
	lsblkOutput := string(lsblkOutputBytes)
	l.logger.Info(lsblkOutput)

	disks := make(map[string]*model.Disk)

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
			Name:       match[1],
			Size:       match[2],
			DeviceType: match[3],
			MountPoint: match[4],
			PartType:   match[5],
			FsType:     match[6],
			Model:      strings.TrimSpace(match[7]),
		}

		if lsblkEntry.IsSupportedType() && lsblkEntry.IsSupportedFsType() {
			device := lsblkEntry.Name
			diskName := lsblkEntry.Model
			l.logger.Info("adding", zap.String("disk", diskName))
			disk := model.NewDisk(diskName, device, lsblkEntry.Size, []model.Partition{})
			if lsblkEntry.IsSinglePartitionDisk() {
				l.logger.Info("adding single partition", zap.String("disk", device))
				disk.Name = lsblkEntry.DeviceType
				partition := l.createPartition(lsblkEntry)
				disk.AddPartition(partition)
			}

			disks[device] = disk

		} else if lsblkEntry.DeviceType == "part" {
			l.logger.Info("adding", zap.String("regular partition", lsblkEntry.Name))
			partition := l.createPartition(lsblkEntry)
			parentDevice := lsblkEntry.ParentDevice()

			if _, ok := disks[parentDevice]; ok {
				disk := disks[parentDevice]
				disk.AddPartition(partition)
			}
		}
	}

	var results []model.Disk
	for _, disk := range disks {
		results = append(results, *disk)
	}
	return &results, nil
}

func (l *Lsblk) createPartition(lsblkEntry model.LsblkEntry) model.Partition {
	mountable := false
	mountPoint := lsblkEntry.MountPoint
	if !lsblkEntry.IsExtendedPartition() {
		if mountPoint == "" || mountPoint == l.systemConfig.ExternalDiskDir() {
			mountable = true
		}
	}

	if lsblkEntry.IsBootDisk() {
		mountable = false
	}
	active := false
	if mountPoint == l.systemConfig.ExternalDiskDir() && l.pathChecker.ExternalDiskLinkExists() {
		active = true
	}

	return model.Partition{
		Size:       lsblkEntry.Size,
		Device:     lsblkEntry.Name,
		MountPoint: mountPoint,
		Active:     active,
		FsType:     lsblkEntry.GetFsType(),
		Mountable:  mountable}
}

func (l *Lsblk) FindPartitionByDevice(device string) (*model.Partition, error) {
	disks, err := l.AllDisks()
	if err != nil {
		return nil, err
	}
	for _, disk := range *disks {
		for _, partition := range disk.Partitions {
			if partition.Device == device {
				l.logger.Info("partition found")
				return &partition, nil
			}
		}
	}
	return nil, ErrNotFound
}
