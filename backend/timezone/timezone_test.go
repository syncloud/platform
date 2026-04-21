package timezone

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeExecutor struct {
	calls [][]string
	out   []byte
	err   error
}

func (f *fakeExecutor) CombinedOutput(name string, arg ...string) ([]byte, error) {
	f.calls = append(f.calls, append([]string{name}, arg...))
	return f.out, f.err
}

type fakeStore struct {
	tz string
}

func (f *fakeStore) SetTimezone(tz string) { f.tz = tz }
func (f *fakeStore) GetTimezone() string   { return f.tz }

func pickValidTZ(t *testing.T) string {
	t.Helper()
	for _, name := range []string{"UTC", "Europe/London", "America/New_York"} {
		if info, err := os.Stat(filepath.Join("/usr/share/zoneinfo", name)); err == nil && !info.IsDir() {
			return name
		}
	}
	t.Skip("no zoneinfo database available on this host")
	return ""
}

func TestApplier_Apply_Valid(t *testing.T) {
	tz := pickValidTZ(t)
	exec := &fakeExecutor{}
	store := &fakeStore{}
	applier := NewApplier(exec, store)

	err := applier.Apply(tz)

	assert.NoError(t, err)
	assert.Equal(t, tz, store.tz)
	assert.Equal(t, [][]string{{"timedatectl", "set-timezone", tz}}, exec.calls)
}

func TestApplier_Apply_Invalid(t *testing.T) {
	exec := &fakeExecutor{}
	store := &fakeStore{}
	applier := NewApplier(exec, store)

	err := applier.Apply("Not/A/Real/Zone")

	assert.Error(t, err)
	assert.Empty(t, store.tz)
	assert.Empty(t, exec.calls, "timedatectl should not be invoked for invalid tz")
}

func TestApplier_Apply_PathTraversal(t *testing.T) {
	exec := &fakeExecutor{}
	store := &fakeStore{}
	applier := NewApplier(exec, store)

	for _, evil := range []string{"../etc/passwd", "/etc/passwd", ""} {
		err := applier.Apply(evil)
		assert.Error(t, err, "expected error for %q", evil)
	}
	assert.Empty(t, exec.calls)
}

func TestApplier_Current(t *testing.T) {
	store := &fakeStore{tz: "Asia/Tokyo"}
	applier := NewApplier(&fakeExecutor{}, store)

	assert.Equal(t, "Asia/Tokyo", applier.Current())
}
