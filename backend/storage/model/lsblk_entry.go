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
	Active     bool
	Uuid       string
}

func (e *LsblkEntry) IsExtendedPartition() bool {
	return e.PartType == PartTypeExtended
}

func (e *LsblkEntry) IsSupportedType() bool {
	if slices.Contains(SupportedDeviceTypes, e.DeviceType) {
		return true
	}
	if e.IsRaid() {
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

func (e *LsblkEntry) IsRaid() bool {
	if strings.HasPrefix(e.DeviceType, "raid") {
		return true
	}
	return false
}

func (e *LsblkEntry) ParentDevice() (string, error) {
	r, err := regexp.Compile(`(.*?)p?\d*$`)
	if err != nil {
		return "", err
	}

	match := r.FindStringSubmatch(e.Name)
	return match[1], nil
}

func (e *LsblkEntry) GetFsType() string {
	if e.IsRaid() {
		return "raid"
	}
	return e.FsType
}
