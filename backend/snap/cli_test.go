package snap

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"strings"
	"testing"
)

type ExecutorStub struct {
	executions []string
}

func (e *ExecutorStub) CombinedOutput(name string, arg ...string) ([]byte, error) {
	e.executions = append(e.executions, fmt.Sprintf("%s %s", name, strings.Join(arg, " ")))
	return make([]byte, 0), nil
}

func TestStart(t *testing.T) {
	executor := &ExecutorStub{}
	service := NewCli(executor, log.Default())
	err := service.Start("service1")
	assert.Nil(t, err)
	assert.Len(t, executor.executions, 1)
	assert.Equal(t, "snap start service1", executor.executions[0])
}

func TestStop(t *testing.T) {
	executor := &ExecutorStub{}
	service := NewCli(executor, log.Default())
	err := service.Stop("service1")
	assert.Nil(t, err)
	assert.Len(t, executor.executions, 1)
	assert.Equal(t, "snap stop service1", executor.executions[0])
}
