package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/storage/model"
	"os"
	"testing"
)

func TestLsblk_AvailableDisks_HideDiskWithRootPartition(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/sdb", "size":"232.9G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"TOSHIBA MK2552GS", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/sdb1", "size":"100M", "type":"part", "mountpoint":null, "parttype":"0x7", "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/sdb2", "size":"184.1G", "type":"part", "mountpoint":"/opt/disk/external", "parttype":"0x83", "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      },
      {"name":"/dev/sdc", "size":"55.9G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"INTEL SSDSC2CW06", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/sdc1", "size":"48.5G", "type":"part", "mountpoint":"/", "parttype":"0x83", "fstype":"ntfs", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output: output}, log.Default())

	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, "/dev/sdb", disks[0].Device)
	assert.Equal(t, "/opt/disk/external", disks[0].Partitions[1].MountPoint)
	assert.Equal(t, 2, len(disks[0].Partitions))

}

func TestLsblk_AvailableDisks_LoopSupport(t *testing.T) {

	output := `{
   "blockdevices": [
      {"name":"/dev/loop0", "size":"10M", "type":"loop", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, 0, len(disks[0].Partitions))
	assert.Equal(t, "Disk loop0", disks[0].Name)
}

func TestLsblk_AvailableDisks_BlankDiskSupport_NotActive(t *testing.T) {

	output := `{
   "blockdevices": [
      {"name":"/dev/sdd", "size":"3.7G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"BLANK DISK", "uuid":null, "serial":null, "wwn":null, "hctl":null}
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, 0, len(disks[0].Partitions))
	assert.False(t, disks[0].Active)
}

func TestLsblk_AvailableDisks_BlankDiskSupport_Active(t *testing.T) {

	output := `{
   "blockdevices": [
      {"name":"/dev/sdd", "size":"3.7G", "type":"disk", "mountpoint":"/opt/disk/external", "parttype":null, "fstype":null, "model":"BLANK DISK", "uuid":null, "serial":null, "wwn":null, "hctl":null}
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, 0, len(disks[0].Partitions))
	assert.True(t, disks[0].Active)
}

func TestLsblk_AvailableDisks_EmptyDisk(t *testing.T) {

	output := `{
   "blockdevices": [
      {"name":"/dev/sda", "size":"3.7G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"BLANK DISK", "uuid":null, "serial":null, "wwn":null, "hctl":null}
   ]
}`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, 0, len(disks[0].Partitions))

}

func TestLsblk_AvailableDisks_DoNotShowSquashfs(t *testing.T) {

	output := `{
   "blockdevices": [
      {"name":"/dev/loop1", "size":"41.1M", "type":"loop", "mountpoint":"/snap/platform/180821", "parttype":null, "fstype":"squashfs", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
   ]
}`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(disks))
}

func TestLsblk_AvailableDisks_DoNotShowSnaps_BrokenFsTypeSquashfs(t *testing.T) {

	output := `{
   "blockdevices": [
      {"name":"/dev/loop1", "size":"41.1M", "type":"loop", "mountpoint":"/snap/platform/180821", "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
   ]
}`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(disks))
}

func TestLsblk_AvailableDisks_DoNotShowZram(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/zram0", "size":"2G", "type":"disk", "mountpoint":"[SWAP]", "parttype":null, "fstype":"swap", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
   ]
}`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(disks))
}

func TestLsblk_AvailableDisks_DoNotShowInternalDisks(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/mmcblk0", "size":"14.4G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/mmcblk0p2", "size":"14.4G", "type":"part", "mountpoint":"/", "parttype":"0x83", "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/mmcblk0p1", "size":"41.8M", "type":"part", "mountpoint":null, "parttype":"0xc", "fstype":"vfat", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(disks))
}

func TestLsblk_AvailableDisks_DoNotShowDisksWithRootPartition(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/sdb", "size":"14.4G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/sdb1", "size":"14.4G", "type":"part", "mountpoint":"/", "parttype":"0x83", "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/sdb2", "size":"41.8M", "type":"part", "mountpoint":null, "parttype":"0xc", "fstype":"vfat", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(disks))
}

func TestLsblk_AvailableDisks_DefaultEmptyDiskName(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/sdb", "size":"14.4G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/sdb1", "size":"14.4G", "type":"part", "mountpoint":null, "parttype":"0x83", "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/sdb2", "size":"41.8M", "type":"part", "mountpoint":null, "parttype":"0xc", "fstype":"vfat", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, "Disk sdb", disks[0].Name)
}

func TestLsblk_AvailableDisks_IsExternalDiskAttached(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/sdb", "size":"232.9G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"TOSHIBA MK2552GS", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/sdb1", "size":"100M", "type":"part", "mountpoint":null, "parttype":"0x7", "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/sdb2", "size":"48.7G", "type":"part", "mountpoint":null, "parttype":"0x7", "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/sdb3", "size":"184.1G", "type":"part", "mountpoint":"/opt/disk/external", "parttype":"0x83", "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      },
      {"name":"/dev/sdc", "size":"55.9G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"INTEL SSDSC2CW06", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/sdc1", "size":"48.5G", "type":"part", "mountpoint":"/", "parttype":"0x83", "fstype":"ntfs", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/sdc2", "size":"1K", "type":"part", "mountpoint":null, "parttype":"0x5", "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`

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
	output := `{
   "blockdevices": [
      {"name":"/dev/sdb", "size":"232.9G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"TOSHIBA MK2552GS", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/sdb1", "size":"100M", "type":"part", "mountpoint":null, "parttype":"0x7", "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/sdb2", "size":"48.7G", "type":"part", "mountpoint":null, "parttype":"0x7", "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/sdb3", "size":"184.1G", "type":"part", "mountpoint":"/opt/disk/external", "parttype":"0x83", "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`

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
	output := `{
   "blockdevices": [
      {"name":"/dev/sda", "size":"1.8T", "type":"disk", "mountpoint":null, "parttype":null, "fstype":"linux_raid_member", "model":"WDC WD20EFRX-68E", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/md0", "size":"3.7T", "type":"raid10", "mountpoint":"/mnt/md0", "parttype":null, "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      },
      {"name":"/dev/sdb", "size":"1.8T", "type":"disk", "mountpoint":null, "parttype":null, "fstype":"linux_raid_member", "model":"WDC WD20EFRX-68E", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/md0", "size":"3.7T", "type":"raid10", "mountpoint":"/mnt/md0", "parttype":null, "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      },
      {"name":"/dev/sdc", "size":"1.8T", "type":"disk", "mountpoint":null, "parttype":null, "fstype":"linux_raid_member", "model":"WDC WD20EFRX-68E", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/md0", "size":"3.7T", "type":"raid10", "mountpoint":"/mnt/md0", "parttype":null, "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      },
      {"name":"/dev/sdd", "size":"1.8T", "type":"disk", "mountpoint":null, "parttype":null, "fstype":"linux_raid_member", "model":"WDC WD20EFRX-68E", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/md0", "size":"3.7T", "type":"raid10", "mountpoint":"/mnt/md0", "parttype":null, "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/detached"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, 1, len(disks[0].Partitions))
	assert.Equal(t, "/mnt/md0", disks[0].Partitions[0].MountPoint)
}

func TestLsblk_AvailableDisks_UnsupportedDevicesWithPartitions(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/loop16", "size":"61.9M", "type":"loop", "mountpoint":null, "parttype":null, "fstype":"squashfs", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/loop16p1", "size":"953M", "type":"part", "mountpoint":null, "parttype":null, "fstype":"squashfs", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/loop16p2", "size":"3G", "type":"part", "mountpoint":null, "parttype":null, "fstype":"squashfs", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`

	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/detached"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(disks))
}

func TestLsblk_AllDisks_InvalidOutput(t *testing.T) {
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{"lsblk: unknown column"}, log.Default())
	_, err := lsblk.AllDisks()
	assert.ErrorContains(t, err, "unable to parse lsblk output")
}

func TestLsblk_FindPartitionByDevice_Found(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/sdb", "size":"232.9G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"TOSHIBA MK2552GS", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/sdb1", "size":"100M", "type":"part", "mountpoint":null, "parttype":"0x7", "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/sdb2", "size":"48.7G", "type":"part", "mountpoint":null, "parttype":"0x7", "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/sdb3", "size":"184.1G", "type":"part", "mountpoint":"/opt/disk/external", "parttype":"0x83", "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      },
      {"name":"/dev/sdc", "size":"55.9G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"INTEL SSDSC2CW06", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/sdc1", "size":"48.5G", "type":"part", "mountpoint":"/test", "parttype":"0x83", "fstype":"vfat", "model":" ", "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`

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
func TestLsblk_AvailableDisks_BtrfsSupport(t *testing.T) {

	output := `{
   "blockdevices": [
      {"name":"/dev/sda", "size":"1.8T", "type":"disk", "mountpoint":null, "parttype":null, "fstype":"btrfs", "model":"CT2000BX500SSD1", "uuid":null, "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/sdb", "size":"111.8G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":"btrfs", "model":"KINGSTON_SA400S37120G", "uuid":null, "serial":null, "wwn":null, "hctl":null}
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(disks))
}

func TestLsblk_AvailableDisks_BtrfsMultiDiskSupport(t *testing.T) {

	output := `{
   "blockdevices": [
      {"name":"/dev/sda", "size":"1.8T", "type":"disk", "mountpoint":"/opt/disk/external", "parttype":null, "fstype":"btrfs", "model":"CT2000BX500SSD1", "uuid":"1", "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/sdb", "size":"111.8G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":"btrfs", "model":"KINGSTON_SA400S37120G", "uuid":"1", "serial":null, "wwn":null, "hctl":null}
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(disks))
	assert.True(t, disks[0].Active)
	assert.True(t, disks[1].Active)
}

func TestLsblk_AvailableDisks_Sorted(t *testing.T) {

	output := `{
   "blockdevices": [
      {"name":"/dev/sdc", "size":"111.8G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":"btrfs", "model":"KINGSTON_SA400S37120G", "uuid":"1", "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/sda", "size":"1.8T", "type":"disk", "mountpoint":null, "parttype":null, "fstype":"btrfs", "model":"CT2000BX500SSD1", "uuid":"1", "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/sdb", "size":"111.8G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":"btrfs", "model":"KINGSTON_SA400S37120G", "uuid":"1", "serial":null, "wwn":null, "hctl":null}
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(disks))
	assert.Equal(t, "/dev/sda", disks[0].Device)
	assert.Equal(t, "/dev/sdb", disks[1].Device)
	assert.Equal(t, "/dev/sdc", disks[2].Device)

}

func TestLsblk_AllDisks_RootPartition(t *testing.T) {

	output := `{
   "blockdevices": [
      {"name":"/dev/sda", "size":"232.9G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"TOSHIBA MK2552GS", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/sda1", "size":"184.1G", "type":"part", "mountpoint":"/opt/disk/external", "parttype":"0x83", "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      },
      {"name":"/dev/sdb", "size":"55.9G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"INTEL SSDSC2CW06", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/sdb1", "size":"48.5G", "type":"part", "mountpoint":"/", "parttype":"0x83", "fstype":"ntfs", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`
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

	output := `{
   "blockdevices": [
      {"name":"/dev/mmcblk1", "size":"1.8G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/mmcblk1p1", "size":"1.8G", "type":"part", "mountpoint":"/", "parttype":"0x83", "fstype":"ext4", "model":null, "uuid":"08579268-0e8b-416b-981e-ed9867f213e4", "serial":null, "wwn":null, "hctl":null}
         ]
      },
      {"name":"/dev/mmcblk1boot0", "size":"1M", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/mmcblk1boot1", "size":"1M", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
   ]
}`
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

	output := `{
   "blockdevices": [
      {"name":"/dev/mmcblk1", "size":"14.6G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/mmcblk1boot0", "size":"4M", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/mmcblk1boot1", "size":"4M", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(disks))
}

func TestLsblk_AvailableDisks_NotBoot_NotHide(t *testing.T) {

	output := `{
   "blockdevices": [
      {"name":"/dev/mmcblk1", "size":"14.6G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/mmcblk1boo0", "size":"4M", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/mmcblk1boo1", "size":"4M", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(disks))
	assert.False(t, disks[0].Boot)
	assert.False(t, disks[1].Boot)
	assert.False(t, disks[2].Boot)
}

func TestLsblk_AllDisks_SerialWwnHctl(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/sda", "size":"1.8T", "type":"disk", "mountpoint":"/opt/disk/external", "parttype":null, "fstype":"btrfs", "model":"CT2000BX500SSD1", "uuid":"b7020451-ac9c-4195-a7e8-f072458b1bae", "serial":"2112E58EB6C6", "wwn":"0x500a0751e58eb6c6", "hctl":"0:0:0:0"},
      {"name":"/dev/sdb", "size":"931.5G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":"btrfs", "model":"External_USB_3.0", "uuid":"b7020451-ac9c-4195-a7e8-f072458b1bae", "serial":"20170617007242F", "wwn":null, "hctl":"2:0:0:0"},
      {"name":"/dev/mmcblk0", "size":"29.7G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":"0xbf138c98", "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/mmcblk0p2", "size":"29.5G", "type":"part", "mountpoint":"/", "parttype":"0x83", "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output: output}, log.Default())

	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(disks))
	assert.Equal(t, "CT2000BX500SSD1", disks[0].Name)
	assert.Equal(t, "2112E58EB6C6", disks[0].Serial)
	assert.Equal(t, "0x500a0751e58eb6c6", disks[0].Wwn)
	assert.Equal(t, "0:0:0:0", disks[0].Hctl)
	assert.True(t, disks[0].Active)
	assert.Equal(t, "20170617007242F", disks[1].Serial)
	assert.Equal(t, "2:0:0:0", disks[1].Hctl)
	assert.True(t, disks[1].Active)
}

func TestLsblk_AvailableDisks_Buster(t *testing.T) {
	output, err := os.ReadFile("testdata/lsblk/buster.json")
	assert.Nil(t, err)
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{string(output)}, log.Default())

	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, "/dev/sdb", disks[0].Device)
	assert.Equal(t, "External_USB_3.0", disks[0].Name)
	assert.Equal(t, "931.5G", disks[0].Size)
	assert.Equal(t, "SERIAL00000002", disks[0].Serial)
	assert.Equal(t, "", disks[0].Wwn)
	assert.Equal(t, "2:0:0:0", disks[0].Hctl)
	assert.Equal(t, 1, len(disks[0].Partitions))
	assert.Equal(t, "/dev/sdb1", disks[0].Partitions[0].Device)
	assert.Equal(t, "ext4", disks[0].Partitions[0].FsType)
	assert.True(t, disks[0].Partitions[0].Active)

	all, err := lsblk.AllDisks()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(all))
	assert.Equal(t, "/dev/sda2", all[0].FindRootPartition().Device)
}

func TestLsblk_AvailableDisks_Bookworm(t *testing.T) {
	output, err := os.ReadFile("testdata/lsblk/bookworm.json")
	assert.Nil(t, err)
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{string(output)}, log.Default())

	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(disks))
	assert.Equal(t, "/dev/loop3", disks[0].Device)
	assert.Equal(t, "/dev/loop4", disks[1].Device)
	disk := disks[2]
	assert.Equal(t, "/dev/sda", disk.Device)
	assert.Equal(t, "CT2000BX500SSD1", disk.Name)
	assert.Equal(t, "1.8T", disk.Size)
	assert.Equal(t, "SERIAL00000003", disk.Serial)
	assert.Equal(t, "0x5000000000000003", disk.Wwn)
	assert.Equal(t, "0:0:0:0", disk.Hctl)
	assert.Equal(t, "f85fb0ab-e643-4266-9f61-0f6e4980b871", disk.Uuid)
	assert.True(t, disk.Active)
	assert.Empty(t, disk.Partitions)

	all, err := lsblk.AllDisks()
	assert.Nil(t, err)
	assert.Equal(t, 4, len(all))
	assert.Equal(t, "/dev/mmcblk0", all[2].Device)
	assert.Equal(t, 2, len(all[2].Partitions))
	assert.Equal(t, "/dev/mmcblk0p2", all[2].FindRootPartition().Device)
}

func TestLsblk_AvailableDisks_SdCardNextToEmmcRoot(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/loop0", "size":"109M", "type":"loop", "mountpoint":"/snap/platform/1914", "parttype":null, "fstype":"squashfs", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/mmcblk1", "size":"3.6G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/mmcblk1p1", "size":"3.6G", "type":"part", "mountpoint":"/", "parttype":"0x83", "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      },
      {"name":"/dev/mmcblk1boot0", "size":"2M", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/mmcblk1boot1", "size":"2M", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/mmcblk0", "size":"3.8G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/mmcblk0p1", "size":"3.8G", "type":"part", "mountpoint":null, "parttype":"0x83", "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, "/dev/mmcblk0", disks[0].Device)
	assert.Equal(t, "/dev/mmcblk0p1", disks[0].Partitions[0].Device)
}

func TestLsblk_AvailableDisks_UnassembledMdadmMembers(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/loop0", "size":"164M", "type":"loop", "mountpoint":"/snap/platform/2261", "parttype":null, "fstype":"squashfs", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/sda", "size":"3.7T", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/sdb", "size":"3.7T", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
      {"name":"/dev/mmcblk0", "size":"59.5G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/mmcblk0p1", "size":"256M", "type":"part", "mountpoint":null, "parttype":"0x83", "fstype":"ext2", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/mmcblk0p2", "size":"59.2G", "type":"part", "mountpoint":"/", "parttype":"0x83", "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(disks))
	assert.Equal(t, "/dev/sda", disks[0].Device)
	assert.Equal(t, "/dev/sdb", disks[1].Device)
}

func TestLsblk_AvailableDisks_RootOnSataGptNextToWholeDiskBtrfs(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/sda", "size":"3.6T", "type":"disk", "mountpoint":null, "parttype":null, "fstype":"btrfs", "model":"One Touch HDD", "uuid":"9c1f0c4e-0000-4000-8000-000000000001", "serial":null, "wwn":null, "hctl":"6:0:0:0"},
      {"name":"/dev/sdb", "size":"298.1G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"TOSHIBA MK3255GS", "uuid":null, "serial":null, "wwn":null, "hctl":"0:0:0:0",
         "children": [
            {"name":"/dev/sdb1", "size":"512M", "type":"part", "mountpoint":"/boot/efi", "parttype":"c12a7328-f81f-11d2-ba4b-00a0c93ec93b", "fstype":"vfat", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/sdb2", "size":"296.6G", "type":"part", "mountpoint":"/", "parttype":"0fc63daf-8483-4772-8e79-3d69d8477de4", "fstype":"btrfs", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null},
            {"name":"/dev/sdb3", "size":"977M", "type":"part", "mountpoint":"[SWAP]", "parttype":"0657fd6d-a4ab-43c4-84e5-0933c84b4f4f", "fstype":"swap", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.Equal(t, "/dev/sda", disks[0].Device)
	assert.Equal(t, "One Touch HDD", disks[0].Name)
	assert.Empty(t, disks[0].Partitions)
}

func TestLsblk_AvailableDisks_GptPartitionWithGuidPartType(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/sda", "size":"923.6G", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"Expansion", "uuid":null, "serial":null, "wwn":null, "hctl":null,
         "children": [
            {"name":"/dev/sda1", "size":"923.5G", "type":"part", "mountpoint":"/opt/disk/external", "parttype":"ebd0a0a2-b9e5-4433-87c0-68b6b72699c7", "fstype":"ext3", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
         ]
      }
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(disks))
	assert.False(t, disks[0].Active)
	assert.Equal(t, 1, len(disks[0].Partitions))
	assert.True(t, disks[0].Partitions[0].Active)
	assert.Equal(t, "ext3", disks[0].Partitions[0].FsType)
}

func TestLsblk_AvailableDisks_AssembledMdadmOnPartitions(t *testing.T) {
	output := `{
   "blockdevices": [
      {"name":"/dev/sda", "size":"3.7T", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"ST4000VN008", "uuid":null, "serial":null, "wwn":null, "hctl":"0:0:0:0",
         "children": [
            {"name":"/dev/sda1", "size":"3.7T", "type":"part", "mountpoint":null, "parttype":"0xfd", "fstype":"linux_raid_member", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null,
               "children": [
                  {"name":"/dev/md0", "size":"3.6T", "type":"raid1", "mountpoint":"/opt/disk/external", "parttype":null, "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
               ]
            }
         ]
      },
      {"name":"/dev/sdb", "size":"3.7T", "type":"disk", "mountpoint":null, "parttype":null, "fstype":null, "model":"ST4000VN008", "uuid":null, "serial":null, "wwn":null, "hctl":"1:0:0:0",
         "children": [
            {"name":"/dev/sdb1", "size":"3.7T", "type":"part", "mountpoint":null, "parttype":"0xfd", "fstype":"linux_raid_member", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null,
               "children": [
                  {"name":"/dev/md0", "size":"3.6T", "type":"raid1", "mountpoint":"/opt/disk/external", "parttype":null, "fstype":"ext4", "model":null, "uuid":null, "serial":null, "wwn":null, "hctl":null}
               ]
            }
         ]
      }
   ]
}`
	lsblk := NewLsblk(&ConfigStub{diskDir: "/opt/disk/external"}, &PathCheckerStub{exists: true}, &ExecutorStub{output}, log.Default())
	disks, err := lsblk.AvailableDisks()
	assert.Nil(t, err)
	var md model.Disk
	for _, disk := range disks {
		if disk.Device == "/dev/md0" {
			md = disk
		}
	}
	assert.Equal(t, "raid1", md.Name)
	assert.Equal(t, 1, len(md.Partitions))
	assert.True(t, md.Partitions[0].Active)
	assert.Equal(t, "raid", md.Partitions[0].FsType)
}
