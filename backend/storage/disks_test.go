package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/storage/model"
	"testing"
)

type CallOrder struct {
	order int
}

func (c *CallOrder) inc() int {
	c.order++
	return c.order
}

type DisksConfigStub struct {
	diskDir string
}

func (c *DisksConfigStub) DiskLink() string {
	return "disk"
}

func (c *DisksConfigStub) InternalDiskDir() string {
	return "internaldisk"
}

func (c *DisksConfigStub) ExternalDiskDir() string {
	return c.diskDir
}

type TriggerStub struct {
	error           bool
	callOrderShared *CallOrder
	callOrder       int
}

func (t *TriggerStub) RunDiskChangeEvent() error {
	if t.callOrderShared != nil {
		t.callOrder = t.callOrderShared.inc()
	}
	if t.error {
		return fmt.Errorf("error")
	}
	return nil
}

type LsblkDisksStub struct {
	disks []model.Disk
}

func (l *LsblkDisksStub) FindPartitionByDevice(_ string) (*model.Partition, error) {
	return &l.disks[0].Partitions[0], nil
}

func (l *LsblkDisksStub) AvailableDisks() ([]model.Disk, error) {
	return l.disks, nil
}

func (l *LsblkDisksStub) AllDisks() ([]model.Disk, error) {
	return l.disks, nil
}

type SystemdStub struct {
	callOrderShared   *CallOrder
	callOrder         int
	addMountCalled    bool
	removeMountCalled bool
}

func (s *SystemdStub) AddMount(_ string) error {
	s.addMountCalled = true
	return nil
}

func (s *SystemdStub) RemoveMount() error {
	s.removeMountCalled = true
	if s.callOrderShared != nil {
		s.callOrder = s.callOrderShared.inc()
	}
	return nil
}

type DisksFreeSpaceCheckerStub struct {
	freeSpace bool
}

func (d DisksFreeSpaceCheckerStub) HasFreeSpace(_ string) (bool, error) {
	return d.freeSpace, nil
}

type DisksLinkerStub struct {
	error bool
}

func (d DisksLinkerStub) RelinkDisk(_ string, _ string) error {
	if d.error {
		return fmt.Errorf("error")
	}
	return nil
}

type DisksExecutorStub struct {
	command string
	args    []string
}

func (e *DisksExecutorStub) CombinedOutput(command string, args ...string) ([]byte, error) {
	e.command = command
	e.args = args
	return []byte(""), nil
}

type BtrfsDisksStub struct {
	existingDevices []string
	newDevices      []string
	uuid            string
	format          bool
	error           bool
}

func (b *BtrfsDisksStub) Update(existingDevices []string, newDevices []string, uuid string, format bool) (string, error) {
	if !b.error {
		b.existingDevices = existingDevices
		b.newDevices = newDevices
		b.uuid = uuid
		b.format = format
		return uuid, nil
	} else {
		return "", fmt.Errorf("expected error")
	}
}

func TestDisks_RootPartition_HasFreeSpace_Extendable(t *testing.T) {

	allDisks := []model.Disk{
		{Partitions: []model.Partition{{MountPoint: "/", Active: true}}},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{freeSpace: true}, &DisksLinkerStub{}, &DisksExecutorStub{}, &BtrfsDisksStub{}, &BtrfsDiskStatsStub{}, log.Default())
	partition, err := disks.RootPartition()
	assert.Nil(t, err)
	assert.True(t, partition.Extendable)
}

func TestDisks_RootPartition_HasNoFreeSpace_NonExtendable(t *testing.T) {

	allDisks := []model.Disk{
		{Partitions: []model.Partition{{MountPoint: "/", Active: true}}},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{freeSpace: false}, &DisksLinkerStub{}, &DisksExecutorStub{}, &BtrfsDisksStub{}, &BtrfsDiskStatsStub{}, log.Default())
	partition, err := disks.RootPartition()
	assert.Nil(t, err)
	assert.False(t, partition.Extendable)
}

func TestDisks_DeactivateDisk_TriggerError_NotFail(t *testing.T) {

	allDisks := []model.Disk{
		{Partitions: []model.Partition{{MountPoint: "/", Active: true}}},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: true}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, &BtrfsDisksStub{}, &BtrfsDiskStatsStub{}, log.Default())
	err := disks.Deactivate()
	assert.Nil(t, err)
}

func TestDisks_DeactivateDisk_TriggerNotError_NotFail(t *testing.T) {

	allDisks := []model.Disk{
		{Partitions: []model.Partition{{MountPoint: "/", Active: true}}},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, &BtrfsDisksStub{}, &BtrfsDiskStatsStub{}, log.Default())
	err := disks.Deactivate()
	assert.Nil(t, err)
}

func TestDisks_DeactivateDisk_TriggerEventBeforeRemove(t *testing.T) {

	allDisks := []model.Disk{
		{Partitions: []model.Partition{{MountPoint: "/", Active: true}}},
	}
	callOrder := &CallOrder{order: 0}
	trigger := &TriggerStub{error: false, callOrderShared: callOrder}
	systemd := &SystemdStub{callOrderShared: callOrder}
	disks := NewDisks(&DisksConfigStub{}, trigger, &LsblkDisksStub{disks: allDisks}, systemd, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, &BtrfsDisksStub{}, &BtrfsDiskStatsStub{}, log.Default())
	err := disks.Deactivate()
	assert.Nil(t, err)
	assert.Less(t, trigger.callOrder, systemd.callOrder)

}

func TestDisks_ActivatePartition_SupportedFs(t *testing.T) {

	allDisks := []model.Disk{
		{Device: "/dev/sda", Partitions: []model.Partition{{Device: "/dev/sda1", MountPoint: "/", Active: true, FsType: "ext4"}}},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, &BtrfsDisksStub{}, &BtrfsDiskStatsStub{}, log.Default())
	err := disks.ActivatePartition("/dev/sda1")
	assert.Nil(t, err)
}

func TestDisks_ActivatePartition_Btrfs(t *testing.T) {

	allDisks := []model.Disk{
		{Device: "/dev/sda", Partitions: []model.Partition{{Device: "/dev/sda1", Active: true, FsType: "btrfs"}}},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, &BtrfsDisksStub{}, &BtrfsDiskStatsStub{}, log.Default())
	err := disks.ActivatePartition("/dev/sda1")
	assert.Nil(t, err)
}

func TestDisks_ActivatePartition_NotSupportedFs(t *testing.T) {

	allDisks := []model.Disk{
		{Device: "/dev/sda", Partitions: []model.Partition{{Device: "/dev/sda1", MountPoint: "/", Active: true, FsType: "fat32"}}},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, &BtrfsDisksStub{}, &BtrfsDiskStatsStub{}, log.Default())
	err := disks.ActivatePartition("/dev/sda1")
	assert.NotNil(t, err)
}

func TestDisks_ActivateDisks_None(t *testing.T) {
	executor := &DisksExecutorStub{}

	allDisks := []model.Disk{
		{Device: "/dev/sda", Partitions: []model.Partition{{Device: "/dev/sda1", MountPoint: "/", Active: true, FsType: "fat32"}}},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, executor, &BtrfsDisksStub{}, &BtrfsDiskStatsStub{}, log.Default())
	err := disks.ActivateDisks([]string{}, false)
	assert.NotNil(t, err)
	assert.Equal(t, err, disks.GetLastError())

}

func TestDisks_ActivateDisks_UseUuid(t *testing.T) {
	allDisks := []model.Disk{
		{Device: "/dev/sda", Partitions: []model.Partition{}, Active: true, Uuid: "uuid1"},
		{Device: "/dev/sdb", Partitions: []model.Partition{}, Uuid: "uuid2"},
	}
	btrfs := &BtrfsDisksStub{}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, btrfs, &BtrfsDiskStatsStub{}, log.Default())
	err := disks.ActivateDisks([]string{"/dev/sdb"}, true)
	assert.Nil(t, err)
	assert.Nil(t, disks.GetLastError())
	assert.Equal(t, "uuid2", btrfs.uuid)
	assert.Equal(t, []string{"/dev/sdb"}, btrfs.newDevices)

}

func TestDisks_ActivateDisks_UseUuidExpand(t *testing.T) {
	allDisks := []model.Disk{
		{Device: "/dev/sda", Partitions: []model.Partition{}, Active: true, Uuid: "uuid1"},
		{Device: "/dev/sdb", Partitions: []model.Partition{}, Uuid: "uuid2"},
	}
	btrfs := &BtrfsDisksStub{}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, btrfs, &BtrfsDiskStatsStub{}, log.Default())
	err := disks.ActivateDisks([]string{"/dev/sda", "/dev/sdb"}, true)
	assert.Nil(t, err)
	assert.Nil(t, disks.GetLastError())
	assert.Equal(t, "uuid1", btrfs.uuid)
	assert.Equal(t, []string{"/dev/sda", "/dev/sdb"}, btrfs.newDevices)

}

func TestDisks_ActivateDisks_0_To_2_UseFirstUuid(t *testing.T) {
	allDisks := []model.Disk{
		{Device: "/dev/sda", Partitions: []model.Partition{}, Uuid: "uuid1"},
		{Device: "/dev/sdb", Partitions: []model.Partition{}},
	}
	btrfs := &BtrfsDisksStub{}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, btrfs, &BtrfsDiskStatsStub{}, log.Default())
	err := disks.ActivateDisks([]string{"/dev/sda", "/dev/sdb"}, true)
	assert.Nil(t, err)
	assert.Nil(t, disks.GetLastError())
	assert.Equal(t, "uuid1", btrfs.uuid)
	assert.Equal(t, []string{"/dev/sda", "/dev/sdb"}, btrfs.newDevices)

}

func TestDisks_ActivateDisks_PartitionToDisk_Deactivate(t *testing.T) {
	allDisks := []model.Disk{
		{Device: "/dev/sda", Partitions: []model.Partition{{Device: "/dev/sda1", MountPoint: "/", Active: true, FsType: "fat32"}}},
	}
	btrfs := &BtrfsDisksStub{}
	systemd := &SystemdStub{}

	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, systemd, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, btrfs, &BtrfsDiskStatsStub{}, log.Default())
	err := disks.ActivateDisks([]string{"/dev/sda"}, true)
	assert.Nil(t, err)
	assert.Nil(t, disks.GetLastError())
	assert.True(t, systemd.removeMountCalled)
	//assert.True(t, systemd.addMountCalled)
	assert.Equal(t, []string{"/dev/sda"}, btrfs.newDevices)

}

func TestDisks_ActivateDisks_BterfsError(t *testing.T) {
	allDisks := []model.Disk{
		{Device: "/dev/sda", Partitions: []model.Partition{}},
	}
	btrfs := &BtrfsDisksStub{error: true}
	systemd := &SystemdStub{}

	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, systemd, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, btrfs, &BtrfsDiskStatsStub{}, log.Default())
	err := disks.ActivateDisks([]string{"/dev/sda"}, true)
	assert.NotNil(t, err)
	assert.Equal(t, err, disks.GetLastError())
	//assert.True(t, systemd.removeMountCalled)

}

func TestDisks_ClearLastError(t *testing.T) {
	var allDisks []model.Disk
	btrfs := &BtrfsDisksStub{error: true}
	systemd := &SystemdStub{}

	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, systemd, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, btrfs, &BtrfsDiskStatsStub{}, log.Default())
	assert.Nil(t, disks.GetLastError())
	err := disks.ActivateDisks([]string{"/dev/sda"}, true)
	assert.NotNil(t, err)
	assert.Equal(t, err, disks.GetLastError())
	disks.ClearLastError()
	assert.Nil(t, disks.GetLastError())
}

func TestDisks_AvailableDisks(t *testing.T) {
	allDisks := []model.Disk{
		{Device: "/dev/loop0", Partitions: []model.Partition{}, Uuid: "uuid1", Active: true},
		{Device: "/dev/loop1", Partitions: []model.Partition{}, Uuid: "uuid1", Active: true},
		{Device: "/dev/loop2", Partitions: []model.Partition{}, Uuid: "uuid2"},
	}
	btrfs := &BtrfsDisksStub{error: true}
	systemd := &SystemdStub{}

	stats := &BtrfsDiskStatsStub{
		profiles: map[string]model.RaidProfiles{
			"uuid1": {Data: "raid1", Metadata: "raid1c3"},
			"uuid2": {Data: "single", Metadata: "dup"},
		},
		devices: []model.BtrfsDevice{
			{Device: "/dev/loop0", Devid: "1", Errors: model.DeviceErrors{Read: 4}},
			{Device: "/dev/loop1", Devid: "2"},
		},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, systemd, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, btrfs, stats, log.Default())
	available, err := disks.AvailableDisks()
	assert.Nil(t, err)
	assert.Len(t, available, 3)
	assert.Equal(t, "raid1", available[0].Raid)
	assert.Equal(t, "raid1c3", available[0].RaidMetadata)
	assert.Equal(t, "1", available[0].Devid)
	assert.Equal(t, &model.DeviceErrors{Read: 4}, available[0].Errors)
	assert.True(t, available[0].HasErrors)
	assert.Equal(t, "2", available[1].Devid)
	assert.Equal(t, &model.DeviceErrors{}, available[1].Errors)
	assert.False(t, available[1].HasErrors)
	assert.Equal(t, "single", available[2].Raid)
	assert.Equal(t, "", available[2].Devid)
	assert.Nil(t, available[2].Errors)
	assert.False(t, available[2].HasErrors)
}

func TestDisks_AvailableDisks_NoActiveDiskSkipsDeviceStats(t *testing.T) {
	allDisks := []model.Disk{
		{Device: "/dev/loop0", Partitions: []model.Partition{}},
	}
	stats := &BtrfsDiskStatsStub{}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &DisksExecutorStub{}, &BtrfsDisksStub{}, stats, log.Default())
	_, err := disks.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, stats.devicesCalls)
}
