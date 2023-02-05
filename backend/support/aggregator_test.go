package support

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"testing"
)

func TestLogAggregator_GetLogs(t *testing.T) {
	aggregator := NewAggregator(log.Default())
	logs := aggregator.GetLogs()
	assert.NotEmpty(t, logs)
}
