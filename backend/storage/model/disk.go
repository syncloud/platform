package model

import (
	"fmt"
	"strings"
)

type Disk struct {
	Name   string
	Device     string
	Size       int64
	Partitions []Partition
	Active     bool
}

func NewDisk(name string, device string, size int64, partitions []Partition) *Disk {
	if name == "" {
		name = "Disk"
	}
	return &Disk{
		Name: name,
		Device: device,
		Size: size,
		Partitions: partitions,
		Active: false,
	}
}

func (d *Disk) isInternal() bool {
	return strings.HasPrefix(d.Device, "/dev/mmcblk")
}

func (d *Disk)  hasRootPartition() bool {
	return find_root_partition()
	is
	not
	None
}

func (d *Disk) AddPartition(partition) {
	if partition.active:
	self.active = True
	self.partitions.append(partition)
}

func (d *Disk) findRootPartition() *Partition {
	return next((p
	for p
	in
	self.partitions
	if p.is_root_fs()), None)
}

func (d *Disk) String() string {
	var partitionStrings []string
	for _, v := range d.Partitions {
		v.ToString()
	}
	return fmt.Sprintf("%s: %s", d.Name, partitionStrings)
}