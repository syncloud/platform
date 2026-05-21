package stability

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestMigrateMovesWhenOnlyOldExists(t *testing.T) {
	dir := t.TempDir()
	old := filepath.Join(dir, "common", "events.jsonl")
	newP := filepath.Join(dir, "data", "events.jsonl")
	require.NoError(t, os.MkdirAll(filepath.Dir(old), 0755))
	require.NoError(t, os.MkdirAll(filepath.Dir(newP), 0755))
	require.NoError(t, os.WriteFile(old, []byte("{\"kind\":\"x\"}\n"), 0644))

	MigrateEventLog(old, newP, zap.NewNop())

	_, errOld := os.Stat(old)
	assert.True(t, os.IsNotExist(errOld))
	body, err := os.ReadFile(newP)
	require.NoError(t, err)
	assert.Equal(t, "{\"kind\":\"x\"}\n", string(body))
}

func TestMigrateSkipsWhenNewAlreadyExists(t *testing.T) {
	dir := t.TempDir()
	old := filepath.Join(dir, "events.jsonl")
	newP := filepath.Join(dir, "new.jsonl")
	require.NoError(t, os.WriteFile(old, []byte("old"), 0644))
	require.NoError(t, os.WriteFile(newP, []byte("new"), 0644))

	MigrateEventLog(old, newP, zap.NewNop())

	oldBody, _ := os.ReadFile(old)
	newBody, _ := os.ReadFile(newP)
	assert.Equal(t, "old", string(oldBody), "old file untouched if new exists")
	assert.Equal(t, "new", string(newBody), "new file untouched if new exists")
}

func TestMigrateNoopWhenOldMissing(t *testing.T) {
	dir := t.TempDir()
	old := filepath.Join(dir, "missing.jsonl")
	newP := filepath.Join(dir, "new.jsonl")

	MigrateEventLog(old, newP, zap.NewNop())

	_, err := os.Stat(newP)
	assert.True(t, os.IsNotExist(err))
}
