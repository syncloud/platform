package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAppDir(t *testing.T) {

	dir, err := GetAppDir("platform")

	assert.Nil(t, err)
	assert.Equal(t, "/snap/platform/current", dir)
}

func TestRestart(t *testing.T) {

	status, err := Restart("platform.nginx-public")

	assert.Nil(t, err)
	assert.Contains(t, "OK", status)
}

func TestConfigDkimKey(t *testing.T) {
	result, err := SetDkimKey("dkim123")
	assert.Nil(t, err)
	assert.Contains(t, "OK", result)

	key, err := GetDkimKey()
	assert.Nil(t, err)
	assert.Equal(t, "dkim123", key)
}

func TestDataPath(t *testing.T) {
	dataDir, err := GetDataDir("platform")
	assert.Nil(t, err)
	assert.Equal(t, "/var/snap/platform/current", dataDir)
}

func TestUrl(t *testing.T) {
	url, err := GetAppUrl("platform")
	assert.Nil(t, err)
	assert.Contains(t, ".syncloud.info", url)
}
