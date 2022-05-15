package executor

import "os/exec"

type CliExecutor struct {
}

type Executor interface {
	CommandOutput(name string, arg ...string) ([]byte, error)
}

func (e *CliExecutor) CommandOutput(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).CombinedOutput()
}
