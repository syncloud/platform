package system

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type ExecutorStub struct {
	commands chan []string
	err      error
}

func (e *ExecutorStub) CombinedOutput(name string, arg ...string) ([]byte, error) {
	e.commands <- append([]string{name}, arg...)
	return []byte("output"), e.err
}

func power(err error) (*Power, *ExecutorStub) {
	executor := &ExecutorStub{commands: make(chan []string, 1), err: err}
	return &Power{executor: executor, delay: time.Millisecond, logger: zap.NewNop()}, executor
}

func waitCommand(t *testing.T, executor *ExecutorStub) []string {
	select {
	case command := <-executor.commands:
		return command
	case <-time.After(5 * time.Second):
		t.Fatal("command was not executed")
		return nil
	}
}

func TestPower_Restart(t *testing.T) {
	p, executor := power(nil)
	p.Restart()
	assert.Equal(t, []string{"shutdown", "-r", "now"}, waitCommand(t, executor))
}

func TestPower_Shutdown(t *testing.T) {
	p, executor := power(nil)
	p.Shutdown()
	assert.Equal(t, []string{"shutdown", "now"}, waitCommand(t, executor))
}

func TestPower_DoesNotBlockOnFailure(t *testing.T) {
	p, executor := power(fmt.Errorf("no shutdown"))
	p.Restart()
	assert.Equal(t, []string{"shutdown", "-r", "now"}, waitCommand(t, executor))
}
