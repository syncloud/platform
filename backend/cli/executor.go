package cli

import "os/exec"

type Executor struct {
}

type CommandExecutor interface {
	CommandOutput(name string, arg ...string) ([]byte, error)
}

func (e *Executor) CommandOutput(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).CombinedOutput()
}
