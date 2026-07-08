package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type SnapdStub struct {
	versions []string
	index    int
	err      error
}

func (s *SnapdStub) InstalledVersion() (string, error) {
	if s.err != nil {
		return "", s.err
	}
	v := s.versions[s.index]
	if s.index < len(s.versions)-1 {
		s.index++
	}
	return v, nil
}

func newInstaller(snapd SnapdVersion) *Installer {
	i := New(snapd, zap.NewNop())
	i.sleep = func(time.Duration) {}
	return i
}

func TestPruneKeepsTwoNewest(t *testing.T) {
	dir := t.TempDir()
	base := time.Now()
	for idx, name := range []string{"snapd-1-amd64.tar.gz", "snapd-2-amd64.tar.gz", "snapd-3-amd64.tar.gz", "snapd-4-amd64.tar.gz"} {
		path := filepath.Join(dir, name)
		assert.NoError(t, os.WriteFile(path, []byte("x"), 0644))
		assert.NoError(t, os.Chtimes(path, base, base.Add(time.Duration(idx)*time.Hour)))
	}

	newInstaller(&SnapdStub{}).prune(dir, 2)

	left, err := filepath.Glob(filepath.Join(dir, "snapd-*.tar.gz"))
	assert.NoError(t, err)
	names := []string{}
	for _, l := range left {
		names = append(names, filepath.Base(l))
	}
	assert.ElementsMatch(t, []string{"snapd-3-amd64.tar.gz", "snapd-4-amd64.tar.gz"}, names)
}

func TestVerifySucceedsWhenVersionMatches(t *testing.T) {
	err := newInstaller(&SnapdStub{versions: []string{"660"}}).verify("660")
	assert.NoError(t, err)
}

func TestVerifyRetriesUntilVersionMatches(t *testing.T) {
	err := newInstaller(&SnapdStub{versions: []string{"659", "659", "660"}}).verify("660")
	assert.NoError(t, err)
}

func TestVerifyFailsWhenVersionNeverMatches(t *testing.T) {
	err := newInstaller(&SnapdStub{versions: []string{"659"}}).verify("660")
	assert.Error(t, err)
}

func TestVerifyFailsOnError(t *testing.T) {
	err := newInstaller(&SnapdStub{err: fmt.Errorf("snapd down")}).verify("660")
	assert.Error(t, err)
}
