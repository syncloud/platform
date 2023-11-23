package model

import (
	"fmt"
	"strings"
)

type Disk struct {
	Name       string      `json:"name"`
	Device     string      `json:"device"`
	Size       string      `json:"size"`
	Partitions []Partition `json:"partitions"`
	Active     bool        `json:"active"`
	Uuid       string      `json:"uuid"`
	MountPoint string      `json:"mount_point"`
	Raid       string      `json:"raid"`
	HasErrors  bool        `json:"has_errors"`
	Boot       bool        `json:"boot"`
}

type UiDeviceEntry struct {
	Name   string `json:"name"`
	Device string `json:"device"`
	Size   string `json:"size"`
	Active bool   `json:"active"`
}

func NewDisk(name string, device string, size string, active bool, uuid string, mountPoint string, boot bool, partitions []Partition) *Disk {
	if name == "" {
		name = fmt.Sprintf("Disk %s", strings.TrimPrefix(device, "/dev/"))
	}
	return &Disk{
		Name:       name,
		Device:     device,
		Size:       size,
		Partitions: partitions,
		Active:     active,
		Uuid:       uuid,
		MountPoint: mountPoint,
		Boot:       boot,
	}
}

func (d *Disk) HasRootPartition() bool {
	return d.FindRootPartition() != nil
}

func (d *Disk) IsAvailable() bool {
	if d.HasRootPartition() {
		return false
	}
	if d.Boot {
		return false
	}
	return true
}

func (d *Disk) AddPartition(partition Partition) {
	d.Partitions = append(d.Partitions, partition)
}

func (d *Disk) FindRootPartition() *Partition {
	for _, v := range d.Partitions {
		if v.isRootFs() {
			return &v
		}
	}
	return nil
}

func (d *Disk) String() string {
	var partitionStrings []string
	for _, v := range d.Partitions {
		v.ToString()
	}
	return fmt.Sprintf("%s: %s", d.Name, partitionStrings)
}
