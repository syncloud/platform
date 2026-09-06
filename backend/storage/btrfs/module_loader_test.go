package btrfs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type LoaderExecutorStub struct {
	commands []string
	err      error
	output   string
}

func (e *LoaderExecutorStub) CombinedOutput(command string, args ...string) ([]byte, error) {
	e.commands = append(e.commands, command+" "+args[0])
	return []byte(e.output), e.err
}

func TestModuleLoader_Load(t *testing.T) {
	executor := &LoaderExecutorStub{}
	err := NewModuleLoader(executor).Load("btrfs")

	assert.Nil(t, err)
	assert.Equal(t, []string{"modprobe btrfs"}, executor.commands)
}

func TestModuleLoader_Load_Error(t *testing.T) {
	executor := &LoaderExecutorStub{err: fmt.Errorf("exit status 1"), output: "modprobe: FATAL: Module btrfs not found"}
	err := NewModuleLoader(executor).Load("btrfs")

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Module btrfs not found")
}
