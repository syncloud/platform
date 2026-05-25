package stability

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeProc struct {
	pid     int
	name    string
	rssKB   uint64
	oomAdj  int
	cgroup  string
	noVmRSS bool
}

func writeFakeProc(t *testing.T, procDir string, p fakeProc) {
	t.Helper()
	base := filepath.Join(procDir, strconv.Itoa(p.pid))
	require.NoError(t, os.MkdirAll(base, 0755))
	status := "Name:\t" + p.name + "\nState:\tS (sleeping)\n"
	if !p.noVmRSS {
		status += "VmRSS:\t" + strconv.FormatUint(p.rssKB, 10) + " kB\n"
	}
	require.NoError(t, os.WriteFile(filepath.Join(base, "status"), []byte(status), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(base, "oom_score_adj"), []byte(strconv.Itoa(p.oomAdj)+"\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(base, "cgroup"), []byte(p.cgroup+"\n"), 0644))
}

func TestCandidatesSkipsProtectedKthreadAndSelf(t *testing.T) {
	dir := t.TempDir()
	writeFakeProc(t, dir, fakeProc{pid: 2, name: "kthreadd", rssKB: 0, cgroup: "0::/"})
	writeFakeProc(t, dir, fakeProc{pid: 100, name: "kworker", noVmRSS: true, cgroup: "0::/"})
	writeFakeProc(t, dir, fakeProc{pid: 200, name: "sshd", rssKB: 5000, cgroup: "0::/system.slice/ssh.service"})
	writeFakeProc(t, dir, fakeProc{pid: 300, name: "backend", rssKB: 20000, cgroup: "0::/system.slice/snap.platform.backend.service"})
	writeFakeProc(t, dir, fakeProc{pid: 400, name: "photoprism", rssKB: 250000, cgroup: "0::/system.slice/snap.photoprism.web.service"})
	writeFakeProc(t, dir, fakeProc{pid: 500, name: "mysqld.bin", rssKB: 200000, cgroup: "0::/system.slice/snap.photoprism.mariadb.service"})

	cands, err := NewProcScanner(dir).Candidates(DefaultProtect(), 999)
	require.NoError(t, err)
	require.Len(t, cands, 2)
	assert.Equal(t, "photoprism", cands[0].Comm)
	assert.Equal(t, "mysqld.bin", cands[1].Comm)
}

func TestCandidatesScoreFavoursHighAdj(t *testing.T) {
	dir := t.TempDir()
	writeFakeProc(t, dir, fakeProc{pid: 100, name: "small_high_adj", rssKB: 50000, oomAdj: 500, cgroup: "0::/x"})
	writeFakeProc(t, dir, fakeProc{pid: 200, name: "bigger_zero_adj", rssKB: 60000, oomAdj: 0, cgroup: "0::/x"})
	cands, err := NewProcScanner(dir).Candidates(DefaultProtect(), 0)
	require.NoError(t, err)
	require.Len(t, cands, 2)
	assert.Equal(t, "small_high_adj", cands[0].Comm)
}

func TestParseAppLabel(t *testing.T) {
	assert.Equal(t, "photoprism.web", parseAppLabel("0::/system.slice/snap.photoprism.web.service"))
	assert.Equal(t, "photoprism.mariadb", parseAppLabel("0::/system.slice/snap.photoprism.mariadb.service"))
	assert.Equal(t, "platform.backend", parseAppLabel("12:devices:/system.slice/snap.platform.backend.service\n0::/system.slice/snap.platform.backend.service"))
	assert.Equal(t, "photoprism.hook.configure", parseAppLabel("0::/system.slice/snap.photoprism.hook.configure.scope"))
	assert.Equal(t, "sshd", parseAppLabel("0::/system.slice/sshd.service"))
	assert.Equal(t, "ssh", parseAppLabel("0::/system.slice/ssh.service"))
	assert.Equal(t, "init", parseAppLabel("0::/init.scope"))
	assert.Equal(t, "nginx", parseAppLabel("11:devices:/system.slice/nginx.service\n0::/system.slice/nginx.service"))
	assert.Equal(t, "", parseAppLabel("0::/"))
	assert.Equal(t, "", parseAppLabel(""))
}

func TestCandidatesPopulatesApp(t *testing.T) {
	dir := t.TempDir()
	writeFakeProc(t, dir, fakeProc{pid: 400, name: "ld.so", rssKB: 250000, cgroup: "0::/system.slice/snap.photoprism.web.service"})
	writeFakeProc(t, dir, fakeProc{pid: 500, name: "nginx", rssKB: 100000, cgroup: "0::/system.slice/nginx.service"})
	writeFakeProc(t, dir, fakeProc{pid: 600, name: "free", rssKB: 50000, cgroup: "0::/"})
	cands, err := NewProcScanner(dir).Candidates(DefaultProtect(), 999)
	require.NoError(t, err)
	require.Len(t, cands, 3)
	byPID := map[int]Victim{}
	for _, c := range cands {
		byPID[c.PID] = c
	}
	assert.Equal(t, "photoprism.web", byPID[400].App)
	assert.Equal(t, "nginx", byPID[500].App)
	assert.Equal(t, "", byPID[600].App)
}

func TestScoreFormula(t *testing.T) {
	assert.Equal(t, 100000.0, score(100000, 0))
	assert.Equal(t, 150000.0, score(100000, 500))
	assert.Equal(t, 100000.0, score(100000, -100))
}
