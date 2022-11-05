package du

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type ExecutorStub struct {
	output string
}

func (e *ExecutorStub) CommandOutput(_ string, _ ...string) ([]byte, error) {
	return []byte(e.output), nil
}

func TestShellDiskUsage_Used(t *testing.T) {
	output := "125    ."

	usage := New(&ExecutorStub{output: output})

	bytes, err := usage.Used("")
	assert.Nil(t, err)
	assert.Equal(t, uint64(125), bytes)

}
