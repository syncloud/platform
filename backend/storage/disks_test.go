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

type StorageExecutorStub struct {
	command string
	args    []string
}

func (e *StorageExecutorStub) CommandOutput(command string, args ...string) ([]byte, error) {
	e.command = command
	e.args = args
	return []byte(""), nil
}

type BtrfsDisksStub struct {
	existingDevices []string
	newDevices      []string
	uuid            string
	format          bool
}

func (b *BtrfsDisksStub) Update(existingDevices []string, newDevices []string, uuid string, format bool) (string, error) {
	b.existingDevices = existingDevices
	b.newDevices = newDevices
	b.uuid = uuid
	b.format = format
	return uuid, nil
}

func TestDisks_RootPartition_HasFreeSpace_Extendable(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false}}, false, "", ""},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{freeSpace: true}, &DisksLinkerStub{}, &StorageExecutorStub{}, &BtrfsDisksStub{}, log.Default())
	partition, err := disks.RootPartition()
	assert.Nil(t, err)
	assert.True(t, partition.Extendable)
}

func TestDisks_RootPartition_HasNoFreeSpace_NonExtendable(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false}}, false, "", ""},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{freeSpace: false}, &DisksLinkerStub{}, &StorageExecutorStub{}, &BtrfsDisksStub{}, log.Default())
	partition, err := disks.RootPartition()
	assert.Nil(t, err)
	assert.False(t, partition.Extendable)
}

func TestDisks_DeactivateDisk_TriggerError_NotFail(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false}}, false, "", ""},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: true}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, &BtrfsDisksStub{}, log.Default())
	err := disks.Deactivate()
	assert.Nil(t, err)
}

func TestDisks_DeactivateDisk_TriggerNotError_NotFail(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false}}, false, "", ""},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, &BtrfsDisksStub{}, log.Default())
	err := disks.Deactivate()
	assert.Nil(t, err)
}

func TestDisks_DeactivateDisk_TriggerEventBeforeRemove(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false}}, false, "", ""},
	}
	callOrder := &CallOrder{order: 0}
	trigger := &TriggerStub{error: false, callOrderShared: callOrder}
	systemd := &SystemdStub{callOrderShared: callOrder}
	disks := NewDisks(&DisksConfigStub{}, trigger, &LsblkDisksStub{disks: allDisks}, systemd, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, &BtrfsDisksStub{}, log.Default())
	err := disks.Deactivate()
	assert.Nil(t, err)
	assert.Less(t, trigger.callOrder, systemd.callOrder)

}

func TestDisks_ActivateDisk_SupportedFs(t *testing.T) {

	allDisks := []model.Disk{
		{"", "/dev/sda", "", []model.Partition{{"", "/dev/sda1", "/", true, "ext4", false}}, false, "", ""},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, &BtrfsDisksStub{}, log.Default())
	err := disks.ActivatePartition("/dev/sda1")
	assert.Nil(t, err)
}

func TestDisks_ActivateDisk_Btrfs(t *testing.T) {

	allDisks := []model.Disk{
		{"", "/dev/sda", "", []model.Partition{{"", "/dev/sda1", "", true, "btrfs", false}}, false, "", ""},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, &BtrfsDisksStub{}, log.Default())
	err := disks.ActivatePartition("/dev/sda1")
	assert.Nil(t, err)
}

func TestDisks_ActivateDisk_NotSupportedFs(t *testing.T) {

	allDisks := []model.Disk{
		{"", "/dev/sda", "", []model.Partition{{"", "/dev/sda1", "/", true, "fat32", false}}, false, "", ""},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, &BtrfsDisksStub{}, log.Default())
	err := disks.ActivatePartition("/dev/sda1")
	assert.NotNil(t, err)
}

func TestDisks_ActivateMultiDisk_None(t *testing.T) {
	executor := &StorageExecutorStub{}

	allDisks := []model.Disk{
		{"", "/dev/sda", "", []model.Partition{{"", "/dev/sda1", "/", true, "fat32", false}}, false, "", ""},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, executor, &BtrfsDisksStub{}, log.Default())
	err := disks.ActivateDisks([]string{}, false)
	assert.NotNil(t, err)

}

func TestDisks_ActivateMultiDisk_UseUuid(t *testing.T) {
	allDisks := []model.Disk{
		{"", "/dev/sda", "", []model.Partition{}, true, "uuid1", ""},
		{"", "/dev/sdb", "", []model.Partition{}, false, "uuid2", ""},
	}
	btrfs := &BtrfsDisksStub{}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, btrfs, log.Default())
	err := disks.ActivateDisks([]string{"/dev/sdb"}, true)
	assert.Nil(t, err)
	assert.Equal(t, "uuid2", btrfs.uuid)
	assert.Equal(t, []string{"/dev/sdb"}, btrfs.newDevices)

}

func TestDisks_ActivateMultiDisk_PartitionToDisk_Deactivate(t *testing.T) {
	allDisks := []model.Disk{
		{"", "/dev/sda", "", []model.Partition{{"", "/dev/sda1", "/", true, "fat32", false}}, false, "", ""},
	}
	btrfs := &BtrfsDisksStub{}
	systemd := &SystemdStub{}

	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, systemd, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, btrfs, log.Default())
	err := disks.ActivateDisks([]string{"/dev/sda"}, true)
	assert.Nil(t, err)
	assert.True(t, systemd.removeMountCalled)
	//assert.True(t, systemd.addMountCalled)
	assert.Equal(t, []string{"/dev/sda"}, btrfs.newDevices)

}

/*
func TestDisks_ActivateMultiDisk_DiskToDisk_NotDeactivate(t *testing.T) {
	allDisks := []model.Disk{
		{"", "/dev/sda", "", []model.Partition{}, true, "", ""},
	}
	btrfs := &BtrfsDisksStub{}
	systemd := &SystemdStub{}

	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, systemd, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, btrfs, log.Default())
	err := disks.ActivateDisks([]string{"/dev/sda"}, true)
	assert.Nil(t, err)
	assert.True(t, systemd.removeMountCalled)
	assert.True(t, systemd.addMountCalled)
	assert.Equal(t, []string{"/dev/sda"}, btrfs.updated)
}*/
