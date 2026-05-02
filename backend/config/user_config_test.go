package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"path"
	"testing"
	"time"
)

func newTestUserConfig(t *testing.T) (*UserConfig, *Db) {
	db := NewDb(path.Join(t.TempDir(), "db"), log.Default())
	assert.NoError(t, NewMigrator(db).Migrate())
	return NewUserConfig(db, log.Default()), db
}

func TestDeviceDomain_NonActivated(t *testing.T) {
	config, _ := newTestUserConfig(t)
	assert.Equal(t, "www.localhost", config.GetDeviceDomain())
}

func TestDeviceDomain_Free(t *testing.T) {
	config, _ := newTestUserConfig(t)
	config.SetRedirectEnabled(true)
	config.SetDomain("test.example.com")
	assert.Equal(t, "test.example.com", config.GetDeviceDomain())
}

func TestDeviceDomain_Custom(t *testing.T) {
	config, _ := newTestUserConfig(t)
	config.SetRedirectEnabled(false)
	config.SetCustomDomain("example.com")
	assert.Equal(t, "example.com", config.GetDeviceDomain())
}

func TestPublicIp_Empty(t *testing.T) {
	config, _ := newTestUserConfig(t)
	config.SetPublicIp(nil)

	assert.Nil(t, config.GetPublicIp())
}

func TestPublicIp_Valid(t *testing.T) {
	config, _ := newTestUserConfig(t)
	ip := "1.1.1.1"
	config.SetPublicIp(&ip)
	assert.Equal(t, "1.1.1.1", *config.GetPublicIp())
}

func TestBackupAppTime(t *testing.T) {
	config, _ := newTestUserConfig(t)
	zero := config.GetBackupAppTime("app1", "backup")
	assert.True(t, zero.IsZero())
	timesatamp := time.Now()
	config.SetBackupAppTime("app1", "backup", timesatamp)
	assert.Equal(t, time.Unix(timesatamp.Unix(), 0), config.GetBackupAppTime("app1", "backup"))
}

func TestDeviceUrl(t *testing.T) {
	config, _ := newTestUserConfig(t)
	config.SetCustomDomain("domain.tld")
	port := 443
	config.SetPublicPort(&port)
	url := config.DeviceUrl()
	assert.Equal(t, "https://domain.tld", url)
}

func TestDeviceUrl_StandardPort(t *testing.T) {
	config, _ := newTestUserConfig(t)
	config.SetCustomDomain("domain.tld")
	port := 443
	config.SetPublicPort(&port)
	url := config.Url("app1")
	assert.Equal(t, "https://app1.domain.tld", url)
}

func TestDeviceUrl_NonStandardPort(t *testing.T) {
	config, _ := newTestUserConfig(t)
	config.SetCustomDomain("domain.tld")
	port := 10000
	config.SetPublicPort(&port)
	url := config.Url("app1")
	assert.Equal(t, "https://app1.domain.tld:10000", url)
}

