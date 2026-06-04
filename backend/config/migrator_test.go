package config

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	_ "modernc.org/sqlite"
	"path"
	"testing"
)

func TestMigrator_CreatesSchemaFromScratch(t *testing.T) {
	db := NewDb(path.Join(t.TempDir(), "db"), log.Default())
	m := NewMigrator(db)

	assert.NoError(t, m.Migrate())

	_, err := db.Exec("INSERT INTO custom_proxy(name, host, port, https, authelia) VALUES ('p', 'h', 1, 0, 1)")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO oidc_client(id, secret, require_pkce, token_endpoint_auth_method) VALUES ('id', 's', 0, 'm')")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO oidc_redirect_uri(client_id, redirect_uri) VALUES ('id', '/cb')")
	assert.NoError(t, err)
}

func TestMigrator_MigratesLegacyRedirectUriIntoTableAndDropsColumn(t *testing.T) {
	db := NewDb(path.Join(t.TempDir(), "db"), log.Default())
	m := NewMigrator(db)

	assert.NoError(t, m.MigrateTo(5))
	_, err := db.Exec("INSERT INTO oidc_client(id, secret, redirect_uri, require_pkce, token_endpoint_auth_method) VALUES ('app', 's', '/old/callback', 1, 'client_secret_basic')")
	assert.NoError(t, err)

	assert.NoError(t, m.Migrate())

	conn := db.Open()
	defer conn.Close()

	rows, err := conn.Query("select client_id, redirect_uri from oidc_redirect_uri")
	assert.NoError(t, err)
	defer rows.Close()
	assert.True(t, rows.Next())
	var clientID, redirectURI string
	assert.NoError(t, rows.Scan(&clientID, &redirectURI))
	assert.Equal(t, "app", clientID)
	assert.Equal(t, "/old/callback", redirectURI)
	assert.False(t, rows.Next())

	_, err = conn.Query("select redirect_uri from oidc_client")
	assert.Error(t, err, "legacy redirect_uri column must be dropped from oidc_client")
}

func TestMigrator_IsIdempotent(t *testing.T) {
	db := NewDb(path.Join(t.TempDir(), "db"), log.Default())
	m := NewMigrator(db)

	assert.NoError(t, m.Migrate())
	assert.NoError(t, m.Migrate())
	assert.NoError(t, m.Migrate())
}

func TestMigrator_AddsHttpsAndAutheliaColumnsToLegacyTable(t *testing.T) {
	db := NewDb(path.Join(t.TempDir(), "db"), log.Default())
	assert.NoError(t, db.Init())
	_, err := db.Exec("create table custom_proxy (name varchar primary key, host varchar, port integer)")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO custom_proxy VALUES ('legacy', 'h', 1)")
	assert.NoError(t, err)

	assert.NoError(t, NewMigrator(db).Migrate())

	_, err = db.Exec("UPDATE custom_proxy SET https = 1, authelia = 1 WHERE name = 'legacy'")
	assert.NoError(t, err)
}

func TestMigrator_PreservesPreExistingProdRowWithAutheliaDefault(t *testing.T) {
	dbFile := path.Join(t.TempDir(), "db")

	pre, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", dbFile))
	assert.NoError(t, err)
	_, err = pre.Exec("create table config (key varchar primary key, value varchar)")
	assert.NoError(t, err)
	_, err = pre.Exec("create table custom_proxy (name varchar primary key, host varchar, port integer, https integer not null default 0)")
	assert.NoError(t, err)
	_, err = pre.Exec("INSERT INTO custom_proxy(name, host, port, https) VALUES ('prod-entry', '192.168.1.5', 8080, 1)")
	assert.NoError(t, err)
	assert.NoError(t, pre.Close())

	db := NewDb(dbFile, log.Default())
	assert.NoError(t, NewMigrator(db).Migrate())

	rows, err := db.Open().Query("select name, host, port, https, authelia from custom_proxy")
	assert.NoError(t, err)
	defer rows.Close()
	assert.True(t, rows.Next())
	var name, host string
	var port, https, authelia int
	assert.NoError(t, rows.Scan(&name, &host, &port, &https, &authelia))
	assert.Equal(t, "prod-entry", name)
	assert.Equal(t, "192.168.1.5", host)
	assert.Equal(t, 8080, port)
	assert.Equal(t, 1, https)
	assert.Equal(t, 0, authelia, "existing prod rows must default to authelia=0 after upgrade")
}
