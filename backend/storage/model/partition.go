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
	Size       string `json:"size"`
	Device     string `json:"device"`
	MountPoint string `json:"mount_point"`
	Active     bool `json:"active"`
	FsType     string `json:"fs_type"`
	Mountable  bool `json:"mountable"`
	Extendable bool `json:"extendable"`
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
