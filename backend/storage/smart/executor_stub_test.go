package smart

import (
	"fmt"
	"os"
	"strings"
)

type ExecutorStub struct {
	responses map[string]string
	calls     []string
}

func (e *ExecutorStub) CombinedOutput(_ string, args ...string) ([]byte, error) {
	call := strings.Join(args, " ")
	e.calls = append(e.calls, call)
	fixture, ok := e.responses[call]
	if !ok {
		return []byte("smartctl: command not found"), fmt.Errorf("exit status 127")
	}
	content, err := os.ReadFile("testdata/" + fixture)
	if err != nil {
		return nil, err
	}
	return content, nil
}
