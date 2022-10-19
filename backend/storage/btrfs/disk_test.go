package btrfs

import (
	"fmt"
	"github.com/prometheus/procfs/btrfs"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"strings"
	"testing"
)

type ConfigStub struct {
}

func (c *ConfigStub) ExternalDiskDir() string {
	return "/mnt"
}

type ExecutorStub struct {
	commands []string
}

func (e *ExecutorStub) CommandOutput(command string, args ...string) ([]byte, error) {
	e.commands = append(e.commands, fmt.Sprintf("%s %s", command, strings.Join(args, " ")))
	return []byte(""), nil
}

type StatsStub struct {
}

func (s *StatsStub) Stats() ([]*btrfs.Stats, error) {
	//TODO implement me
	panic("implement me")
}

type SystemdStub struct {
	addMountCalled bool
}

func (s *SystemdStub) AddMount(_ string) error {
	s.addMountCalled = true
	return nil
}

func Test_DetectChange_Replaced(t *testing.T) {
	executor := &ExecutorStub{}
	disks := &Disks{&ConfigStub{}, executor, &StatsStub{}, &SystemdStub{}, log.Default()}
	err := disks.Apply([]string{"/dev/loop1", "/dev/loop2"}, []string{"/dev/loop1", "/dev/loop3"}, "", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 3)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device add /dev/loop3 /mnt", executor.commands[0])
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh balance start -dconvert=raid1 -mconvert=raid1 /mnt", executor.commands[1])
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device delete /dev/loop2 /mnt", executor.commands[2])
}

func Test_DetectChange_Add_One(t *testing.T) {
	executor := &ExecutorStub{}
	disks := &Disks{&ConfigStub{}, executor, &StatsStub{}, &SystemdStub{}, log.Default()}
	err := disks.Apply([]string{"/dev/loop1"}, []string{"/dev/loop1", "/dev/loop2"}, "", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 2)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device add /dev/loop2 /mnt", executor.commands[0])
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh balance start -dconvert=raid1 -mconvert=raid1 /mnt", executor.commands[1])
}

func Test_DetectChange_Create_One(t *testing.T) {
	executor := &ExecutorStub{}
	systemd := &SystemdStub{}
	disks := &Disks{&ConfigStub{}, executor, &StatsStub{}, systemd, log.Default()}
	err := disks.Apply([]string{}, []string{"/dev/loop1"}, "uuid", true)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 1)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/mkfs.sh -U uuid -f -m single -d single /dev/loop1", executor.commands[0])
	assert.True(t, systemd.addMountCalled)
}

func Test_DetectChange_Create_One_NotFormat(t *testing.T) {
	executor := &ExecutorStub{}
	systemd := &SystemdStub{}
	disks := &Disks{&ConfigStub{}, executor, &StatsStub{}, systemd, log.Default()}
	err := disks.Apply([]string{}, []string{"/dev/loop1"}, "uuid", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 0)
	//assert.Equal(t, "/snap/platform/current/btrfs/bin/mkfs.sh -U uuid -f -m single -d single /dev/loop1", executor.commands[0])
	assert.True(t, systemd.addMountCalled)
}

func Test_DetectChange_Create_Two(t *testing.T) {
	executor := &ExecutorStub{}
	disks := &Disks{&ConfigStub{}, executor, &StatsStub{}, &SystemdStub{}, log.Default()}
	err := disks.Apply([]string{}, []string{"/dev/loop1", "/dev/loop2"}, "uuid", true)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 1)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/mkfs.sh -U uuid -f -m raid1 -d raid1 /dev/loop1 /dev/loop2", executor.commands[0])
}

func Test_DetectChange_Remove_One(t *testing.T) {
	executor := &ExecutorStub{}
	disks := &Disks{&ConfigStub{}, executor, &StatsStub{}, &SystemdStub{}, log.Default()}
	err := disks.Apply([]string{"/dev/loop1", "/dev/loop2"}, []string{"/dev/loop1"}, "", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 2)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh balance start -dconvert=single -mconvert=single /mnt", executor.commands[0])
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device delete /dev/loop2 /mnt", executor.commands[1])

}

func Test_DetectChange_Disable_One(t *testing.T) {
	executor := &ExecutorStub{}
	disks := &Disks{&ConfigStub{}, executor, &StatsStub{}, &SystemdStub{}, log.Default()}
	err := disks.Apply([]string{"/dev/loop1"}, []string{}, "", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 0)
}

func Test_DetectChange_Disable_Two(t *testing.T) {
	executor := &ExecutorStub{}
	disks := &Disks{&ConfigStub{}, executor, &StatsStub{}, &SystemdStub{}, log.Default()}
	err := disks.Apply([]string{"/dev/loop1", "/dev/loop2"}, []string{}, "", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 0)
}

func Test_DetectChange_0_To_3(t *testing.T) {
	executor := &ExecutorStub{}
	disks := &Disks{&ConfigStub{}, executor, &StatsStub{}, &SystemdStub{}, log.Default()}
	err := disks.Apply([]string{}, []string{"/dev/loop1", "/dev/loop2", "/dev/loop3"}, "uuid", true)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 1)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/mkfs.sh -U uuid -f -m raid10 -d raid10 /dev/loop1 /dev/loop2 /dev/loop3", executor.commands[0])
}

func Test_DetectChange_1_To_3(t *testing.T) {
	executor := &ExecutorStub{}
	disks := &Disks{&ConfigStub{}, executor, &StatsStub{}, &SystemdStub{}, log.Default()}
	err := disks.Apply([]string{"/dev/loop1"}, []string{"/dev/loop1", "/dev/loop2", "/dev/loop3"}, "uuid", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 2)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device add /dev/loop2 /dev/loop3 /mnt", executor.commands[0])
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh balance start -dconvert=raid10 -mconvert=raid10 /mnt", executor.commands[1])
}
