package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestRedirectDomain(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())

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

	assert.Equal(t, "localhost", config.GetDeviceDomain())
}

func TestDeviceDomain_Free(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())
	config.SetRedirectDomain("example.com")

	config.SetRedirectEnabled(true)
	config.SetDomain("test.example.com")
	assert.Equal(t, "test.example.com", config.GetDeviceDomain())
}

func TestDeviceBackwardsCompatibleDomain_Free(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())
	config.SetRedirectDomain("example.com")

	config.SetRedirectEnabled(true)
	config.setDeprecatedUserDomain("test")
	assert.Equal(t, "test.example.com", config.GetDeviceDomain())
}

func TestDeviceDomain_Custom(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, tempFile().Name())
	config.SetRedirectDomain("wrong")

	config.SetRedirectEnabled(false)
	config.SetCustomDomain("example.com")
	assert.Equal(t, "example.com", config.GetDeviceDomain())
}

func tempFile() *os.File {
	tmpFile, err := ioutil.TempFile("", "")
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
manual_certificate_port = 80
manual_access_port = 443
activated = True

[redirect]
domain = syncloud.it
api_url = https://api.syncloud.it
user_email = user@example.com
user_update_token = token2
`

	err := ioutil.WriteFile(oldConfigFile.Name(), []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	db := tempFile().Name()
	_ = os.Remove(db)
	config := NewUserConfig(db, oldConfigFile.Name())
	config.SetRedirectDomain("syncloud.it")

	assert.Equal(t, "syncloud.it", config.GetRedirectDomain())
	assert.True(t, config.GetUpnp())
	assert.True(t, config.IsRedirectEnabled())
	assert.False(t, config.GetExternalAccess())

	_, err = os.Stat(oldConfigFile.Name())
	assert.False(t, os.IsExist(err))
}
