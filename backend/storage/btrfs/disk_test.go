package btrfs

import (
	"github.com/prometheus/procfs/btrfs"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"testing"
)

type ConfigStub struct {
}

func (c *ConfigStub) ExternalDiskDir() string {
	return "/mnt"
}

type ExecutorStub struct {
	command string
	args    []string
}

func (e *ExecutorStub) CommandOutput(command string, args ...string) ([]byte, error) {
	e.command = command
	e.args = args
	return []byte(""), nil
}

type StatsStub struct {
}

func (s *StatsStub) Stats() ([]*btrfs.Stats, error) {
	//TODO implement me
	panic("implement me")
}

func Test_DetectChange_Replaced(t *testing.T) {
	disks := &Disks{&ConfigStub{}, &ExecutorStub{}, &StatsStub{}, log.Default()}
	changes, err := disks.Apply([]string{"/dev/loop1", "/dev/loop2"}, []string{"/dev/loop1", "/dev/loop3"}, "")

	assert.Nil(t, err)
	assert.Len(t, changes, 3)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device add /dev/loop3 /mnt", changes[0].ToString())
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh balance start -dconvert=raid1 -mconvert=raid1 /mnt", changes[1].ToString())
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device delete /dev/loop2 /mnt", changes[2].ToString())
}

func Test_DetectChange_Add_One(t *testing.T) {
	disks := &Disks{&ConfigStub{}, &ExecutorStub{}, &StatsStub{}, log.Default()}
	changes, err := disks.Apply([]string{"/dev/loop1"}, []string{"/dev/loop1", "/dev/loop2"}, "")

	assert.Nil(t, err)
	assert.Len(t, changes, 2)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device add /dev/loop2 /mnt", changes[0].ToString())
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh balance start -dconvert=raid1 -mconvert=raid1 /mnt", changes[1].ToString())
}

func Test_DetectChange_Create_One(t *testing.T) {
	disks := &Disks{&ConfigStub{}, &ExecutorStub{}, &StatsStub{}, log.Default()}
	changes, err := disks.Apply([]string{}, []string{"/dev/loop1"}, "uuid")

	assert.Nil(t, err)
	assert.Len(t, changes, 1)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/mkfs.sh -U uuid -f -m single -d single /dev/loop1", changes[0].ToString())
}

func Test_DetectChange_Create_Two(t *testing.T) {
	disks := &Disks{&ConfigStub{}, &ExecutorStub{}, &StatsStub{}, log.Default()}
	changes, err := disks.Apply([]string{}, []string{"/dev/loop1", "/dev/loop2"}, "uuid")

	assert.Nil(t, err)
	assert.Len(t, changes, 1)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/mkfs.sh -U uuid -f -m raid1 -d raid1 /dev/loop1 /dev/loop2", changes[0].ToString())
}

func Test_DetectChange_Remove_One(t *testing.T) {
	disks := &Disks{&ConfigStub{}, &ExecutorStub{}, &StatsStub{}, log.Default()}
	changes, err := disks.Apply([]string{"/dev/loop1", "/dev/loop2"}, []string{"/dev/loop1"}, "")

	assert.Nil(t, err)
	assert.Len(t, changes, 2)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh balance start -dconvert=single -mconvert=single /mnt", changes[0].ToString())
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device delete /dev/loop2 /mnt", changes[1].ToString())

}

func Test_DetectChange_Disable_One(t *testing.T) {
	disks := &Disks{&ConfigStub{}, &ExecutorStub{}, &StatsStub{}, log.Default()}
	changes, err := disks.Apply([]string{"/dev/loop1"}, []string{}, "")

	assert.Nil(t, err)
	assert.Len(t, changes, 0)
}

func Test_DetectChange_Disable_Two(t *testing.T) {
	disks := &Disks{&ConfigStub{}, &ExecutorStub{}, &StatsStub{}, log.Default()}
	changes, err := disks.Apply([]string{"/dev/loop1", "/dev/loop2"}, []string{}, "")

	assert.Nil(t, err)
	assert.Len(t, changes, 0)
}

func Test_DetectChange_0_To_3(t *testing.T) {
	disks := &Disks{&ConfigStub{}, &ExecutorStub{}, &StatsStub{}, log.Default()}
	changes, err := disks.Apply([]string{}, []string{"/dev/loop1", "/dev/loop2", "/dev/loop3"}, "uuid")

	assert.Nil(t, err)
	assert.Len(t, changes, 1)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/mkfs.sh -U uuid -f -m raid10 -d raid10 /dev/loop1 /dev/loop2 /dev/loop3", changes[0].ToString())
}

func Test_DetectChange_1_To_3(t *testing.T) {
	disks := &Disks{&ConfigStub{}, &ExecutorStub{}, &StatsStub{}, log.Default()}
	changes, err := disks.Apply([]string{"/dev/loop1"}, []string{"/dev/loop1", "/dev/loop2", "/dev/loop3"}, "uuid")

	assert.Nil(t, err)
	assert.Len(t, changes, 2)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device add /dev/loop2 /dev/loop3 /mnt", changes[0].ToString())
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh balance start -dconvert=raid10 -mconvert=raid10 /mnt", changes[1].ToString())
}
