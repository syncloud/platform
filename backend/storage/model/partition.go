package model

import (
	"fmt"
	"golang.org/x/exp/slices"
)

var NoPartitionFsTypes []string

func init() {
	NoPartitionFsTypes = []string{"vfat", "exfat"}
}

type Partition struct {
	Size       int64
	Device     string
	MountPoint string
	Active     bool
	FsType     string
	Mountable  bool
	Extendable bool
}

func NewPartition(size int64, device string, mountPoint string, active bool, fsType string, mountable bool) *Partition {
	return &Partition{
		Size:       size,
		Device:     device,
		MountPoint: mountPoint,
		Active:     active,
		FsType:     fsType,
		Mountable:  mountable,
		Extendable: false,
	}
}

func (p *Partition) PermissionsSupport() bool {
	return !slices.Contains(NoPartitionFsTypes, p.FsType)
}

func (p *Partition) isRootFs() bool {
	return p.MountPoint == "/"
}

func (p *Partition) ToString() string {
	return fmt.Sprintf("%s, %d, %s, %t", p.Device, p.Size, p.MountPoint, p.Active)
}
