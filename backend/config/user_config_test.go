package config

import (
	"log"
	"os"
	"testing"
 "time"
	"github.com/stretchr/testify/assert"
)

func TestRedirectDomain(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())
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
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())
	config.Load()
	assert.Equal(t, "localhost", config.GetDeviceDomain())
}

func TestDeviceDomain_Free(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())
	config.Load()
	config.SetRedirectDomain("example.com")

	config.SetRedirectEnabled(true)
	config.SetDomain("test.example.com")
	assert.Equal(t, "test.example.com", config.GetDeviceDomain())
}

func TestDeviceBackwardsCompatibleDomain_Free(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())
	config.Load()
	config.SetRedirectDomain("example.com")

	config.SetRedirectEnabled(true)
	config.setDeprecatedUserDomain("test")
	assert.Equal(t, "test.example.com", config.GetDeviceDomain())
}

func TestDeviceDomain_Custom(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())
	config.Load()
	config.SetRedirectDomain("wrong")

	config.SetRedirectEnabled(false)
	config.SetCustomDomain("example.com")
	assert.Equal(t, "example.com", config.GetDeviceDomain())
}

func tempFile() *os.File {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal(err)
	}
	return tmpFile
}

func TestMigrate(t *testing.T) {
	oldConfigFile := tempFile()
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

	err := os.WriteFile(oldConfigFile.Name(), []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, oldConfigFile.Name())
	config.Load()
	config.SetRedirectDomain("syncloud.it")

	assert.Equal(t, "syncloud.it", config.GetRedirectDomain())

	assert.True(t, config.IsRedirectEnabled())
	assert.False(t, config.IsIpv4Public())

	_, err = os.Stat(oldConfigFile.Name())
	assert.False(t, os.IsExist(err))
}

func TestMigratev2_ExternalFalse(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())
	config.Load()
	config.Upsert("platform.external_access", "false")
	config.Load()

	assert.False(t, config.IsIpv4Public())
	assert.Nil(t, config.GetOrNilString("platform.external_access"))
}

func TestMigratev2_ExternalTrue(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())
	config.Load()
	config.Upsert("platform.external_access", "true")
	config.Load()

	assert.True(t, config.IsIpv4Public())
	assert.Nil(t, config.GetOrNilString("platform.external_access"))
}

func TestPublicIp_Empty(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())
	config.Load()
	config.SetPublicIp(nil)

	assert.Nil(t, config.GetPublicIp())
}

func TestPublicIp_Valid(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())
	config.Load()
	ip := "1.1.1.1"
	config.SetPublicIp(&ip)
	assert.Equal(t, "1.1.1.1", *config.GetPublicIp())
}

func TestBackupAppTime(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())
	config.Load()
 timesatamp := time.Now().Unix()
	config.SetBackupAppTime("app1", "backup", timesatamp)
	assert.Equal(t, timesatamp, *config.GetBackupAppTime("app1", "backup"))
}

