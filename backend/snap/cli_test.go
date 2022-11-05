package snap

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type ExecutorStub struct {
	executions []string
}

func (e *ExecutorStub) CommandOutput(name string, arg ...string) ([]byte, error) {
	e.executions = append(e.executions, fmt.Sprintf("%s %s", name, strings.Join(arg, " ")))
	return make([]byte, 0), nil
}

func TestStart(t *testing.T) {
	executor := &ExecutorStub{}
	service := NewCli(executor)
	err := service.Start("service1")
	assert.Nil(t, err)
	assert.Len(t, executor.executions, 1)
	assert.Equal(t, "snap start service1", executor.executions[0])
}

func TestStop(t *testing.T) {
	executor := &ExecutorStub{}
	service := NewCli(executor)
	err := service.Stop("service1")
	assert.Nil(t, err)
	assert.Len(t, executor.executions, 1)
	assert.Equal(t, "snap stop service1", executor.executions[0])
}
