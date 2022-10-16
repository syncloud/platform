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
}

type UiDeviceEntry struct {
	Name   string `json:"name"`
	Device string `json:"device"`
	Size   string `json:"size"`
	Active bool   `json:"active"`
}

func NewDisk(name string, device string, size string, active bool, uuid string, partitions []Partition) *Disk {
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
	}
}

func (d *Disk) IsInternal() bool {
	return strings.HasPrefix(d.Device, "/dev/mmcblk")
}

func (d *Disk) HasRootPartition() bool {
	return d.FindRootPartition() != nil
}

func (d *Disk) IsAvailable() bool {
	available := false
	for _, v := range d.Partitions {
		if v.Active {
			available = true
		} else if v.MountPoint == "" {
			available = true
		}
	}
	if d.IsInternal() {
		available = false
	}
	return available
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
