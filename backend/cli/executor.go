package cli

import (
	"fmt"
	"go.uber.org/zap"
	"os/exec"
)

type ShellExecutor struct {
	logger *zap.Logger
}

type Executor interface {
	CombinedOutput(name string, arg ...string) ([]byte, error)
}

func New(logger *zap.Logger) *ShellExecutor {
	return &ShellExecutor{logger: logger}
}

func (e *ShellExecutor) CombinedOutput(name string, arg ...string) ([]byte, error) {
	command := exec.Command(name, arg...)
	e.logger.Info("execute", zap.String("cmd", command.String()))
	output, err := command.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("%v: %s", err, string(output))
	}
	return output, err
}
