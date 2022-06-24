package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"testing"
)

func TestExecutor_CommandOutput(t *testing.T) {
	executor := NewExecutor(log.Default())
	output, err := executor.CommandOutput("date")
	assert.Nil(t, err)
	assert.Greater(t, len(string(output)), 0)
}
