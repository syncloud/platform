package stability

import (
	"errors"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type fakeKill struct {
	calls []struct {
		pid int
		sig syscall.Signal
	}
	alive   map[int]bool
	hookErr error
}

func (k *fakeKill) fn(pid int, sig syscall.Signal) error {
	k.calls = append(k.calls, struct {
		pid int
		sig syscall.Signal
	}{pid, sig})
	if sig == 0 {
		if k.alive[pid] {
			return nil
		}
		return syscall.ESRCH
	}
	if k.hookErr != nil {
		return k.hookErr
	}
	return nil
}

func (k *fakeKill) signalled(sig syscall.Signal) int {
	n := 0
	for _, c := range k.calls {
		if c.sig == sig {
			n++
		}
	}
	return n
}

func newWatcherWithProc(t *testing.T, memTotal, memAvail uint64, procDir string) *Watcher {
	t.Helper()
	return newWatcherWithPSI(t, memTotal, memAvail, procDir, "")
}

func newWatcherWithPSI(t *testing.T, memTotal, memAvail uint64, procDir, psi string) *Watcher {
	t.Helper()
	procRoot := t.TempDir()
	writeProcFile(t, procRoot, "meminfo", "MemTotal: "+strconvUint(memTotal)+" kB\nMemAvailable: "+strconvUint(memAvail)+" kB\n")
	if psi != "" {
		writeProcFile(t, procRoot, "pressure/memory", psi)
	}
	return NewWatcher(NewMemInfo(procRoot), NewProcScanner(procDir), nil, nil, zap.NewNop())
}

func TestTickNoActionWhenHealthy(t *testing.T) {
	procDir := t.TempDir()
	writeFakeProc(t, procDir, fakeProc{pid: 100, name: "photoprism", rssKB: 500000, cgroup: "0::/x"})
	w := newWatcherWithProc(t, 4000000, 3000000, procDir)
	k := &fakeKill{alive: map[int]bool{}}
	w.kill = k.fn
	require.NoError(t, w.tick())
	assert.Empty(t, k.calls)
}

func TestTickSigtermsThenSigkills(t *testing.T) {
	procDir := t.TempDir()
	writeFakeProc(t, procDir, fakeProc{pid: 100, name: "photoprism", rssKB: 500000, cgroup: "0::/x"})
	w := newWatcherWithProc(t, 4000000, 100000, procDir)
	k := &fakeKill{alive: map[int]bool{100: true}}
	w.kill = k.fn
	w.grace = 500_000_000
	require.NoError(t, w.killWorst())
	require.NotEmpty(t, k.calls)
	assert.Equal(t, syscall.SIGTERM, k.calls[0].sig)
	last := k.calls[len(k.calls)-1]
	assert.Equal(t, syscall.SIGKILL, last.sig)
}

func TestTickSkipsKillIfVictimExitedAfterTerm(t *testing.T) {
	procDir := t.TempDir()
	writeFakeProc(t, procDir, fakeProc{pid: 100, name: "photoprism", rssKB: 500000, cgroup: "0::/x"})
	w := newWatcherWithProc(t, 4000000, 100000, procDir)
	k := &fakeKill{alive: map[int]bool{}}
	w.kill = k.fn
	require.NoError(t, w.killWorst())
	for _, c := range k.calls {
		assert.NotEqual(t, syscall.SIGKILL, c.sig)
	}
}

func TestKillWorstReturnsErrNoVictim(t *testing.T) {
	procDir := t.TempDir()
	writeFakeProc(t, procDir, fakeProc{pid: 100, name: "sshd", rssKB: 5000, cgroup: "0::/system.slice/ssh.service"})
	w := newWatcherWithProc(t, 4000000, 100000, procDir)
	k := &fakeKill{alive: map[int]bool{}}
	w.kill = k.fn
	err := w.killWorst()
	assert.True(t, errors.Is(err, ErrNoVictim))
}

func TestPressureExceededOnLowAvailRegardlessOfPSI(t *testing.T) {
	w := NewWatcher(NewMemInfo(t.TempDir()), nil, nil, nil, zap.NewNop())
	assert.True(t, w.pressureExceeded(0.05, 0, false))
	assert.True(t, w.pressureExceeded(0.05, 0, true))
	assert.False(t, w.pressureExceeded(0.30, 0, false))
}

func TestPressureNotExceededByPSIAloneWhenMemoryAvailable(t *testing.T) {
	w := NewWatcher(NewMemInfo(t.TempDir()), nil, nil, nil, zap.NewNop())
	assert.False(t, w.pressureExceeded(0.423, 43.3, true))
	assert.False(t, w.pressureExceeded(0.605, 44.3, true))
	assert.False(t, w.pressureExceeded(0.731, 42.6, true))
}

func TestPressureExceededWhenPSIHighAndMemoryLow(t *testing.T) {
	w := NewWatcher(NewMemInfo(t.TempDir()), nil, nil, nil, zap.NewNop())
	assert.True(t, w.pressureExceeded(0.15, 30, true))
	assert.False(t, w.pressureExceeded(0.15, 5, true))
	assert.False(t, w.pressureExceeded(0.15, 30, false))
}

func TestTickIgnoresHighSomePSIWhenMemoryAvailable(t *testing.T) {
	procDir := t.TempDir()
	writeFakeProc(t, procDir, fakeProc{pid: 100, name: "photoprism", rssKB: 1779000, cgroup: "0::/system.slice/snap.photoprism.web.service"})
	w := newWatcherWithPSI(t, 12887000, 5448000, procDir,
		"some avg10=43.30 avg60=20.00 avg300=5.00 total=1\nfull avg10=0.12 avg60=0.05 avg300=0.01 total=2\n")
	k := &fakeKill{alive: map[int]bool{100: true}}
	w.kill = k.fn
	require.NoError(t, w.tick())
	assert.Empty(t, k.calls)
}

func TestTickKillsWhenFullPSIHighAndMemoryLow(t *testing.T) {
	procDir := t.TempDir()
	writeFakeProc(t, procDir, fakeProc{pid: 100, name: "photoprism", rssKB: 1779000, cgroup: "0::/system.slice/snap.photoprism.web.service"})
	w := newWatcherWithPSI(t, 12887000, 1288000, procDir,
		"some avg10=90.00 avg60=80.00 avg300=70.00 total=1\nfull avg10=55.00 avg60=40.00 avg300=30.00 total=2\n")
	k := &fakeKill{alive: map[int]bool{}}
	w.kill = k.fn
	require.NoError(t, w.tick())
	require.NotEmpty(t, k.calls)
	assert.Equal(t, syscall.SIGTERM, k.calls[0].sig)
}

func TestTickWithoutPSIStillKillsOnLowAvail(t *testing.T) {
	procDir := t.TempDir()
	writeFakeProc(t, procDir, fakeProc{pid: 100, name: "photoprism", rssKB: 500000, cgroup: "0::/x"})
	w := newWatcherWithProc(t, 4000000, 100000, procDir)
	assert.False(t, w.mem.PSIAvailable())
	k := &fakeKill{alive: map[int]bool{}}
	w.kill = k.fn
	require.NoError(t, w.tick())
	require.NotEmpty(t, k.calls)
	assert.Equal(t, syscall.SIGTERM, k.calls[0].sig)
}

func TestTickCooldownPreventsKillCascade(t *testing.T) {
	procDir := t.TempDir()
	writeFakeProc(t, procDir, fakeProc{pid: 100, name: "photoprism", rssKB: 1779000, cgroup: "0::/a"})
	writeFakeProc(t, procDir, fakeProc{pid: 101, name: "jellyfin", rssKB: 1635000, cgroup: "0::/b"})
	writeFakeProc(t, procDir, fakeProc{pid: 102, name: "syncthing", rssKB: 1398000, cgroup: "0::/c"})
	w := newWatcherWithProc(t, 4000000, 100000, procDir)
	k := &fakeKill{alive: map[int]bool{}}
	w.kill = k.fn
	require.NoError(t, w.tick())
	require.NoError(t, w.tick())
	require.NoError(t, w.tick())
	assert.Equal(t, 1, k.signalled(syscall.SIGTERM))
	assert.Equal(t, 100, k.calls[0].pid)
}

func TestTickKillsAgainAfterCooldownExpires(t *testing.T) {
	procDir := t.TempDir()
	writeFakeProc(t, procDir, fakeProc{pid: 100, name: "photoprism", rssKB: 1779000, cgroup: "0::/a"})
	w := newWatcherWithProc(t, 4000000, 100000, procDir)
	k := &fakeKill{alive: map[int]bool{}}
	w.kill = k.fn
	require.NoError(t, w.tick())
	w.lastAction = time.Now().Add(-2 * w.cooldown)
	require.NoError(t, w.tick())
	assert.Equal(t, 2, k.signalled(syscall.SIGTERM))
}
