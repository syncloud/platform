package stability

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeProcFile(t *testing.T, dir, rel, contents string) {
	t.Helper()
	full := filepath.Join(dir, rel)
	require.NoError(t, os.MkdirAll(filepath.Dir(full), 0755))
	require.NoError(t, os.WriteFile(full, []byte(contents), 0644))
}

func TestSnapshotParsesMemTotalAndAvailable(t *testing.T) {
	dir := t.TempDir()
	writeProcFile(t, dir, "meminfo", `MemTotal:        3789348 kB
MemFree:          203484 kB
MemAvailable:    2998888 kB
Buffers:           18472 kB
`)
	m := NewMemInfo(dir)
	s, err := m.Snapshot()
	require.NoError(t, err)
	assert.Equal(t, uint64(3789348), s.TotalKB)
	assert.Equal(t, uint64(2998888), s.AvailableKB)
	assert.InDelta(t, 0.79, s.AvailableRatio(), 0.01)
}

func TestSnapshotFailsWithoutMemTotal(t *testing.T) {
	dir := t.TempDir()
	writeProcFile(t, dir, "meminfo", "MemFree: 100 kB\n")
	_, err := NewMemInfo(dir).Snapshot()
	assert.Error(t, err)
}

func TestPSIAvailability(t *testing.T) {
	dir := t.TempDir()
	m := NewMemInfo(dir)
	assert.False(t, m.PSIAvailable())
	writeProcFile(t, dir, "pressure/memory", "")
	assert.True(t, m.PSIAvailable())
}

func TestPSIMemoryAvg10(t *testing.T) {
	dir := t.TempDir()
	writeProcFile(t, dir, "pressure/memory", `some avg10=12.34 avg60=3.21 avg300=1.00 total=12345
full avg10=4.50 avg60=1.00 avg300=0.50 total=678
`)
	v, err := NewMemInfo(dir).PSIMemoryAvg10()
	require.NoError(t, err)
	assert.InDelta(t, 12.34, v, 0.001)
}
