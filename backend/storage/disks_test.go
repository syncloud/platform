package storage

import (
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
}

func (t TriggerStub) RunDiskChangeEvent() error {
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

func (l LsblkDisksStub) AvailableDisks() (*[]model.Disk, error) {
	return &l.disks, nil
}

func (l LsblkDisksStub) AllDisks() (*[]model.Disk, error) {
	return &l.disks, nil
}

type DisksExecutorStub struct {
	output string
}

func (e *DisksExecutorStub) CommandOutput(_ string, _ ...string) ([]byte, error) {
	return []byte(e.output), nil
}

type SystemdStub struct {
}

func (s SystemdStub) RemoveMount() error {
	return nil
}

func TestDisks_RootPartition_Low_NotExtendable(t *testing.T) {

	parted := `BYT;
/dev/sda:100%:scsi:512:512:msdos:ATA KINGSTON SV300S3:;
1:0.00%:0.00%:0.00%:free;
1:0.00%:1.67%:1.67%:::;
2:1.67%:100%:98.3%:ext4::;
1:100%:100%:0.00%:free;`
	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false, false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksExecutorStub{output: parted}, log.Default())
	partition, err := disks.RootPartition()
	assert.Nil(t, err)
	assert.False(t, partition.Extendable)

}

func TestDisks_RootPartition_High_Extendable(t *testing.T) {

	parted := `BYT;
/dev/sda:100%:scsi:512:512:msdos:ATA KINGSTON SV300S3:;
1:0.00%:0.00%:0.00%:free;
1:0.00%:1.67%:1.67%:::;
2:1.67%:100%:98.3%:ext4::;
1:100%:100%:12.34%:free;`
	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false, false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksExecutorStub{output: parted}, log.Default())
	partition, err := disks.RootPartition()
	assert.Nil(t, err)
	assert.True(t, partition.Extendable)
}

func TestDisks_RootPartition_No_NotExtendable(t *testing.T) {

	parted := `BYT;
/dev/sda:100%:scsi:512:512:msdos:ATA KINGSTON SV300S3:;
1:0.00%:0.00%:0.00%:free;
1:0.00%:1.67%:1.67%:::;
2:1.67%:100%:98.3%:ext4::;`
	allDisks := []model.Disk{
		{"", "", "", []model.Partition{{"", "", "/", true, "", false, false}}, false},
	}
	disks := NewDisks(&DisksConfigStub{}, &TriggerStub{}, &LsblkDisksStub{disks: allDisks}, &SystemdStub{}, &DisksExecutorStub{output: parted}, log.Default())
	partition, err := disks.RootPartition()
	assert.Nil(t, err)
	assert.False(t, partition.Extendable)
}
