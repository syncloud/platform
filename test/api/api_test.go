package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"regexp"
	"testing"
)

func TestGetAppDir(t *testing.T) {

	dir, err := GetAppDir("platform")

	assert.Nil(t, err)
	assert.Equal(t, "/snap/platform/current", *dir)
}

func TestRestart(t *testing.T) {

	status, err := Restart("platform.nginx-public")

	assert.Nil(t, err)
	assert.Contains(t, "OK", *status)
}

func TestConfigDkimKey(t *testing.T) {
	result, err := SetDkimKey("dkim123")
	assert.Nil(t, err)
	assert.Contains(t, "OK", *result)

	key, err := GetDkimKey()
	assert.Nil(t, err)
	assert.Equal(t, "dkim123", *key)
}

func TestDataPath(t *testing.T) {
	dataDir, err := GetDataDir("platform")
	assert.Nil(t, err)
	assert.Equal(t, "/var/snap/platform/common", *dataDir)
}

func TestUrl(t *testing.T) {
	url, err := GetAppUrl("platform")
	assert.Nil(t, err)
	assert.Regexp(t, regexp.MustCompile(`\.redirect$`), *url)
	assert.Regexp(t, regexp.MustCompile(`^https://platform\.`), *url)
}

func TestGetDomainName(t *testing.T) {
	domain, err := GetDomainName("platform")
	assert.Nil(t, err)
	assert.Regexp(t, regexp.MustCompile(`\.redirect$`), *domain)
	assert.Regexp(t, regexp.MustCompile(`^platform\.`), *domain)
}

func TestGetDeviceDomainName(t *testing.T) {
	domain, err := GetDeviceDomainName()
	assert.Nil(t, err)
	assert.Regexp(t, regexp.MustCompile(`\.redirect$`), *domain)
	assert.NotRegexp(t, regexp.MustCompile(`^platform\.`), *domain)
}

func TestAppInitStorage(t *testing.T) {
	dir, err := AppInitStorage("app1", "root")
	assert.Nil(t, err)
	assert.Equal(t, "/data/app1", *dir)
}

func TestAppStorageDir(t *testing.T) {
	dir, err := AppStorageDir("app1")
	assert.Nil(t, err)
	assert.Equal(t, "/data/app1", *dir)
}

func TestUserEmail(t *testing.T) {
	email, err := UserEmail()
	assert.Nil(t, err)
	//TODO: not sure
	assert.Equal(t, "redirect", *email)
}

func TestRegisterOIDCClient(t *testing.T) {
	password, err := do(http.MethodPost, "/oidc/register", url.Values{
		"id":                         {"app1"},
		"redirect_uri":               {"/callback"},
		"require_pkce":               {"true"},
		"token_endpoint_auth_method": {"client_secret_post"},
	})
	assert.Nil(t, err)
	assert.True(t, len(*password) > 50, *password)
}
