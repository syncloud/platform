package snap

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type InstalledVersionStub struct {
	versions []string
	index    int
	err      error
}

func (s *InstalledVersionStub) InstalledVersion() (string, error) {
	if s.err != nil {
		return "", s.err
	}
	v := s.versions[s.index]
	if s.index < len(s.versions)-1 {
		s.index++
	}
	return v, nil
}

func newSnapd(version InstalledVersionProvider) *Snapd {
	s := NewSnapd(version, nil, zap.NewNop())
	s.sleep = func(time.Duration) {}
	return s
}

func TestSnapdPruneKeepsTwoNewest(t *testing.T) {
	dir := t.TempDir()
	base := time.Now()
	for idx, name := range []string{"snapd-1-amd64.tar.gz", "snapd-2-amd64.tar.gz", "snapd-3-amd64.tar.gz", "snapd-4-amd64.tar.gz"} {
		path := filepath.Join(dir, name)
		assert.NoError(t, os.WriteFile(path, []byte("x"), 0644))
		assert.NoError(t, os.Chtimes(path, base, base.Add(time.Duration(idx)*time.Hour)))
	}

	newSnapd(&InstalledVersionStub{}).prune(dir, 2)

	left, err := filepath.Glob(filepath.Join(dir, "snapd-*.tar.gz"))
	assert.NoError(t, err)
	names := []string{}
	for _, l := range left {
		names = append(names, filepath.Base(l))
	}
	assert.ElementsMatch(t, []string{"snapd-3-amd64.tar.gz", "snapd-4-amd64.tar.gz"}, names)
}

func TestSnapdVerifySucceedsWhenVersionMatches(t *testing.T) {
	err := newSnapd(&InstalledVersionStub{versions: []string{"660"}}).verify("660")
	assert.NoError(t, err)
}

func TestSnapdVerifyRetriesUntilVersionMatches(t *testing.T) {
	err := newSnapd(&InstalledVersionStub{versions: []string{"659", "659", "660"}}).verify("660")
	assert.NoError(t, err)
}

func TestSnapdVerifyFailsWhenVersionNeverMatches(t *testing.T) {
	err := newSnapd(&InstalledVersionStub{versions: []string{"659"}}).verify("660")
	assert.Error(t, err)
}

func TestSnapdVerifyFailsOnError(t *testing.T) {
	err := newSnapd(&InstalledVersionStub{err: fmt.Errorf("snapd down")}).verify("660")
	assert.Error(t, err)
}
