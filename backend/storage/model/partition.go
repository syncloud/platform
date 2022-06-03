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
	Size       string
	Device     string
	MountPoint string
	Active     bool
	FsType     string
	Mountable  bool
	Extendable bool
}

func (p *Partition) PermissionsSupport() bool {
	return !slices.Contains(NoPartitionFsTypes, p.FsType)
}

func (p *Partition) isRootFs() bool {
	return p.MountPoint == "/"
}

func (p *Partition) ToString() string {
	return fmt.Sprintf("%s, %s, %s, %t", p.Device, p.Size, p.MountPoint, p.Active)
}
