package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/storage/model"
)

type FileSystemStub struct {
	stats map[string][3]uint64
	err   error
}

func (f *FileSystemStub) Stat(path string) (uint64, uint64, uint64, error) {
	if f.err != nil {
		return 0, 0, 0, f.err
	}
	stat, found := f.stats[path]
	if !found {
		return 0, 0, 0, fmt.Errorf("no such path: %s", path)
	}
	return stat[0], stat[1], stat[2], nil
}

type DiskSpaceConfigStub struct {
}

func (c *DiskSpaceConfigStub) DiskLink() string {
	return "/data"
}

func diskSpace(stats map[string][3]uint64, err error) *DiskSpace {
	return NewDiskSpace(&FileSystemStub{stats: stats, err: err}, &DiskSpaceConfigStub{}, log.Default())
}

func TestDiskSpace_NotLowWhenEnoughFree(t *testing.T) {
	status := diskSpace(map[string][3]uint64{
		"/":     {15 * 1024 * 1024, 8 * 1024 * 1024, 1},
		"/data": {900 * 1024 * 1024, 400 * 1024 * 1024, 2},
	}, nil).Status()

	assert.False(t, status.Low)
	assert.Len(t, status.Mounts, 2)
	assert.False(t, status.Mounts[0].Low)
}

func TestDiskSpace_LowOnSystemDisk(t *testing.T) {
	status := diskSpace(map[string][3]uint64{
		"/":     {15 * 1024 * 1024, 855 * 1024, 1},
		"/data": {900 * 1024 * 1024, 400 * 1024 * 1024, 2},
	}, nil).Status()

	assert.True(t, status.Low)
	assert.True(t, status.Mounts[0].Low)
	assert.False(t, status.Mounts[1].Low)
}

func TestDiskSpace_LowOnDataDisk(t *testing.T) {
	status := diskSpace(map[string][3]uint64{
		"/":     {15 * 1024 * 1024, 8 * 1024 * 1024, 1},
		"/data": {900 * 1024 * 1024, 1024 * 1024, 2},
	}, nil).Status()

	assert.True(t, status.Low)
	assert.True(t, status.Mounts[1].Low)
	assert.Equal(t, model.DiskSpaceData, status.Mounts[1].Kind)
}

func TestDiskSpace_SameDeviceReportedOnce(t *testing.T) {
	status := diskSpace(map[string][3]uint64{
		"/":     {15 * 1024 * 1024, 8 * 1024 * 1024, 1},
		"/data": {15 * 1024 * 1024, 8 * 1024 * 1024, 1},
	}, nil).Status()

	assert.Len(t, status.Mounts, 1)
	assert.Equal(t, "/", status.Mounts[0].Path)
	assert.Equal(t, model.DiskSpaceSystem, status.Mounts[0].Kind)
}

func TestDiskSpace_MissingPathIsSkipped(t *testing.T) {
	status := diskSpace(map[string][3]uint64{
		"/": {15 * 1024 * 1024, 8 * 1024 * 1024, 1},
	}, nil).Status()

	assert.False(t, status.Low)
	assert.Len(t, status.Mounts, 1)
}

func TestDiskSpace_StatErrorIsNotLow(t *testing.T) {
	status := diskSpace(nil, fmt.Errorf("no such file")).Status()

	assert.False(t, status.Low)
	assert.Empty(t, status.Mounts)
}
