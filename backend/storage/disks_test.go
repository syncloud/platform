package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/storage/model"
	"testing"
)

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
	error bool
}

func (t TriggerStub) RunDiskChangeEvent() error {
	if t.error {
		return fmt.Errorf("error")
	}
	return nil
}

//type PathCheckerStub struct {
//	exists bool
//}

//func (p *PathCheckerStub) ExternalDiskLinkExists() bool {
//	return p.exists
//}

type LsblkDisksStub struct {
	disks []model.Disk
}

func (l LsblkDisksStub) FindPartitionByDevice(_ string) (*model.Partition, error) {
	return &model.Partition{}, nil
}

func (l LsblkDisksStub) AvailableDisks() (*[]model.Disk, error) {
	return &l.disks, nil
}

func (l LsblkDisksStub) AllDisks() (*[]model.Disk, error) {
	return &l.disks, nil
}

//type DisksExecutorStub struct {
//	output string
//}

//func (e *DisksExecutorStub) CommandOutput(_ string, _ ...string) ([]byte, error) {
//	return []byte(e.output), nil
//}

type SystemdStub struct {
}

func (s SystemdStub) AddMount(_ string) error {
	return nil
}

func (s SystemdStub) RemoveMount() error {
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

func TestDisks_RootPartition_HasFreeSpace_Extendable(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false, false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{freeSpace: true}, &DisksLinkerStub{}, log.Default())
	partition, err := disks.RootPartition()
	assert.Nil(t, err)
	assert.True(t, partition.Extendable)
}

func TestDisks_RootPartition_HasNoFreeSpace_NonExtendable(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false, false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{freeSpace: false}, &DisksLinkerStub{}, log.Default())
	partition, err := disks.RootPartition()
	assert.Nil(t, err)
	assert.False(t, partition.Extendable)
}

func TestDisks_DeactivateDisk_TriggerError_NotFail(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false, false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: true}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, log.Default())
	err := disks.DeactivateDisk()
	assert.Nil(t, err)
}

func TestDisks_DeactivateDisk_TriggerNotError_NotFail(t *testing.T) {

	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false, false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{error: false}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksFreeSpaceCheckerStub{}, &DisksLinkerStub{}, log.Default())
	err := disks.DeactivateDisk()
	assert.Nil(t, err)
}
