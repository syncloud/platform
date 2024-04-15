package auth

import (
	"fmt"
	"github.com/syncloud/platform/log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserConfigStub struct {
	domain *string
}

func (u *UserConfigStub) Url(app string) string {
	if u.domain != nil {
		return fmt.Sprintf("https://%s.%s", app, *u.domain)
	}
	return fmt.Sprintf("https://%s.localhost", app)
}

func (u *UserConfigStub) DeviceUrl() string {
	if u.domain != nil {
		return fmt.Sprintf("https://auth.%s", *u.domain)
	}

	return "https://localhost"
}

func (u *UserConfigStub) GetDeviceDomainNil() *string {
	return u.domain
}

type SystemdStub struct {
}

func (s *SystemdStub) RestartService(service string) error {
	return nil
}

func TestWebInit(t *testing.T) {
	config := &UserConfigStub{domain: nil}
	outDir := t.TempDir()
	secretDir := t.TempDir()
	web := NewWeb("../../config/authelia", outDir, secretDir, config, &SystemdStub{}, log.Default())
	err := web.InitConfig(false)
	assert.Nil(t, err)

	assert.FileExists(t, path.Join(secretDir, KeyFile))
	assert.FileExists(t, path.Join(secretDir, SecretFile))

	body, err := os.ReadFile(path.Join(outDir, "config.yml"))
	assert.Nil(t, err)
	assert.Contains(t, string(body), `auth.www.localhost`)
}

func TestWebReInit(t *testing.T) {
	domain := "example.com"
	config := &UserConfigStub{domain: &domain}
	outDir := t.TempDir()
	secretDir := t.TempDir()

	keyFilePath := path.Join(secretDir, KeyFile)
	err := os.WriteFile(keyFilePath, []byte("key"), 0644)
	assert.Nil(t, err)

	secretFilePath := path.Join(secretDir, SecretFile)
	err = os.WriteFile(secretFilePath, []byte("secret"), 0644)
	assert.Nil(t, err)

	web := NewWeb("../../config/authelia", outDir, secretDir, config, &SystemdStub{}, log.Default())
	err = web.InitConfig(true)
	assert.Nil(t, err)

	body, err := os.ReadFile(keyFilePath)
	assert.Nil(t, err)
	assert.Contains(t, string(body), `key`)

	body, err = os.ReadFile(secretFilePath)
	assert.Nil(t, err)
	assert.Contains(t, string(body), `secret`)

	body, err = os.ReadFile(path.Join(outDir, "config.yml"))
	assert.Nil(t, err)
	assert.Contains(t, string(body), `auth.example.com`)
}
