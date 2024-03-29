package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type FreeSpaceCheckerExecutorStub struct {
	output string
}

func (e *FreeSpaceCheckerExecutorStub) CombinedOutput(_ string, _ ...string) ([]byte, error) {
	return []byte(e.output), nil
}

func TestDisks_RootPartition_Low_NotExtendable(t *testing.T) {

	parted := `BYT;
/dev/sda:100%:scsi:512:512:msdos:ATA KINGSTON SV300S3:;
1:0.00%:0.00%:0.00%:free;
1:0.00%:1.67%:1.67%:::;
2:1.67%:100%:98.3%:ext4::;
1:100%:100%:0.00%:free;`

	checker := NewFreeSpaceChecker(&FreeSpaceCheckerExecutorStub{output: parted})
	hasFree, err := checker.HasFreeSpace("/dev/sda")
	assert.Nil(t, err)
	assert.False(t, hasFree)

}

func TestDisks_RootPartition_High_Extendable(t *testing.T) {

	parted := `BYT;
/dev/sda:100%:scsi:512:512:msdos:ATA KINGSTON SV300S3:;
1:0.00%:0.00%:0.00%:free;
1:0.00%:1.67%:1.67%:::;
2:1.67%:100%:98.3%:ext4::;
1:100%:100%:12.34%:free;`

	checker := NewFreeSpaceChecker(&FreeSpaceCheckerExecutorStub{output: parted})
	hasFree, err := checker.HasFreeSpace("/dev/sda")
	assert.Nil(t, err)
	assert.True(t, hasFree)
}

func TestDisks_RootPartition_No_NotExtendable(t *testing.T) {

	parted := `BYT;
/dev/sda:100%:scsi:512:512:msdos:ATA KINGSTON SV300S3:;
1:0.00%:0.00%:0.00%:free;
1:0.00%:1.67%:1.67%:::;
2:1.67%:100%:98.3%:ext4::;`
	checker := NewFreeSpaceChecker(&FreeSpaceCheckerExecutorStub{output: parted})
	hasFree, err := checker.HasFreeSpace("/dev/sda")
	assert.Nil(t, err)
	assert.False(t, hasFree)
}

func TestDisks_RootPartition_Extendable(t *testing.T) {

	parted := `
BYT;
/dev/mmcblk0:100%:sd/mmc:512:512:msdos:SD SD32G:;
1:0.00%:0.01%:0.01%:free;
1:0.01%:0.87%:0.86%:fat32::lba;
1:0.87%:0.87%:0.00%:free;
2:0.87%:11.2%:10.3%:ext4::;
1:11.2%:100%:88.8%:free;
`

	checker := NewFreeSpaceChecker(&FreeSpaceCheckerExecutorStub{output: parted})
	hasFree, err := checker.HasFreeSpace("/dev/sda")
	assert.Nil(t, err)
	assert.True(t, hasFree)
}
