package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"testing"
)

func TestExecutor_CommandOutput(t *testing.T) {
	executor := New(log.Default())
	output, err := executor.CombinedOutput("date")
	assert.Nil(t, err)
	assert.Greater(t, len(string(output)), 0)
}

func TestShellExecutor_Output_IgnoresStderr(t *testing.T) {
	output, err := New(log.Default()).Output("sh", "-c", "echo out; echo warning >&2")
	assert.Nil(t, err)
	assert.Equal(t, "out\n", string(output))
}

func TestShellExecutor_Output_ErrorIncludesStderr(t *testing.T) {
	_, err := New(log.Default()).Output("sh", "-c", "echo out; echo failed >&2; exit 1")
	assert.EqualError(t, err, "exit status 1: failed\n")
}
