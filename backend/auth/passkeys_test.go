package auth

import (
	"database/sql"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func passkeysDB(t *testing.T) *Passkeys {
	dataDir := t.TempDir()
	db, err := sql.Open("sqlite", autheliaDSN(path.Join(dataDir, "authelia.sqlite3")))
	assert.NoError(t, err)
	defer db.Close()
	_, err = db.Exec(`CREATE TABLE webauthn_credentials (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_used_at DATETIME NULL DEFAULT NULL,
		username VARCHAR(100) NOT NULL,
		description VARCHAR(64) NOT NULL
	)`)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO webauthn_credentials (created_at, last_used_at, username, description) VALUES
		('2026-11-14 21:45:00.123456789+00:00', '2026-11-15 09:12:00+00:00', 'user', 'Firefox on Android'),
		('2026-11-15 10:00:00+00:00', NULL, 'user', 'Chrome on Windows'),
		('2026-11-15 11:00:00+00:00', NULL, 'other', 'Safari on iOS')`)
	assert.NoError(t, err)
	return NewPasskeys(dataDir)
}

func TestPasskeys_List(t *testing.T) {
	passkeys, err := passkeysDB(t).List("user")
	assert.NoError(t, err)
	assert.Len(t, passkeys, 2)
	assert.Equal(t, "Firefox on Android", passkeys[0].Description)
	assert.Equal(t, 2026, passkeys[0].Created.Year())
	assert.NotNil(t, passkeys[0].LastUsed)
	assert.Equal(t, 9, passkeys[0].LastUsed.Hour())
	assert.Nil(t, passkeys[1].LastUsed)
}

func TestPasskeys_List_None(t *testing.T) {
	passkeys, err := passkeysDB(t).List("nobody")
	assert.NoError(t, err)
	assert.NotNil(t, passkeys)
	assert.Empty(t, passkeys)
}

func TestPasskeys_Remove(t *testing.T) {
	p := passkeysDB(t)
	assert.NoError(t, p.Remove("user", 1))
	passkeys, _ := p.List("user")
	assert.Len(t, passkeys, 1)
	assert.Equal(t, "Chrome on Windows", passkeys[0].Description)
}

func TestPasskeys_Remove_OtherUser(t *testing.T) {
	p := passkeysDB(t)
	assert.NoError(t, p.Remove("user", 3))
	passkeys, _ := p.List("other")
	assert.Len(t, passkeys, 1)
}

func TestPasskeys_Reset(t *testing.T) {
	p := passkeysDB(t)
	assert.NoError(t, p.Reset("user"))
	passkeys, _ := p.List("user")
	assert.Empty(t, passkeys)
	others, _ := p.List("other")
	assert.Len(t, others, 1)
}

func TestPasskeys_Missing(t *testing.T) {
	_, err := NewPasskeys(t.TempDir()).List("user")
	assert.Error(t, err)
}
