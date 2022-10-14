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

func (l *LsblkDisksStub) AllDisks() (*[]model.Disk, error) {
	return &l.disks, nil
}

type SystemdStub struct {
	callOrderShared *CallOrder
	callOrder       int
}

func (s *SystemdStub) AddMount(_ string) error {
	return nil
}

func (s *SystemdStub) RemoveMount() error {
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

func TestDisks_RootPartition_HasFreeSpace_Extendable(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{freeSpace: true}, &DisksLinkerStub{}, &StorageExecutorStub{}, log.Default())
	partition, err := disks.RootPartition()
	assert.Nil(t, err)
	assert.True(t, partition.Extendable)
}

func TestDisks_RootPartition_HasNoFreeSpace_NonExtendable(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{freeSpace: false}, &DisksLinkerStub{}, &StorageExecutorStub{}, log.Default())
	partition, err := disks.RootPartition()
	assert.Nil(t, err)
	assert.False(t, partition.Extendable)
}

func TestDisks_DeactivateDisk_TriggerError_NotFail(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: true}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, log.Default())
	err := disks.DeactivateDisk()
	assert.Nil(t, err)
}

func TestDisks_DeactivateDisk_TriggerNotError_NotFail(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, log.Default())
	err := disks.DeactivateDisk()
	assert.Nil(t, err)
}

func TestDisks_DeactivateDisk_TriggerEventBeforeRemove(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false}}, false},
	}
	callOrder := &CallOrder{order: 0}
	trigger := &TriggerStub{error: false, callOrderShared: callOrder}
	systemd := &SystemdStub{callOrderShared: callOrder}
	disks := NewDisks(&DisksConfigStub{}, trigger, &LsblkDisksStub{disks: allDisks}, systemd, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, log.Default())
	err := disks.DeactivateDisk()
	assert.Nil(t, err)
	assert.Less(t, trigger.callOrder, systemd.callOrder)

}

func TestDisks_ActivateDisk_SupportedFs(t *testing.T) {

	allDisks := []model.Disk{
		{"", "/dev/sda", "", []model.Partition{{"", "/dev/sda1", "/", true, "ext4", false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, log.Default())
	err := disks.ActivateDisk("/dev/sda1")
	assert.Nil(t, err)
}

func TestDisks_ActivateDisk_Btrfs(t *testing.T) {

	allDisks := []model.Disk{
		{"", "/dev/sda", "", []model.Partition{{"", "/dev/sda1", "", true, "btrfs", false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, log.Default())
	err := disks.ActivateDisk("/dev/sda1")
	assert.Nil(t, err)
}

func TestDisks_ActivateDisk_NotSupportedFs(t *testing.T) {

	allDisks := []model.Disk{
		{"", "/dev/sda", "", []model.Partition{{"", "/dev/sda1", "/", true, "fat32", false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, &StorageExecutorStub{}, log.Default())
	err := disks.ActivateDisk("/dev/sda1")
	assert.NotNil(t, err)
}

func TestDisks_ActivateMultiDisk_None(t *testing.T) {
	executor := &StorageExecutorStub{}

	allDisks := []model.Disk{
		{"", "/dev/sda", "", []model.Partition{{"", "/dev/sda1", "/", true, "fat32", false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, executor, log.Default())
	err := disks.ActivateMultiDisk([]string{})
	assert.NotNil(t, err)

}

func TestDisks_ActivateMultiDisk_Single(t *testing.T) {
	executor := &StorageExecutorStub{}

	allDisks := []model.Disk{
		{"", "/dev/sda", "", []model.Partition{{"", "/dev/sda1", "/", true, "fat32", false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, executor, log.Default())
	err := disks.ActivateMultiDisk([]string{"/dev/sda"})
	assert.Nil(t, err)
	assert.Equal(t, "single", executor.args[6])
	assert.Equal(t, "/dev/sda", executor.args[7])

}

func TestDisks_ActivateMultiDisk_Multiple(t *testing.T) {
	executor := &StorageExecutorStub{}

	allDisks := []model.Disk{
		{"", "/dev/sda", "", []model.Partition{{"", "/dev/sda1", "/", true, "fat32", false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, executor, log.Default())
	err := disks.ActivateMultiDisk([]string{"/dev/sda", "/dev/sdb"})
	assert.Nil(t, err)
	assert.Equal(t, "raid1", executor.args[6])
	assert.Equal(t, "/dev/sda", executor.args[7])
	assert.Equal(t, "/dev/sdb", executor.args[8])

}
