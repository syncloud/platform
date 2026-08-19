package system

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func uptimeFile(t *testing.T, content string) *Uptime {
	file := path.Join(t.TempDir(), "uptime")
	err := os.WriteFile(file, []byte(content), 0644)
	assert.NoError(t, err)
	return &Uptime{file: file}
}

func TestUptime_Seconds(t *testing.T) {
	seconds, err := uptimeFile(t, "12345.67 98765.43\n").Seconds()
	assert.NoError(t, err)
	assert.Equal(t, 12345, seconds)
}

func TestUptime_Empty(t *testing.T) {
	_, err := uptimeFile(t, "\n").Seconds()
	assert.Error(t, err)
}

func TestUptime_NotANumber(t *testing.T) {
	_, err := uptimeFile(t, "not a number\n").Seconds()
	assert.Error(t, err)
}

func TestUptime_NoFile(t *testing.T) {
	uptime := &Uptime{file: path.Join(t.TempDir(), "missing")}
	_, err := uptime.Seconds()
	assert.Error(t, err)
}
