package btrfs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type StatsConfigStub struct {
}

func (c *StatsConfigStub) ExternalDiskDir() string {
	return "/mnt"
}

type StatsExecutorStub struct {
	output string
}

func (e *StatsExecutorStub) CombinedOutput(_ string, _ ...string) ([]byte, error) {
	return []byte(e.output), nil
}

func Test_HasErrors(t *testing.T) {
	executor := &StatsExecutorStub{output: `
{
  "__header": {
    "version": "1"
  },
  "device-stats": [
    {
      "device": "/dev/loop0",
      "devid": "1",
      "write_io_errs": "0",
      "read_io_errs": "0",
      "flush_io_errs": "0",
      "corruption_errs": "0",
      "generation_errs": "1"
    },
    {
      "device": "/dev/loop1",
      "devid": "2",
      "write_io_errs": "0",
      "read_io_errs": "0",
      "flush_io_errs": "0",
      "corruption_errs": "0",
      "generation_errs": "0"
    }
  ]
}
`}
	stats := NewStats(&StatsConfigStub{}, executor)

	loop0Errors, err := stats.HasErrors("/dev/loop0")
	assert.Nil(t, err)
	assert.True(t, loop0Errors)

	loop1Errors, err := stats.HasErrors("/dev/loop1")
	assert.Nil(t, err)
	assert.False(t, loop1Errors)
}
