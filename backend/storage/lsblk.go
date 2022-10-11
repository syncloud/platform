package storage

import (
	"fmt"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/storage/model"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

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

func (l *Lsblk) parseLsblkOutput() ([]model.LsblkEntry, error) {
	var entries []model.LsblkEntry

	lsblkOutputBytes, err := l.executor.CommandOutput("lsblk", "-Pp", "-o", "NAME,SIZE,TYPE,MOUNTPOINT,PARTTYPE,FSTYPE,MODEL,UUID")
	if err != nil {
		return nil, err
	}
	lsblkOutput := string(lsblkOutputBytes)
	l.logger.Info(lsblkOutput)
	lsblkLines := strings.Split(lsblkOutput, "\n")
	for _, rawLine := range lsblkLines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}
		l.logger.Info("parsing", zap.String("line", line))
		r := *regexp.MustCompile(`NAME="(.*)" SIZE="(.*)" TYPE="(.*)" MOUNTPOINT="(.*)" PARTTYPE="(.*)" FSTYPE="(.*)" MODEL="(.*)" UUID="(.*)"`)
		match := r.FindStringSubmatch(line)
		mountPoint := match[4]

		entries = append(entries, model.LsblkEntry{
			Name:       match[1],
			Size:       match[2],
			DeviceType: match[3],
			MountPoint: mountPoint,
			PartType:   match[5],
			FsType:     match[6],
			Model:      strings.TrimSpace(match[7]),
			Active:     l.isActive(mountPoint),
			Uuid:       match[8],
		})

	}

	return entries, nil
}

func (l *Lsblk) extractActiveUuid(entries []model.LsblkEntry) map[string]bool {
	uuids := make(map[string]bool)
	for _, entry := range entries {
		if entry.Active && entry.Uuid != "" {
			uuids[entry.Uuid] = true
		}
	}
	return uuids
}

func (l *Lsblk) AllDisks() (*[]model.Disk, error) {
	disks := make(map[string]*model.Disk)
	entries, err := l.parseLsblkOutput()
	activeUuids := l.extractActiveUuid(entries)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsSupportedType() && entry.IsSupportedFsType() {
			device := entry.Name
			diskName := entry.Model
			l.logger.Info("adding", zap.String("disk", diskName))
			active := entry.Active
			if !active {
				active = activeUuids[entry.Uuid]
			}
			disk := model.NewDisk(diskName, device, entry.Size, active, []model.Partition{})
			if entry.IsSinglePartitionDisk() {
				l.logger.Info("adding single partition", zap.String("disk", device))
				disk.Name = entry.DeviceType
				partition := l.createPartition(entry)
				disk.AddPartition(partition)
			}

			disks[device] = disk

		} else if entry.DeviceType == "part" {
			l.logger.Info("adding", zap.String("regular partition", entry.Name))
			partition := l.createPartition(entry)
			parentDevice := entry.ParentDevice()

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
	return model.Partition{
		Size:       lsblkEntry.Size,
		Device:     lsblkEntry.Name,
		MountPoint: lsblkEntry.MountPoint,
		Active:     lsblkEntry.Active,
		FsType:     lsblkEntry.GetFsType()}
}

func (l *Lsblk) isActive(mountPoint string) bool {
	active := false
	if mountPoint == l.systemConfig.ExternalDiskDir() && l.pathChecker.ExternalDiskLinkExists() {
		active = true
	}
	return active
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
	return nil, fmt.Errorf("unable to find device: %s", device)
}
