package storage

type ExecutorStub struct {
	output string
}

func (e *ExecutorStub) Output(_ string, _ ...string) ([]byte, error) {
	return []byte(e.output), nil
}
