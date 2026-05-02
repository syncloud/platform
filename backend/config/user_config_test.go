package config

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	_ "modernc.org/sqlite"
	"os"
	"path"
	"testing"
	"time"
)

func TestRedirectDomain(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()

	config.SetRedirectDomain("syncloud.it")
	config.UpdateRedirectApiUrl("https://api.syncloud.it")
	assert.Equal(t, "syncloud.it", config.GetRedirectDomain())

	assert.Equal(t, "https://api.syncloud.it", config.GetRedirectApiUrl())

	config.SetRedirectDomain("syncloud.info")
	assert.Equal(t, "syncloud.info", config.GetRedirectDomain())
	assert.Equal(t, "https://api.syncloud.info", config.GetRedirectApiUrl())
}

func TestDeviceDomain_NonActivated(t *testing.T) {
	tempDir := t.TempDir()
	db := path.Join(tempDir, "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	assert.Equal(t, "www.localhost", config.GetDeviceDomain())
}

func TestDeviceDomain_Free(t *testing.T) {
	tempDir := t.TempDir()
	db := path.Join(tempDir, "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	config.SetRedirectDomain("example.com")

	config.SetRedirectEnabled(true)
	config.SetDomain("test.example.com")
	assert.Equal(t, "test.example.com", config.GetDeviceDomain())
}

func TestDeviceBackwardsCompatibleDomain_Free(t *testing.T) {
	tempDir := t.TempDir()
	db := path.Join(tempDir, "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	config.SetRedirectDomain("example.com")

	config.SetRedirectEnabled(true)
	config.setDeprecatedUserDomain("test")
	assert.Equal(t, "test.example.com", config.GetDeviceDomain())
}

func TestDeviceDomain_Custom(t *testing.T) {
	tempDir := t.TempDir()
	db := path.Join(tempDir, "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	config.SetRedirectDomain("wrong")

	config.SetRedirectEnabled(false)
	config.SetCustomDomain("example.com")
	assert.Equal(t, "example.com", config.GetDeviceDomain())
}

func TestMigrate(t *testing.T) {
	tempDir := t.TempDir()
	oldConfigFile := path.Join(tempDir, "old.db")

	content := `
[platform]
redirect_enabled = True
user_domain = test
domain_update_token = token1
external_access = False
manual_access_port = 443
activated = True

[redirect]
domain = syncloud.it
api_url = https://api.syncloud.it
user_email = user@example.com
user_update_token = token2
`

	err := os.WriteFile(oldConfigFile, []byte(content), 0644)
	assert.NoError(t, err)

	db := path.Join(tempDir, "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, oldConfigFile, log.Default())
	config.Load()
	config.SetRedirectDomain("syncloud.it")

	assert.Equal(t, "syncloud.it", config.GetRedirectDomain())

	assert.True(t, config.IsRedirectEnabled())
	assert.False(t, config.IsIpv4Public())

	_, err = os.Stat(oldConfigFile)
	assert.False(t, os.IsExist(err))
}

func TestMigratev2_ExternalFalse(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	config.Upsert("platform.external_access", "false")
	config.Load()

	assert.False(t, config.IsIpv4Public())
	assert.Nil(t, config.GetOrNilString("platform.external_access"))
}

func TestMigratev2_ExternalTrue(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	config.Upsert("platform.external_access", "true")
	config.Load()

	assert.True(t, config.IsIpv4Public())
	assert.Nil(t, config.GetOrNilString("platform.external_access"))
}

func TestPublicIp_Empty(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	config.SetPublicIp(nil)

	assert.Nil(t, config.GetPublicIp())
}

func TestPublicIp_Valid(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	ip := "1.1.1.1"
	config.SetPublicIp(&ip)
	assert.Equal(t, "1.1.1.1", *config.GetPublicIp())
}

func TestBackupAppTime(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	zero := config.GetBackupAppTime("app1", "backup")
	assert.True(t, zero.IsZero())
	timesatamp := time.Now()
	config.SetBackupAppTime("app1", "backup", timesatamp)
	assert.Equal(t, time.Unix(timesatamp.Unix(), 0), config.GetBackupAppTime("app1", "backup"))
}

func TestDefaultInt(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	assert.Equal(t, 0, config.GetOrDefaultInt("unknown", 0))
	config.Upsert("unknown", "1")
	assert.Equal(t, 1, config.GetOrDefaultInt("unknown", 0))
}

func TestDefaultString(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	assert.Equal(t, "default", config.GetOrDefaultString("unknown", "default"))
	config.Upsert("unknown", "test")
	assert.Equal(t, "test", config.GetOrDefaultString("unknown", "test"))
}

func TestDeviceUrl(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	config.SetCustomDomain("domain.tld")
	port := 443
	config.SetPublicPort(&port)
	url := config.DeviceUrl()
	assert.Equal(t, "https://domain.tld", url)
}

func TestDeviceUrl_StandardPort(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	config.SetCustomDomain("domain.tld")
	port := 443
	config.SetPublicPort(&port)
	//userConfig := &UserConfigMock{"domain.tld", 443}
	//device := New(userConfig)
	url := config.Url("app1")
	assert.Equal(t, "https://app1.domain.tld", url)
}

func TestDeviceUrl_NonStandardPort(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	config.SetCustomDomain("domain.tld")
	port := 10000
	config.SetPublicPort(&port)

	//userConfig := &UserConfigMock{"domain.tld", 10000}
	//device := New(userConfig)
	url := config.Url("app1")
	assert.Equal(t, "https://app1.domain.tld:10000", url)
}

func TestCustomProxy_AddAndList_DefaultAutheliaFalse(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()

	err := config.AddCustomProxy("legacy", "10.0.0.1", 8080, false, false)
	assert.NoError(t, err)

	entries, err := config.CustomProxies()
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "legacy", entries[0].Name)
	assert.False(t, entries[0].Authelia)
}

func TestCustomProxy_AddAndList_AutheliaTrue(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()

	err := config.AddCustomProxy("guarded", "10.0.0.2", 9090, true, true)
	assert.NoError(t, err)

	entries, err := config.CustomProxies()
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.True(t, entries[0].Https)
	assert.True(t, entries[0].Authelia)
}

func TestCustomProxy_MigratePreExistingProdRow(t *testing.T) {
	dbFile := path.Join(t.TempDir(), "db")
	_ = os.Remove(dbFile)

	pre, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", dbFile))
	assert.NoError(t, err)
	_, err = pre.Exec("create table config (key varchar primary key, value varchar)")
	assert.NoError(t, err)
	_, err = pre.Exec("create table custom_proxy (name varchar primary key, host varchar, port integer, https integer not null default 0)")
	assert.NoError(t, err)
	_, err = pre.Exec("INSERT INTO custom_proxy(name, host, port, https) VALUES ('prod-entry', '192.168.1.5', 8080, 1)")
	assert.NoError(t, err)
	assert.NoError(t, pre.Close())

	config := NewUserConfig(dbFile, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()

	entries, err := config.CustomProxies()
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "prod-entry", entries[0].Name)
	assert.Equal(t, "192.168.1.5", entries[0].Host)
	assert.Equal(t, 8080, entries[0].Port)
	assert.True(t, entries[0].Https)
	assert.False(t, entries[0].Authelia, "existing prod rows must default to authelia=false after upgrade")
}

func TestUserConfig_OIDCClients(t *testing.T) {
	db := path.Join(t.TempDir(), "db")
	_ = os.Remove(db)
	config := NewUserConfig(db, path.Join(t.TempDir(), "old.db"), log.Default())
	config.Load()
	config.SetCustomDomain("example.com")
	err := config.AddOIDCClient(OIDCClient{
		ID:                      "app1",
		Secret:                  "secret",
		RedirectURI:             "/callback",
		RequirePkce:             true,
		TokenEndpointAuthMethod: "client_secret_post",
	})
	assert.NoError(t, err)

	clients, err := config.OIDCClients()
	assert.NoError(t, err)
	assert.Equal(t, "app1", clients[0].ID)
	assert.Equal(t, "secret", clients[0].Secret)
	assert.Equal(t, "https://app1.example.com/callback", clients[0].RedirectURI)
	assert.True(t, clients[0].RequirePkce)
	assert.Equal(t, "client_secret_post", clients[0].TokenEndpointAuthMethod)
}
