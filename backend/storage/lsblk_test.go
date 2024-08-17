package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/storage/model"
	"testing"
)

type ConfigStub struct {
	diskDir string
}

func (c *ConfigStub) ExternalDiskDir() string {
	return c.diskDir
}

type PathCheckerStub struct {
	exists bool
}

func (p *PathCheckerStub) ExternalDiskLinkExists() bool {
	return p.exists
}

type ExecutorStub struct {
	output string
}

func (e *ExecutorStub) CombinedOutput(_ string, _ ...string) ([]byte, error) {
	return []byte(e.output), nil
}

func TestLsblk_AvailableDisks_HideDiskWithRootPartition(t *testing.T) {
	output := `
NAME="/dev/sdb" SIZE="232.9G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="TOSHIBA MK2552GS" UUID=""
NAME="/dev/sdb1" SIZE="100M" TYPE="part" MOUNTPOINT="" PARTTYPE="0x7" FSTYPE="" MODEL="" UUID=""
NAME="/dev/sdb2" SIZE="184.1G" TYPE="part" MOUNTPOINT="/opt/disk/external" PARTTYPE="0x83" FSTYPE="ext4" MODEL="" UUID=""
NAME="/dev/sdc" SIZE="55.9G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="INTEL SSDSC2CW06" UUID=""
NAME="/dev/sdc1" SIZE="48.5G" TYPE="part" MOUNTPOINT="/" PARTTYPE="0x83" FSTYPE="ntfs" MODEL="" UUID=""
`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output: output}, log.Default())

	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, "/dev/sdb", disks[0].Device)
	assert.Equal(t, "/opt/disk/external", disks[0].Partitions[1].MountPoint)
	assert.Equal(t, 2, len(disks[0].Partitions))

}

func TestLsblk_AvailableDisks_LoopSupport(t *testing.T) {

	output := `NAME="/dev/loop0" SIZE="10M" TYPE="loop" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, 0, len(disks[0].Partitions))
	assert.Equal(t, "Disk loop0", disks[0].Name)
}

func TestLsblk_AvailableDisks_BlankDiskSupport_NotActive(t *testing.T) {

	output := `NAME="/dev/sdd" SIZE="3.7G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="BLANK DISK" UUID=""`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, 0, len(disks[0].Partitions))
	assert.False(t, disks[0].Active)
}

func TestLsblk_AvailableDisks_BlankDiskSupport_Active(t *testing.T) {

	output := `NAME="/dev/sdd" SIZE="3.7G" TYPE="disk" MOUNTPOINT="/opt/disk/external" PARTTYPE="" FSTYPE="" MODEL="BLANK DISK" UUID=""`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, 0, len(disks[0].Partitions))
	assert.True(t, disks[0].Active)
}

func TestLsblk_AvailableDisks_EmptyDisk(t *testing.T) {

	output := `NAME="/dev/sda" SIZE="3.7G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="BLANK DISK" UUID=""`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, 0, len(disks[0].Partitions))

}

func TestLsblk_AvailableDisks_DoNotShowSquashfs(t *testing.T) {

	output := `NAME="/dev/loop1" SIZE="41.1M" TYPE="loop" MOUNTPOINT="/snap/platform/180821" PARTTYPE="" FSTYPE="squashfs" MODEL="" UUID=""`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(disks))
}

func TestLsblk_AvailableDisks_DoNotShowSnaps_BrokenFsTypeSquashfs(t *testing.T) {

	output := `NAME="/dev/loop1" SIZE="41.1M" TYPE="loop" MOUNTPOINT="/snap/platform/180821" PARTTYPE="" FSTYPE="" MODEL="" UUID=""`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(disks))
}

func TestLsblk_AvailableDisks_DoNotShowInternalDisks(t *testing.T) {
	output := `
NAME="/dev/mmcblk0" SIZE="14.4G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""
NAME="/dev/mmcblk0p2" SIZE="14.4G" TYPE="part" MOUNTPOINT="/" PARTTYPE="0x83" FSTYPE="ext4" MODEL="" UUID=""
NAME="/dev/mmcblk0p1" SIZE="41.8M" TYPE="part" MOUNTPOINT="" PARTTYPE="0xc" FSTYPE="vfat" MODEL="" UUID=""
`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(disks))
}

func TestLsblk_AvailableDisks_DoNotShowDisksWithRootPartition(t *testing.T) {
	output := `
NAME="/dev/sdb" SIZE="14.4G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""
NAME="/dev/sdb1" SIZE="14.4G" TYPE="part" MOUNTPOINT="/" PARTTYPE="0x83" FSTYPE="ext4" MODEL="" UUID=""
NAME="/dev/sdb2" SIZE="41.8M" TYPE="part" MOUNTPOINT="" PARTTYPE="0xc" FSTYPE="vfat" MODEL="" UUID=""
`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(disks))
}

func TestLsblk_AvailableDisks_DefaultEmptyDiskName(t *testing.T) {
	output := `
NAME="/dev/sdb" SIZE="14.4G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""
NAME="/dev/sdb1" SIZE="14.4G" TYPE="part" MOUNTPOINT="" PARTTYPE="0x83" FSTYPE="ext4" MODEL="" UUID=""
NAME="/dev/sdb2" SIZE="41.8M" TYPE="part" MOUNTPOINT="" PARTTYPE="0xc" FSTYPE="vfat" MODEL="" UUID=""
`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, "Disk sdb", disks[0].Name)
}

func TestLsblk_AvailableDisks_IsExternalDiskAttached(t *testing.T) {
	output := `
NAME="/dev/sdb" SIZE="232.9G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="TOSHIBA MK2552GS" UUID=""
NAME="/dev/sdb1" SIZE="100M" TYPE="part" MOUNTPOINT="" PARTTYPE="0x7" FSTYPE="" MODEL="" UUID=""
NAME="/dev/sdb2" SIZE="48.7G" TYPE="part" MOUNTPOINT="" PARTTYPE="0x7" FSTYPE="" MODEL="" UUID=""
NAME="/dev/sdb3" SIZE="184.1G" TYPE="part" MOUNTPOINT="/opt/disk/external" PARTTYPE="0x83" FSTYPE="ext4" MODEL="" UUID=""
NAME="/dev/sdc" SIZE="55.9G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="INTEL SSDSC2CW06" UUID=""
NAME="/dev/sdc1" SIZE="48.5G" TYPE="part" MOUNTPOINT="/" PARTTYPE="0x83" FSTYPE="ntfs" MODEL="" UUID=""
NAME="/dev/sdc2" SIZE="1K" TYPE="part" MOUNTPOINT="" PARTTYPE="0x5" FSTYPE="" MODEL="" UUID=""
`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output: output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	var disk model.Disk
	for _, d := range disks {
		if d.Device == "/dev/sdb" {
			disk = d
		}
	}
	assert.False(t, disk.Active)

	var partition model.Partition
	for _, p := range disk.Partitions {
		if p.Device == "/dev/sdb3" {
			partition = p
		}
	}
	assert.True(t, partition.Active)
}

func TestLsblk_AvailableDisks_IsExternalDiskDetached(t *testing.T) {
	output := `
NAME="/dev/sdb" SIZE="232.9G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="TOSHIBA MK2552GS" UUID=""
NAME="/dev/sdb1" SIZE="100M" TYPE="part" MOUNTPOINT="" PARTTYPE="0x7" FSTYPE="" MODEL="" UUID=""
NAME="/dev/sdb2" SIZE="48.7G" TYPE="part" MOUNTPOINT="" PARTTYPE="0x7" FSTYPE="" MODEL="" UUID=""
NAME="/dev/sdb3" SIZE="184.1G" TYPE="part" MOUNTPOINT="/opt/disk/external" PARTTYPE="0x83" FSTYPE="ext4" MODEL="" UUID=""
`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/detached"}, &PathCheckerStub{exists: true}, &ExecutorStub{output: output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	var disk model.Disk
	for _, d := range disks {
		if d.Device == "/dev/sdb" {
			disk = d
		}
	}
	assert.False(t, disk.Active)
}

func TestLsblk_AvailableDisks_Raid(t *testing.T) {
	output := `
NAME="/dev/sda" SIZE="1.8T" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="linux_raid_member" MODEL="WDC WD20EFRX-68E" UUID=""
NAME="/dev/sdb" SIZE="1.8T" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="linux_raid_member" MODEL="WDC WD20EFRX-68E" UUID=""
NAME="/dev/sdc" SIZE="1.8T" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="linux_raid_member" MODEL="WDC WD20EFRX-68E" UUID=""
NAME="/dev/sdd" SIZE="1.8T" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="linux_raid_member" MODEL="WDC WD20EFRX-68E" UUID=""
NAME="/dev/md0" SIZE="3.7T" TYPE="raid10" MOUNTPOINT="/mnt/md0" PARTTYPE="" FSTYPE="ext4" MODEL="" UUID=""
NAME="/dev/md0" SIZE="3.7T" TYPE="raid10" MOUNTPOINT="/mnt/md0" PARTTYPE="" FSTYPE="ext4" MODEL="" UUID=""
NAME="/dev/md0" SIZE="3.7T" TYPE="raid10" MOUNTPOINT="/mnt/md0" PARTTYPE="" FSTYPE="ext4" MODEL="" UUID=""
NAME="/dev/md0" SIZE="3.7T" TYPE="raid10" MOUNTPOINT="/mnt/md0" PARTTYPE="" FSTYPE="ext4" MODEL="" UUID=""
`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/detached"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, 1, len(disks[0].Partitions))
	assert.Equal(t, "/mnt/md0", disks[0].Partitions[0].MountPoint)
}

func TestLsblk_AvailableDisks_UnsupportedDevicesWithPartitions(t *testing.T) {
	output := `
NAME="/dev/loop16" SIZE="61.9M" TYPE="loop" MOUNTPOINT="" PARTTYPE="" FSTYPE="squashfs" MODEL="" UUID=""
NAME="/dev/loop16p1" SIZE="953M" TYPE="part" MOUNTPOINT="" PARTTYPE="" FSTYPE="squashfs" MODEL="" UUID=""
NAME="/dev/loop16p2" SIZE="3G" TYPE="part" MOUNTPOINT="" PARTTYPE="" FSTYPE="squashfs" MODEL="" UUID=""
`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/detached"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(disks))
}

func TestLsblk_AvailableDisks_SkipEmptyLines(t *testing.T) {
	output := `

NAME="/dev/loop0" SIZE="10M" TYPE="loop" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""

`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
}

/*
	func TestLsblk_FindPartitionByDevice_Found(t *testing.T) {
		output := `

NAME="/dev/sdb" SIZE="232.9G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="TOSHIBA MK2552GS" UUID=""
NAME="/dev/sdb1" SIZE="100M" TYPE="part" MOUNTPOINT="" PARTTYPE="0x7" FSTYPE="" MODEL="" UUID=""
NAME="/dev/sdb2" SIZE="48.7G" TYPE="part" MOUNTPOINT="" PARTTYPE="0x7" FSTYPE="" MODEL="" UUID=""
NAME="/dev/sdb3" SIZE="184.1G" TYPE="part" MOUNTPOINT="/opt/disk/external" PARTTYPE="0x83" FSTYPE="ext4" MODEL="" UUID=""
NAME="/dev/sdc" SIZE="55.9G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="INTEL SSDSC2CW06" UUID=""
NAME="/dev/sdc1" SIZE="48.5G" TYPE="part" MOUNTPOINT="/test" PARTTYPE="0x83" FSTYPE="vfat" MODEL=" " UUID=""
`

		lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
		partition, err := lsblk.FindPartitionByDevice("/dev/sdc1")
		assert.Nil(t, err)
		assert.NotNil(t, partition)
		assert.Equal(t, "/dev/sdc1", partition.Device)
		assert.Equal(t, "/test", partition.MountPoint)
	}

	func TestLsblk_FindPartitionByDevice_NotFound(t *testing.T) {
		output := ""

		lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
		_, err := lsblk.FindPartitionByDevice("/dev/sdc1")
		assert.NotNil(t, err)
	}
*/
func TestLsblk_AvailableDisks_BtrfsSupport(t *testing.T) {

	output := `
NAME="/dev/sda" SIZE="1.8T" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="btrfs" MODEL="CT2000BX500SSD1" UUID=""
NAME="/dev/sdb" SIZE="111.8G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="btrfs" MODEL="KINGSTON_SA400S37120G" UUID=""
`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(disks))
}

func TestLsblk_AvailableDisks_BtrfsMultiDiskSupport(t *testing.T) {

	output := `
NAME="/dev/sda" SIZE="1.8T" TYPE="disk" MOUNTPOINT="/opt/disk/external" PARTTYPE="" FSTYPE="btrfs" MODEL="CT2000BX500SSD1" UUID="1"
NAME="/dev/sdb" SIZE="111.8G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="btrfs" MODEL="KINGSTON_SA400S37120G" UUID="1"
`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(disks))
	assert.True(t, disks[0].Active)
	assert.True(t, disks[1].Active)
}

func TestLsblk_AvailableDisks_Sorted(t *testing.T) {

	output := `
NAME="/dev/sdc" SIZE="111.8G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="btrfs" MODEL="KINGSTON_SA400S37120G" UUID="1"
NAME="/dev/sda" SIZE="1.8T" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="btrfs" MODEL="CT2000BX500SSD1" UUID="1"
NAME="/dev/sdb" SIZE="111.8G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="btrfs" MODEL="KINGSTON_SA400S37120G" UUID="1"
`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(disks))
	assert.Equal(t, "/dev/sda", disks[0].Device)
	assert.Equal(t, "/dev/sdb", disks[1].Device)
	assert.Equal(t, "/dev/sdc", disks[2].Device)

}

func TestLsblk_AllDisks_RootPartition(t *testing.T) {

	output := `
NAME="/dev/sda" SIZE="232.9G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="TOSHIBA MK2552GS" UUID=""
NAME="/dev/sda1" SIZE="184.1G" TYPE="part" MOUNTPOINT="/opt/disk/external" PARTTYPE="0x83" FSTYPE="ext4" MODEL="" UUID=""
NAME="/dev/sdb" SIZE="55.9G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="INTEL SSDSC2CW06" UUID=""
NAME="/dev/sdb1" SIZE="48.5G" TYPE="part" MOUNTPOINT="/" PARTTYPE="0x83" FSTYPE="ntfs" MODEL="" UUID=""
`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AllDisks()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(disks))
	assert.Equal(t, "/dev/sda", disks[0].Device)
	assert.False(t, disks[0].HasRootPartition())
	assert.Equal(t, "/dev/sdb", disks[1].Device)
	assert.True(t, disks[1].HasRootPartition())
	assert.Equal(t, "/dev/sdb1", disks[1].FindRootPartition().Device)
}

func TestLsblk_AllDisks_MMCBootPartitions_HidePartition(t *testing.T) {

	output := `
NAME="/dev/mmcblk1" SIZE="1.8G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""
NAME="/dev/mmcblk1p1" SIZE="1.8G" TYPE="part" MOUNTPOINT="/" PARTTYPE="0x83" FSTYPE="ext4" MODEL="" UUID="08579268-0e8b-416b-981e-ed9867f213e4"
NAME="/dev/mmcblk1boot0" SIZE="1M" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""
NAME="/dev/mmcblk1boot1" SIZE="1M" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""
`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AllDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, "/dev/mmcblk1", disks[0].Device)
	assert.True(t, disks[0].HasRootPartition())
	assert.True(t, disks[0].Boot)
	assert.Equal(t, 1, len(disks[0].Partitions))
}

func TestLsblk_AvailableDisks_HideBoot(t *testing.T) {

	output := `
NAME="/dev/mmcblk1" SIZE="14.6G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""
NAME="/dev/mmcblk1boot0" SIZE="4M" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""
NAME="/dev/mmcblk1boot1" SIZE="4M" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""
`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(disks))
}

func TestLsblk_AvailableDisks_NotBoot_NotHide(t *testing.T) {

	output := `
NAME="/dev/mmcblk1" SIZE="14.6G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""
NAME="/dev/mmcblk1boo0" SIZE="4M" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""
NAME="/dev/mmcblk1boo1" SIZE="4M" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="" UUID=""
`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(disks))
	assert.False(t, disks[0].Boot)
	assert.False(t, disks[1].Boot)
	assert.False(t, disks[2].Boot)
}
