package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"path"
	"testing"
)

func newTestRedirect(t *testing.T) *Redirect {
	db := NewDb(path.Join(t.TempDir(), "db"), log.Default())
	assert.NoError(t, NewMigrator(db).Migrate())
	return NewRedirect(db)
}

func TestRedirectDomain(t *testing.T) {
	redirect := newTestRedirect(t)

	redirect.SetDomain("syncloud.it")
	redirect.UpdateApiUrl("https://api.syncloud.it")
	assert.Equal(t, "syncloud.it", redirect.Domain())
	assert.Equal(t, "https://api.syncloud.it", redirect.ApiUrl())

	redirect.SetDomain("syncloud.info")
	assert.Equal(t, "syncloud.info", redirect.Domain())
	assert.Equal(t, "https://api.syncloud.info", redirect.ApiUrl())
}
