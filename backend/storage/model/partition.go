package model

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
	if p.FsType == "vfat" {
		return false
	}
	if p.FsType == "exfat" {
		return false
	}
	return true
}

func (p *Partition) isRootFs() bool {
	return p.MountPoint == "/"
}

func (p *Partition) ToString() bool {
	return '{0}, {1}, {2}, {3}'.format(self.device, self.size, self.mount_point, self.active)
}
