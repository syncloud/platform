package stability

import (
	"errors"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type fakeKill struct {
	calls   []struct{ pid int; sig syscall.Signal }
	alive   map[int]bool
	hookErr error
}

func (k *fakeKill) fn(pid int, sig syscall.Signal) error {
	k.calls = append(k.calls, struct{ pid int; sig syscall.Signal }{pid, sig})
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

func newWatcherWithProc(t *testing.T, memTotal, memAvail uint64, procDir string) *Watcher {
	t.Helper()
	root := t.TempDir()
	procRoot := root
	writeProcFile(t, procRoot, "meminfo", "MemTotal: "+strconvUint(memTotal)+" kB\nMemAvailable: "+strconvUint(memAvail)+" kB\n")
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

func TestPressureExceededByAvailOrPSI(t *testing.T) {
	w := NewWatcher(NewMemInfo(t.TempDir()), nil, nil, nil, zap.NewNop())
	assert.True(t, w.pressureExceeded(0.05, 0, false))
	assert.False(t, w.pressureExceeded(0.30, 0, false))
	assert.True(t, w.pressureExceeded(0.30, 50, true))
	assert.False(t, w.pressureExceeded(0.30, 5, true))
}
