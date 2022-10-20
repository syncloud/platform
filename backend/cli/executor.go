package cli

import (
	"go.uber.org/zap"
	"os/exec"
)

type Executor struct {
	logger *zap.Logger
}

type CommandExecutor interface {
	CommandOutput(name string, arg ...string) ([]byte, error)
}

func NewExecutor(logger *zap.Logger) *Executor {
	return &Executor{logger: logger}
}

func (e *Executor) CommandOutput(name string, arg ...string) ([]byte, error) {
	e.logger.Info("execute", zap.Strings(name, arg))
	return exec.Command(name, arg...).CombinedOutput()
}
