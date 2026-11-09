package storage

import (
	"encoding/json"
	"fmt"
	"github.com/syncloud/platform/storage/model"
	"go.uber.org/zap"
	"sort"
	"strings"
)

const LSBLK = "/snap/platform/current/lsblk/bin/lsblk.sh"

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

func (l *Lsblk) parseLsblkOutput() ([]model.LsblkDevice, error) {
	output, err := l.executor.Output(LSBLK, "-Jp", "-o", LsblkColumns)
	if err != nil {
		return nil, err
	}
	var parsed model.LsblkOutput
	if err := json.Unmarshal(output, &parsed); err != nil {
		return nil, fmt.Errorf("unable to parse lsblk output: %w", err)
	}
	return parsed.BlockDevices, nil
}

func (l *Lsblk) walk(devices []model.LsblkDevice, parent string, visit func(device model.LsblkDevice, parent string)) {
	for _, device := range devices {
		visit(device, parent)
		l.walk(device.Children, device.Name, visit)
	}
}

func (l *Lsblk) AllDisks() ([]model.Disk, error) {
	devices, err := l.parseLsblkOutput()
	if err != nil {
		return nil, err
	}
	activeUuids := make(map[string]bool)
	bootDevices := make(map[string]bool)
	l.walk(devices, "", func(device model.LsblkDevice, _ string) {
		if l.isActive(device.MountPoint) && device.Uuid != "" {
			activeUuids[device.Uuid] = true
		}
		if boot := device.DetectMMCBootDevice(); boot != nil {
			bootDevices[*boot] = true
		}
	})
	disks := make(map[string]*model.Disk)
	l.walk(devices, "", func(device model.LsblkDevice, parent string) {
		if device.IsSupportedType() && device.IsSupportedFsType() {
			active := l.isActive(device.MountPoint) || activeUuids[device.Uuid]
			disk := model.NewDisk(strings.TrimSpace(device.Model), device.Name, device.Size, active, device.Uuid, device.MountPoint, bootDevices[device.Name], []model.Partition{})
			disk.Serial = strings.TrimSpace(device.Serial)
			disk.Wwn = device.Wwn
			disk.Hctl = device.Hctl
			if device.IsRaid() {
				disk.Name = device.Type
				disk.AddPartition(l.createPartition(device))
			}
			disks[device.Name] = disk
		} else if device.Type == "part" {
			if disk, ok := disks[parent]; ok {
				disk.AddPartition(l.createPartition(device))
			}
		}
	})

	var results []model.Disk
	for _, disk := range disks {
		results = append(results, *disk)
	}
	sort.Sort(ByDevice(results))
	return results, nil
}

func (l *Lsblk) createPartition(device model.LsblkDevice) model.Partition {
	return model.Partition{
		Size:       device.Size,
		Device:     device.Name,
		MountPoint: device.MountPoint,
		Active:     l.isActive(device.MountPoint),
		FsType:     device.GetFsType()}
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
