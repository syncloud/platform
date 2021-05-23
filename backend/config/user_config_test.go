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
	config, err := NewUserConfig(db, tempFile().Name(), "syncloud.it", "https://api.syncloud.it")
	assert.Nil(t, err)

	config.UpdateRedirectDomain("syncloud.it")
	config.UpdateRedirectApiUrl("https://api.syncloud.it")
	assert.Equal(t, "syncloud.it", config.GetRedirectDomain())

	assert.Equal(t, "https://api.syncloud.it", config.GetRedirectApiUrl())

	config.UpdateRedirectDomain("syncloud.info")
	config.UpdateRedirectApiUrl("https://api.syncloud.info:81")
	assert.Equal(t, "syncloud.info", config.GetRedirectDomain())
	assert.Equal(t, "https://api.syncloud.info:81", config.GetRedirectApiUrl())
}

func TestDeviceDomain_NonActivated(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config, err := NewUserConfig(db, tempFile().Name(), "", "")
	assert.Nil(t, err)

	assert.Equal(t, "localhost", config.GetDeviceDomain())
}

func TestDeviceDomain_Free(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config, err := NewUserConfig(db, tempFile().Name(), "example.com", "")
	assert.Nil(t, err)

	config.SetRedirectEnabled(true)
	config.SetUserDomain("test")
	assert.Equal(t, "test.example.com", config.GetDeviceDomain())
}

func TestDeviceDomain_Custom(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config, err := NewUserConfig(db, tempFile().Name(), "wrong", "")
	assert.Nil(t, err)

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
	config, err := NewUserConfig(db, oldConfigFile.Name(), "syncloud.it", "https://api.syncloud.it")
	assert.Nil(t, err)

	assert.Equal(t, "syncloud.it", config.GetRedirectDomain())
	assert.True(t, config.GetUpnp())
	assert.True(t, config.IsRedirectEnabled())
	assert.False(t, config.GetExternalAccess())

	_, err = os.Stat(oldConfigFile.Name())
	assert.False(t, os.IsExist(err))
}
