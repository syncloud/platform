package systemd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/log"
	"testing"
)

type ConfigStub struct {
	diskDir string
}

func (c *ConfigStub) ConfigDir() string {
	return "/config"
}

func (c *ConfigStub) ExternalDiskDir() string {
	return c.diskDir
}

type ExecutorStub struct {
	f func(arg string) (string, error)
}

func (e *ExecutorStub) CommandOutput(_ string, args ...string) ([]byte, error) {
	f, err := e.f(args[0])
	return []byte(f), err
}

func ExecutorFunc(f func(arg string) (string, error)) cli.CommandExecutor {
	return &ExecutorStub{f}
}

func TestControl_DirToSystemdMountFilename(t *testing.T) {
	executorFunc := ExecutorFunc(
		func(arg string) (string, error) {
			return "ok", nil
		})
	control := New(executorFunc, &ConfigStub{diskDir: "/opt/disk/external"}, log.Default())
	assert.Equal(t, "dir1-dir2.mount", control.DirToSystemdMountFilename("/dir1/dir2"))
}

func TestControl_RemoveMount_Inactive(t *testing.T) {
	executorFunc := ExecutorFunc(
		func(arg string) (string, error) {
			switch arg {
			case "is-active":
				return "inactive", fmt.Errorf("error")
			case "stop":
				return "unable to stop inactive mount", fmt.Errorf("error")
			case "disable":
				return "unable to disable inactive mount", fmt.Errorf("error")
			}
			return "", fmt.Errorf("unknown command")
		})

	control := New(executorFunc, &ConfigStub{diskDir: "/opt/disk/external"}, log.Default())
	assert.Nil(t, control.RemoveMount())
}

func TestControl_RemoveMount_Active(t *testing.T) {
	executorFunc := ExecutorFunc(
		func(arg string) (string, error) {
			switch arg {
			case "is-active":
				return "active", nil
			case "stop":
				return "ok", nil
			case "disable":
				return "ok", nil
			}
			return "", fmt.Errorf("unknown command")
		})

	control := New(executorFunc, &ConfigStub{diskDir: "/opt/disk/external"}, log.Default())
	assert.Nil(t, control.RemoveMount())
}
