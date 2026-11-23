package ca

type ExecutorStub struct {
	commands []string
}

func (e *ExecutorStub) CombinedOutput(name string, arg ...string) ([]byte, error) {
	e.commands = append(e.commands, name+" "+arg[0])
	return nil, nil
}
