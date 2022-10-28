package btrfs

import (
	"fmt"
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

type SystemdStub struct {
	addMountCalled bool
}

func (s *SystemdStub) AddMount(_ string) error {
	s.addMountCalled = true
	return nil
}

func Test_Update_1_To_1(t *testing.T) {
	executor := &ExecutorStub{}
	disks := &Disks{&ConfigStub{}, executor, &SystemdStub{}, log.Default()}
	uuid, err := disks.Update([]string{"/dev/loop1", "/dev/loop2"}, []string{"/dev/loop1", "/dev/loop3"}, "uuid", false)

	assert.Nil(t, err)
	assert.Equal(t, "uuid", uuid)
	assert.Len(t, executor.commands, 3)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device add --enqueue -f /dev/loop3 /mnt", executor.commands[0])
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh balance start --enqueue -dconvert=raid1 -mconvert=raid1 /mnt", executor.commands[1])
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device delete --enqueue /dev/loop2 /mnt", executor.commands[2])
}

func Test_Update_1_To_2(t *testing.T) {
	executor := &ExecutorStub{}
	disks := &Disks{&ConfigStub{}, executor, &SystemdStub{}, log.Default()}
	_, err := disks.Update([]string{"/dev/loop1"}, []string{"/dev/loop1", "/dev/loop2"}, "", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 2)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device add --enqueue -f /dev/loop2 /mnt", executor.commands[0])
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh balance start --enqueue -dconvert=raid1 -mconvert=raid1 /mnt", executor.commands[1])
}

func Test_Update_1_To_2_Format(t *testing.T) {
	executor := &ExecutorStub{}
	systemd := &SystemdStub{}
	disks := &Disks{&ConfigStub{}, executor, systemd, log.Default()}
	uuid, err := disks.Update([]string{"/dev/loop1"}, []string{"/dev/loop1", "/dev/loop2"}, "", true)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 1)
	assert.Equal(t, fmt.Sprintf("/snap/platform/current/btrfs/bin/mkfs.sh -U %s -f -m raid1 -d raid1 /dev/loop1 /dev/loop2", uuid), executor.commands[0])
	assert.True(t, systemd.addMountCalled)
}

/*func Test_Update_2_To_2_Replace1(t *testing.T) {
	executor := &ExecutorStub{}
	disks := &Disks{&ConfigStub{}, executor, &SystemdStub{}, log.Default()}
	_, err := disks.Update([]string{"/dev/loop1", "/dev/loop2"}, []string{"/dev/loop1", "/dev/loop3"}, "", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 1)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device add --enqueue -f /dev/loop2 /mnt", executor.commands[0])
}*/

func Test_Update_0_To_1(t *testing.T) {
	executor := &ExecutorStub{}
	systemd := &SystemdStub{}
	disks := &Disks{&ConfigStub{}, executor, systemd, log.Default()}
	uuid, err := disks.Update([]string{}, []string{"/dev/loop1"}, "uuid", false)

	assert.Nil(t, err)
	assert.Equal(t, "uuid", uuid)
	assert.Len(t, executor.commands, 0)
	assert.True(t, systemd.addMountCalled)
}

func Test_Update_0_To_1_Format(t *testing.T) {
	executor := &ExecutorStub{}
	systemd := &SystemdStub{}
	disks := &Disks{&ConfigStub{}, executor, systemd, log.Default()}
	uuid, err := disks.Update([]string{}, []string{"/dev/loop1"}, "uuid", true)

	assert.Nil(t, err)
	assert.NotEqual(t, "uuid", uuid)
	assert.Len(t, executor.commands, 1)
	assert.Equal(t, fmt.Sprintf("/snap/platform/current/btrfs/bin/mkfs.sh -U %s -f -m single -d single /dev/loop1", uuid), executor.commands[0])
	assert.True(t, systemd.addMountCalled)
}

func Test_Update_0_To_2_Format(t *testing.T) {
	executor := &ExecutorStub{}
	systemd := &SystemdStub{}
	disks := &Disks{&ConfigStub{}, executor, systemd, log.Default()}
	uuid, err := disks.Update([]string{}, []string{"/dev/loop1", "/dev/loop2"}, "uuid", true)

	assert.Nil(t, err)
	assert.NotEqual(t, "uuid", uuid)
	assert.Len(t, executor.commands, 1)
	assert.Equal(t, fmt.Sprintf("/snap/platform/current/btrfs/bin/mkfs.sh -U %s -f -m raid1 -d raid1 /dev/loop1 /dev/loop2", uuid), executor.commands[0])
	assert.True(t, systemd.addMountCalled)

}

func Test_Update_2_To_1(t *testing.T) {
	executor := &ExecutorStub{}
	systemd := &SystemdStub{}
	disks := &Disks{&ConfigStub{}, executor, systemd, log.Default()}
	_, err := disks.Update([]string{"/dev/loop1", "/dev/loop2"}, []string{"/dev/loop1"}, "", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 2)
	assert.True(t, systemd.addMountCalled)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh balance start --enqueue -f -dconvert=single -mconvert=single /mnt", executor.commands[0])
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device delete --enqueue /dev/loop2 /mnt", executor.commands[1])
	assert.True(t, systemd.addMountCalled)

}

func Test_Update_1_To_0(t *testing.T) {
	executor := &ExecutorStub{}
	systemd := &SystemdStub{}
	disks := &Disks{&ConfigStub{}, executor, systemd, log.Default()}
	_, err := disks.Update([]string{"/dev/loop1"}, []string{}, "", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 0)
	assert.False(t, systemd.addMountCalled)
}

func Test_Update_2_To_0(t *testing.T) {
	executor := &ExecutorStub{}
	systemd := &SystemdStub{}
	disks := &Disks{&ConfigStub{}, executor, systemd, log.Default()}
	_, err := disks.Update([]string{"/dev/loop1", "/dev/loop2"}, []string{}, "", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 0)
	assert.False(t, systemd.addMountCalled)
}

func Test_Update_0_To_3(t *testing.T) {
	executor := &ExecutorStub{}
	systemd := &SystemdStub{}
	disks := &Disks{&ConfigStub{}, executor, systemd, log.Default()}
	uuid, err := disks.Update([]string{}, []string{"/dev/loop1", "/dev/loop2", "/dev/loop3"}, "uuid", true)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 1)
	assert.Equal(t, fmt.Sprintf("/snap/platform/current/btrfs/bin/mkfs.sh -U %s -f -m raid1 -d raid1 /dev/loop1 /dev/loop2 /dev/loop3", uuid), executor.commands[0])
	assert.True(t, systemd.addMountCalled)
}

func Test_Update_1_To_3(t *testing.T) {
	executor := &ExecutorStub{}
	systemd := &SystemdStub{}
	disks := &Disks{&ConfigStub{}, executor, systemd, log.Default()}
	_, err := disks.Update([]string{"/dev/loop1"}, []string{"/dev/loop1", "/dev/loop2", "/dev/loop3"}, "uuid", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 2)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device add --enqueue -f /dev/loop2 /dev/loop3 /mnt", executor.commands[0])
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh balance start --enqueue -dconvert=raid1 -mconvert=raid1 /mnt", executor.commands[1])
	assert.True(t, systemd.addMountCalled)
}

func Test_Update_0_To_4(t *testing.T) {
	executor := &ExecutorStub{}
	systemd := &SystemdStub{}
	disks := &Disks{&ConfigStub{}, executor, systemd, log.Default()}
	uuid, err := disks.Update([]string{}, []string{"/dev/loop1", "/dev/loop2", "/dev/loop3", "/dev/loop4"}, "uuid", true)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 1)
	assert.Equal(t, fmt.Sprintf("/snap/platform/current/btrfs/bin/mkfs.sh -U %s -f -m raid10 -d raid10 /dev/loop1 /dev/loop2 /dev/loop3 /dev/loop4", uuid), executor.commands[0])
	assert.True(t, systemd.addMountCalled)
}

func Test_Update_1_To_4(t *testing.T) {
	executor := &ExecutorStub{}
	systemd := &SystemdStub{}
	disks := &Disks{&ConfigStub{}, executor, systemd, log.Default()}
	_, err := disks.Update([]string{"/dev/loop1"}, []string{"/dev/loop1", "/dev/loop2", "/dev/loop3", "/dev/loop4"}, "uuid", false)

	assert.Nil(t, err)
	assert.Len(t, executor.commands, 2)
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh device add --enqueue -f /dev/loop2 /dev/loop3 /dev/loop4 /mnt", executor.commands[0])
	assert.Equal(t, "/snap/platform/current/btrfs/bin/btrfs.sh balance start --enqueue -dconvert=raid10 -mconvert=raid10 /mnt", executor.commands[1])
	assert.True(t, systemd.addMountCalled)
}
