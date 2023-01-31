package cli

import (
	"fmt"
	"go.uber.org/zap"
	"os/exec"
	"strings"
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
	e.logger.Info("execute", zap.String("cmd", fmt.Sprintf("%s %s", name, strings.Join(arg, " "))))
	output, err := exec.Command(name, arg...).CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("%v: %s", err, string(output))
	}
	return output, err
}
