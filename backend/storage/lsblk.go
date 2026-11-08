package storage

import (
	"encoding/json"
	"fmt"
	"github.com/syncloud/platform/storage/model"
	"go.uber.org/zap"
	"sort"
	"strings"
)

const LsblkColumns = "NAME,SIZE,TYPE,MOUNTPOINT,PARTTYPE,FSTYPE,MODEL,UUID,SERIAL,WWN,HCTL"

type Lsblk struct {
	systemConfig Config
	pathChecker  Checker
	executor     LsblkExecutor
	logger       *zap.Logger
}

type LsblkExecutor interface {
	Output(name string, arg ...string) ([]byte, error)
}

type Config interface {
	ExternalDiskDir() string
}

type ByDevice []model.Disk

func (a ByDevice) Len() int           { return len(a) }
func (a ByDevice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDevice) Less(i, j int) bool { return a[i].Device < a[j].Device }

func NewLsblk(config Config, pathChecker Checker, executor LsblkExecutor, logger *zap.Logger) *Lsblk {
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
	output, err := l.executor.Output("lsblk", "-Jp", "-o", LsblkColumns)
	if err != nil {
		return nil, err
	}
	var parsed model.LsblkOutput
	if err := json.Unmarshal(output, &parsed); err != nil {
		return nil, fmt.Errorf("unable to parse lsblk output: %w", err)
	}
	return l.flatten(parsed.BlockDevices, "", nil), nil
}

func (l *Lsblk) flatten(devices []model.LsblkDevice, parent string, entries []model.LsblkEntry) []model.LsblkEntry {
	for _, device := range devices {
		entries = append(entries, model.LsblkEntry{
			Name:       device.Name,
			Size:       device.Size,
			DeviceType: device.Type,
			MountPoint: device.MountPoint,
			PartType:   device.PartType,
			FsType:     device.FsType,
			Model:      strings.TrimSpace(device.Model),
			Active:     l.isActive(device.MountPoint),
			Uuid:       device.Uuid,
			Serial:     strings.TrimSpace(device.Serial),
			Wwn:        device.Wwn,
			Hctl:       device.Hctl,
			Parent:     parent,
		})
		entries = l.flatten(device.Children, device.Name, entries)
	}
	return entries
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
			if disk, ok := disks[entry.Parent]; ok {
				disk.AddPartition(l.createPartition(entry))
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
