package storage

import (
	"fmt"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/storage/model"
	"go.uber.org/zap"
	"regexp"
	"sort"
	"strings"
)

const LsblkColumns = "NAME,SIZE,TYPE,MOUNTPOINT,PARTTYPE,FSTYPE,MODEL,UUID,SERIAL,WWN,HCTL"

var lsblkField = regexp.MustCompile(`([A-Z:_-]+)="([^"]*)"`)

type Lsblk struct {
	systemConfig Config
	pathChecker  Checker
	executor     cli.Executor
	logger       *zap.Logger
}

type Config interface {
	ExternalDiskDir() string
}

type ByDevice []model.Disk

func (a ByDevice) Len() int           { return len(a) }
func (a ByDevice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDevice) Less(i, j int) bool { return a[i].Device < a[j].Device }

func NewLsblk(config Config, pathChecker Checker, executor cli.Executor, logger *zap.Logger) *Lsblk {
	return &Lsblk{
		systemConfig: config,
		pathChecker:  pathChecker,
		executor:     executor,
		logger:       logger,
	}
}

func (l *Lsblk) AvailableDisks() ([]model.Disk, error) {
	var disks []model.Disk

	allDisks, err := l.AllDisks()
	if err != nil {
		return nil, err
	}
	for _, disk := range allDisks {
		if disk.IsAvailable() {
			disks = append(disks, disk)
		}
	}
	return disks, nil
}

func (l *Lsblk) parseLsblkOutput() ([]model.LsblkEntry, error) {
	var entries []model.LsblkEntry

	lsblkOutputBytes, err := l.executor.CombinedOutput("lsblk", "-Pp", "-o", LsblkColumns)
	if err != nil {
		return nil, err
	}
	for _, rawLine := range strings.Split(string(lsblkOutputBytes), "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}
		fields := make(map[string]string)
		for _, match := range lsblkField.FindAllStringSubmatch(line, -1) {
			fields[match[1]] = match[2]
		}
		mountPoint := fields["MOUNTPOINT"]
		entries = append(entries, model.LsblkEntry{
			Name:       fields["NAME"],
			Size:       fields["SIZE"],
			DeviceType: fields["TYPE"],
			MountPoint: mountPoint,
			PartType:   fields["PARTTYPE"],
			FsType:     fields["FSTYPE"],
			Model:      strings.TrimSpace(fields["MODEL"]),
			Active:     l.isActive(mountPoint),
			Uuid:       fields["UUID"],
			Serial:     strings.TrimSpace(fields["SERIAL"]),
			Wwn:        fields["WWN"],
			Hctl:       fields["HCTL"],
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

func (l *Lsblk) extractBootDevices(entries []model.LsblkEntry) map[string]bool {
	devices := make(map[string]bool)
	for _, entry := range entries {
		device := entry.DetectMMCBootDevice()
		if device != nil {
			devices[*device] = true
		}
	}
	return devices
}

func (l *Lsblk) AllDisks() ([]model.Disk, error) {
	disks := make(map[string]*model.Disk)
	entries, err := l.parseLsblkOutput()
	if err != nil {
		return nil, err
	}
	activeUuids := l.extractActiveUuid(entries)
	bootDevices := l.extractBootDevices(entries)
	for _, entry := range entries {
		device := entry.Name
		if ok, _ := bootDevices[device]; ok {
			entry.Boot = true
		}
		if entry.IsSupportedType() && entry.IsSupportedFsType() {
			diskName := entry.Model
			active := entry.Active
			if !active {
				active = activeUuids[entry.Uuid]
			}
			disk := model.NewDisk(diskName, device, entry.Size, active, entry.Uuid, entry.MountPoint, entry.Boot, []model.Partition{})
			disk.Serial = entry.Serial
			disk.Wwn = entry.Wwn
			disk.Hctl = entry.Hctl
			if entry.IsRaid() {
				disk.Name = entry.DeviceType
				partition := l.createPartition(entry)
				disk.AddPartition(partition)
			}

			disks[device] = disk

		} else if entry.DeviceType == "part" {
			partition := l.createPartition(entry)
			parentDevice, err := entry.ParentDevice()
			if err != nil {
				return nil, err
			}

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
	sort.Sort(ByDevice(results))
	return results, nil
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
	for _, disk := range disks {
		for _, partition := range disk.Partitions {
			if partition.Device == device {
				l.logger.Info("partition found")
				return &partition, nil
			}
		}
	}
	return nil, fmt.Errorf("unable to find device: %s", device)
}
