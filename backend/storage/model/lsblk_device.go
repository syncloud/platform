package model

import (
	"golang.org/x/exp/slices"
	"regexp"
	"strings"
)

var SupportedDeviceTypes = []string{"disk", "loop"}

type LsblkDevice struct {
	Name       string        `json:"name"`
	Size       string        `json:"size"`
	Type       string        `json:"type"`
	MountPoint string        `json:"mountpoint"`
	PartType   string        `json:"parttype"`
	FsType     string        `json:"fstype"`
	Model      string        `json:"model"`
	Uuid       string        `json:"uuid"`
	Serial     string        `json:"serial"`
	Wwn        string        `json:"wwn"`
	Hctl       string        `json:"hctl"`
	Children   []LsblkDevice `json:"children"`
}

func (d *LsblkDevice) IsSupportedType() bool {
	if d.IsMMCBootPartition() {
		return false
	}
	if d.IsZram() {
		return false
	}
	if slices.Contains(SupportedDeviceTypes, d.Type) {
		return true
	}
	return d.IsRaid()
}

func (d *LsblkDevice) IsMMCBootPartition() bool {
	return d.DetectMMCBootDevice() != nil
}

func (d *LsblkDevice) IsZram() bool {
	r := regexp.MustCompile(`^/dev/zram\d+$`)
	return r.MatchString(d.Name)
}

func (d *LsblkDevice) DetectMMCBootDevice() *string {
	r := regexp.MustCompile(`^(/dev/mmcblk\d+)boot\d+$`)
	match := r.FindStringSubmatch(d.Name)
	if match != nil {
		return &match[1]
	}
	return nil
}

func (d *LsblkDevice) IsSupportedFsType() bool {
	if d.FsType == "squashfs" {
		return false
	}
	if strings.HasPrefix(d.MountPoint, "/snap") {
		return false
	}
	return d.FsType != "linux_raid_member"
}

func (d *LsblkDevice) IsRaid() bool {
	return strings.HasPrefix(d.Type, "raid")
}

func (d *LsblkDevice) GetFsType() string {
	if d.IsRaid() {
		return "raid"
	}
	return d.FsType
}
