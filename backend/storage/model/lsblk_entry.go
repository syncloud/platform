package model

import (
	"golang.org/x/exp/slices"
	"regexp"
	"strings"
)

const PartTypeExtended = "0x5"

var SupportedDeviceTypes []string

func init() {
	SupportedDeviceTypes = []string{"disk", "loop"}
}

type LsblkEntry struct {
	Name       string
	Size       string
	DeviceType string
	MountPoint string
	PartType   string
	FsType     string
	Model      string
}

func (e *LsblkEntry) IsExtendedPartition() bool {
	return e.PartType == PartTypeExtended
}

func (e *LsblkEntry) IsBootDisk() bool {
	return strings.HasPrefix(e.Name, "/dev/mmcblk0")
}

func (e *LsblkEntry) IsSupportedType() bool {
	if slices.Contains(SupportedDeviceTypes, e.DeviceType) {
		return true
	}
	if strings.HasPrefix(e.DeviceType, "raid") {
		return true
	}
	return false
}

func (e *LsblkEntry) IsSupportedFsType() bool {
	if e.FsType == "squashfs" {
		return false
	}
	if e.FsType == "linux_raid_member" {
		return false
	}
	return true
}

func (e *LsblkEntry) IsSinglePartitionDisk() bool {
	if e.DeviceType == "loop" {
		return true
	}
	if strings.HasPrefix(e.DeviceType, "raid") {
		return true
	}
	return false
}

func (e *LsblkEntry) ParentDevice() string {
	r := *regexp.MustCompile(`(.*?)p?\d*$`)
	match := r.FindStringSubmatch(e.Name)
	return match[1]
}

func (e *LsblkEntry) GetFsType() string {
	if strings.HasPrefix(e.DeviceType, "raid") {
		return "raid"
	}
	return e.FsType
}
