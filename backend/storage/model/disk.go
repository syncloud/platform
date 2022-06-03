package model

import (
	"fmt"
	"strings"
)

type Disk struct {
	Name       string
	Device     string
	Size       string
	Partitions []Partition
	Active     bool
}

func NewDisk(name string, device string, size string, partitions []Partition) *Disk {
	if name == "" {
		name = "Disk"
	}
	return &Disk{
		Name:       name,
		Device:     device,
		Size:       size,
		Partitions: partitions,
		Active:     false,
	}
}

func (d *Disk) IsInternal() bool {
	return strings.HasPrefix(d.Device, "/dev/mmcblk")
}

func (d *Disk) HasRootPartition() bool {
	return d.findRootPartition() != nil
}

func (d *Disk) AddPartition(partition Partition) {
	if partition.Active {
		d.Active = true
	}
	d.Partitions = append(d.Partitions, partition)
}

func (d *Disk) findRootPartition() *Partition {
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
