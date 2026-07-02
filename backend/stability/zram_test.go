package stability

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newTestZram(t *testing.T, memTotalKB uint64, swapsContent string) (*Zram, string) {
	t.Helper()
	root := t.TempDir()
	procDir := filepath.Join(root, "proc")
	writeProcFile(t, procDir, "meminfo", "MemTotal: "+formatKB(memTotalKB)+" kB\nMemAvailable: 1000000 kB\n")
	swapsPath := filepath.Join(root, "swaps")
	require.NoError(t, os.WriteFile(swapsPath, []byte(swapsContent), 0644))

	devFile := filepath.Join(root, "zram0")
	require.NoError(t, os.WriteFile(devFile, make([]byte, 8192), 0644))

	sysBlock := filepath.Join(root, "sys", "block", "zram0")
	require.NoError(t, os.MkdirAll(sysBlock, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(sysBlock, "comp_algorithm"), []byte("[lzo]"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(sysBlock, "disksize"), []byte("0"), 0644))

	z := &Zram{
		sysBlock:  sysBlock,
		hotAdd:    filepath.Join(root, "hot_add"),
		procSwaps: swapsPath,
		devPath:   devFile,
		mem:       NewMemInfo(procDir),
		log:       zap.NewNop(),
	}
	return z, sysBlock
}

func formatKB(kb uint64) string {
	return strconvUint(kb)
}

func strconvUint(v uint64) string {
	if v == 0 {
		return "0"
	}
	digits := []byte{}
	for v > 0 {
		digits = append([]byte{byte('0' + v%10)}, digits...)
		v /= 10
	}
	return string(digits)
}

func TestEnsureSkipsAboveMemThreshold(t *testing.T) {
	z, sysBlock := newTestZram(t, 8*1024*1024, "Filename...\n")
	require.NoError(t, z.EnsureConfigured())
	size, _ := os.ReadFile(filepath.Join(sysBlock, "disksize"))
	assert.Equal(t, "0", string(size))
}

func TestEnsureSkipsIfAlreadyOn(t *testing.T) {
	z, sysBlock := newTestZram(t, 4*1024*1024, "Filename\t\tType\tSize\tUsed\tPriority\n")
	sc, err := os.ReadFile(z.procSwaps)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(z.procSwaps, []byte(string(sc)+z.devPath+" partition 2097148 0 100\n"), 0644))
	require.NoError(t, z.EnsureConfigured())
	size, _ := os.ReadFile(filepath.Join(sysBlock, "disksize"))
	assert.Equal(t, "0", string(size))
}

func TestFileSwapsSelectsOnlyFileEntries(t *testing.T) {
	content := "Filename\t\tType\tSize\tUsed\tPriority\n" +
		"/swapfile\tfile\t2097148\t100000\t-2\n" +
		"/dev/zram0\tpartition\t1000000\t0\t10\n" +
		"/dev/sda2\tpartition\t1000\t0\t5\n"
	assert.Equal(t, []string{"/swapfile"}, fileSwaps(content))
}

func TestEnsureDeviceReturnsNilWhenPresent(t *testing.T) {
	z, _ := newTestZram(t, 4*1024*1024, "")
	require.NoError(t, z.ensureDevice())
}

func TestConfigureSysfsWritesAlgoAndSize(t *testing.T) {
	z, sysBlock := newTestZram(t, 4*1024*1024, "")
	require.NoError(t, z.configureSysfs(123456))
	algo, _ := os.ReadFile(filepath.Join(sysBlock, "comp_algorithm"))
	assert.Equal(t, "zstd", string(algo))
	size, _ := os.ReadFile(filepath.Join(sysBlock, "disksize"))
	assert.Equal(t, "123456", string(size))
}

func TestSizeBytesCapped(t *testing.T) {
	z, _ := newTestZram(t, 4*1024*1024, "")
	assert.Equal(t, uint64(2*1024*1024*1024), z.sizeBytes(8*1024*1024*1024))
	assert.Equal(t, uint64(1024*1024*1024), z.sizeBytes(2*1024*1024*1024))
}

func TestMkswapInPlaceWritesHeader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "swap")
	require.NoError(t, os.WriteFile(path, make([]byte, 1024*1024), 0644))
	require.NoError(t, mkswapInPlace(path))
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	pageSize := os.Getpagesize()
	assert.Equal(t, swapMagicV1, string(b[pageSize-10:pageSize]))
	assert.Equal(t, uint32(1), binary.LittleEndian.Uint32(b[1024:1028]))
}
