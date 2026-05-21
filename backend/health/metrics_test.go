package health

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeProc(t *testing.T, dir, rel, contents string) {
	t.Helper()
	p := filepath.Join(dir, rel)
	require.NoError(t, os.MkdirAll(filepath.Dir(p), 0755))
	require.NoError(t, os.WriteFile(p, []byte(contents), 0644))
}

func newTestCollector(t *testing.T) (*Collector, string) {
	t.Helper()
	dir := t.TempDir()
	writeProc(t, dir, "stat", "cpu  1000 50 200 5000 30 0 10 0\ncpu0 ...\n")
	writeProc(t, dir, "meminfo", "MemTotal: 3700000 kB\nMemAvailable: 1500000 kB\nMemFree: 200000 kB\nBuffers: 50000 kB\nCached: 900000 kB\nSwapTotal: 2000000 kB\nSwapFree: 1500000 kB\n")
	writeProc(t, dir, "diskstats", "   8       0 sda 100 0 200 0 10 0 20 0 0 0 0 0 0 0\n"+
		"   8       1 sda1 50 0 100 0 5 0 10 0 0 0 0 0 0 0\n"+
		" 179       0 mmcblk0 1000 0 2000 0 100 0 200 0 0 0 0 0 0 0\n"+
		" 179       1 mmcblk0p1 500 0 1000 0 50 0 100 0 0 0 0 0 0 0\n"+
		"   7       0 loop0 1 0 2 0 0 0 0 0 0 0 0 0 0 0\n")
	writeProc(t, dir, "net/dev", `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
    lo: 1000      10    0    0    0     0          0         0     1000      10    0    0    0     0       0          0
  eth0: 5000      20    0    0    0     0          0         0     8000      30    0    0    0     0       0          0
`)
	return NewCollector(dir), dir
}

func TestReadCPU(t *testing.T) {
	c, _ := newTestCollector(t)
	cpu, err := c.readCPU()
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), cpu.User)
	assert.Equal(t, uint64(5000), cpu.Idle)
	assert.Equal(t, uint64(30), cpu.IOWait)
}

func TestReadMemory(t *testing.T) {
	c, _ := newTestCollector(t)
	mem, err := c.readMemory()
	require.NoError(t, err)
	assert.Equal(t, uint64(3700000), mem.TotalKB)
	assert.Equal(t, uint64(1500000), mem.AvailableKB)
	assert.Equal(t, uint64(2000000), mem.SwapTotalKB)
}

func TestReadDisksFiltersPartitionsAndLoops(t *testing.T) {
	c, _ := newTestCollector(t)
	disks, err := c.readDisks()
	require.NoError(t, err)
	names := []string{}
	for _, d := range disks {
		names = append(names, d.Name)
	}
	assert.ElementsMatch(t, []string{"sda", "mmcblk0"}, names)
}

func TestReadNetSkipsLoopback(t *testing.T) {
	c, _ := newTestCollector(t)
	nets, err := c.readNet()
	require.NoError(t, err)
	require.Len(t, nets, 1)
	assert.Equal(t, "eth0", nets[0].Name)
	assert.Equal(t, uint64(5000), nets[0].RxBytes)
	assert.Equal(t, uint64(8000), nets[0].TxBytes)
}

func TestSnapshotEndToEnd(t *testing.T) {
	c, _ := newTestCollector(t)
	s, err := c.Snapshot()
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), s.CPU.User)
	assert.Equal(t, uint64(3700000), s.Memory.TotalKB)
	require.Len(t, s.Disks, 2)
	require.Len(t, s.Net, 1)
}

func TestIsPartition(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{"sda", false},
		{"sda1", true},
		{"mmcblk0", false},
		{"mmcblk0p1", true},
		{"nvme0n1", false},
		{"nvme0n1p1", true},
		{"loop0", true},
		{"dm-0", true},
		{"ram0", true},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, isPartition(c.name), c.name)
	}
}
