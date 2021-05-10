package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestDomain(t *testing.T) {
	db := tempFile().Name()
	_ = os.Remove(db)
	config := New(db, tempFile().Name(), "syncloud.it", "https://api.syncloud.it")
	config.EnsureDb()

	config.UpdateRedirect("syncloud.it", "https://api.syncloud.it")
	assert.Equal(t, "syncloud.it", config.GetRedirectDomain())

	assert.Equal(t, "https://api.syncloud.it", config.GetRedirectApiUrl())

	config.UpdateRedirect("syncloud.info", "https://api.syncloud.info:81")
	assert.Equal(t, "syncloud.info", config.GetRedirectDomain())
	assert.Equal(t, "https://api.syncloud.info:81", config.GetRedirectApiUrl())
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
	config := New(db, oldConfigFile.Name(), "syncloud.it", "https://api.syncloud.it")
	config.EnsureDb()
	assert.Equal(t, "syncloud.it", config.GetRedirectDomain())
	assert.True(t, config.GetUpnp())
	assert.True(t, config.IsRedirectEnabled())
	assert.False(t, config.GetExternalAccess())

	_, err = os.Stat(oldConfigFile.Name())
	assert.False(t, os.IsExist(err))
}
