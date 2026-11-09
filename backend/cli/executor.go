package cli

import (
	"bytes"
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

func (e *ShellExecutor) Output(name string, arg ...string) ([]byte, error) {
	command := exec.Command(name, arg...)
	e.logger.Info("execute", zap.String("cmd", command.String()))
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		return stdout.Bytes(), fmt.Errorf("%v: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		e.logger.Warn("stderr", zap.String("cmd", command.String()), zap.String("output", stderr.String()))
	}
	return stdout.Bytes(), nil
}
