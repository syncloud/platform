package systemd

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"testing"
)

type ConfigStub struct {
	diskDir string
}

func (c *ConfigStub) ExternalDiskDir() string {
	return c.diskDir
}

type ExecutorStub struct {
	output string
}

func (e *ExecutorStub) CommandOutput(_ string, _ ...string) ([]byte, error) {
	return []byte(e.output), nil
}

func TestLsblk_AvailableDisks_FindRootPartitionSome(t *testing.T) {
	control := New(&ExecutorStub{output: ""}, &ConfigStub{diskDir: "/opt/disk/external"}, log.Default())
	assert.Equal(t, "dir1-dir2.mount", control.DirToSystemdMountFilename("/dir1/dir2"))
}
