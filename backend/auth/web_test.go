package auth

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/log"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"testing"
)

type UserConfigStub struct {
	domain    *string
	activated bool
	clients   []config.OIDCClient
}

func (u *UserConfigStub) AddOIDCClient(client config.OIDCClient) error {
	return nil
}

func (u *UserConfigStub) IsActivated() bool {
	return u.activated
}

func (u *UserConfigStub) Url(app string) string {
	if u.domain != nil {
		return fmt.Sprintf("https://%s.%s", app, *u.domain)
	}
	return fmt.Sprintf("https://%s.localhost", app)
}

func (u *UserConfigStub) OIDCClients() ([]config.OIDCClient, error) {
	return u.clients, nil
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

func (s *SystemdStub) RestartService(_ string) error {
	return nil
}

type PasswordGeneratorStub struct {
}

func (p *PasswordGeneratorStub) Generate() (Secret, error) {
	return Secret{Password: "pass", Hash: "hash"}, nil
}

func TestWebInit(t *testing.T) {
	userConfig := &UserConfigStub{domain: nil, activated: false}
	outDir := t.TempDir()
	secretDir := t.TempDir()
	web := NewWeb("../../config/authelia", outDir, secretDir, userConfig, &SystemdStub{}, &PasswordGeneratorStub{}, log.Default())
	err := web.InitConfig()
	assert.NoError(t, err)

	assert.FileExists(t, path.Join(secretDir, KeyFile))
	assert.FileExists(t, path.Join(secretDir, SecretFile))

	body, err := os.ReadFile(path.Join(outDir, "config.yml"))
	assert.Nil(t, err)
	assert.Contains(t, string(body), `auth.www.localhost`)
}

func TestWebReInit(t *testing.T) {
	domain := "example.com"
	userConfig := &UserConfigStub{domain: &domain, activated: true}
	outDir := t.TempDir()
	secretDir := t.TempDir()

	keyFilePath := path.Join(secretDir, KeyFile)
	err := os.WriteFile(keyFilePath, []byte("key"), 0644)
	assert.Nil(t, err)

	secretFilePath := path.Join(secretDir, SecretFile)
	err = os.WriteFile(secretFilePath, []byte("secret"), 0644)
	assert.Nil(t, err)

	web := NewWeb("../../config/authelia", outDir, secretDir, userConfig, &SystemdStub{}, &PasswordGeneratorStub{}, log.Default())
	err = web.InitConfig()
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

type Config struct {
	IdentityProviders IdentityProvider `yaml:"identity_providers"`
}

type IdentityProvider struct {
	OIDC OIDC `yaml:"oidc"`
}

type OIDC struct {
	Clients []Client `yaml:"clients"`
}

type Client struct {
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	RedirectUris []string `yaml:"redirect_uris"`
}

func TestWebClients(t *testing.T) {
	domain := "example.com"
	userConfig := &UserConfigStub{domain: &domain, clients: []config.OIDCClient{
		{ID: "app1", Secret: "app1secret", RedirectURI: "https://app1.example.com/callback1"},
		{ID: "app2", Secret: "app2secret", RedirectURI: "https://app2.example.com/callback2"},
	}, activated: false}
	outDir := t.TempDir()
	secretDir := t.TempDir()
	web := NewWeb("../../config/authelia", outDir, secretDir, userConfig, &SystemdStub{}, &PasswordGeneratorStub{}, log.Default())
	err := web.InitConfig()
	assert.NoError(t, err)

	body, err := os.ReadFile(path.Join(outDir, "config.yml"))
	assert.NoError(t, err)

	gen := Config{}
	err = yaml.Unmarshal(body, &gen)
	assert.NoError(t, err)

	assert.Len(t, gen.IdentityProviders.OIDC.Clients, 3)
	assert.Equal(t, "syncloud", gen.IdentityProviders.OIDC.Clients[0].ClientID)

	assert.Equal(t, "app1", gen.IdentityProviders.OIDC.Clients[1].ClientID)
	assert.Equal(t, "app1secret", gen.IdentityProviders.OIDC.Clients[1].ClientSecret)
	assert.Len(t, gen.IdentityProviders.OIDC.Clients[1].RedirectUris, 1)
	assert.Equal(t, "https://app1.example.com/callback1", gen.IdentityProviders.OIDC.Clients[1].RedirectUris[0])

	assert.Equal(t, "app2", gen.IdentityProviders.OIDC.Clients[2].ClientID)
	assert.Equal(t, "app2secret", gen.IdentityProviders.OIDC.Clients[2].ClientSecret)
	assert.Equal(t, "https://app2.example.com/callback2", gen.IdentityProviders.OIDC.Clients[2].RedirectUris[0])
}
