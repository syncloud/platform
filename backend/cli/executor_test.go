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
