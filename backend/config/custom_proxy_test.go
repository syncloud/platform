package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"path"
	"testing"
)

func newTestCustomProxy(t *testing.T) *CustomProxy {
	db := NewDb(path.Join(t.TempDir(), "db"), log.Default())
	assert.NoError(t, NewMigrator(db).Migrate())
	return NewCustomProxy(db)
}

func TestCustomProxy_AddAndList_DefaultAutheliaFalse(t *testing.T) {
	config := newTestCustomProxy(t)

	err := config.Add("legacy", "10.0.0.1", 8080, false, false)
	assert.NoError(t, err)

	entries, err := config.List()
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "legacy", entries[0].Name)
	assert.False(t, entries[0].Authelia)
}

func TestCustomProxy_AddAndList_AutheliaTrue(t *testing.T) {
	config := newTestCustomProxy(t)

	err := config.Add("guarded", "10.0.0.2", 9090, true, true)
	assert.NoError(t, err)

	entries, err := config.List()
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.True(t, entries[0].Https)
	assert.True(t, entries[0].Authelia)
}
