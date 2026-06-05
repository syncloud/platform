package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"path"
	"testing"
)

func newTestOIDC(t *testing.T) *OIDC {
	db := NewDb(path.Join(t.TempDir(), "db"), log.Default())
	assert.NoError(t, NewMigrator(db).Migrate())
	return NewOIDC(db)
}

func TestOIDC_AddAndList_RoundTripsRedirectURI(t *testing.T) {
	o := newTestOIDC(t)

	err := o.AddClient(OIDCClient{
		ID:                      "app1",
		Secret:                  "secret",
		RedirectURIs:            []string{"/callback"},
		RequirePkce:             true,
		TokenEndpointAuthMethod: "client_secret_post",
	})
	assert.NoError(t, err)

	clients, err := o.Clients()
	assert.NoError(t, err)
	assert.Len(t, clients, 1)
	assert.Equal(t, "app1", clients[0].ID)
	assert.Equal(t, "secret", clients[0].Secret)
	assert.Equal(t, []string{"/callback"}, clients[0].RedirectURIs)
	assert.True(t, clients[0].RequirePkce)
	assert.Equal(t, "client_secret_post", clients[0].TokenEndpointAuthMethod)
}

func TestOIDC_AddClient_OverwritesExistingByID(t *testing.T) {
	o := newTestOIDC(t)

	assert.NoError(t, o.AddClient(OIDCClient{ID: "app1", Secret: "first", RedirectURIs: []string{"/a"}}))
	assert.NoError(t, o.AddClient(OIDCClient{ID: "app1", Secret: "second", RedirectURIs: []string{"/b", "/c"}}))

	clients, err := o.Clients()
	assert.NoError(t, err)
	assert.Len(t, clients, 1)
	assert.Equal(t, "second", clients[0].Secret)
	assert.Equal(t, []string{"/b", "/c"}, clients[0].RedirectURIs)
}

func TestOIDC_Clients_EmptyByDefault(t *testing.T) {
	o := newTestOIDC(t)
	clients, err := o.Clients()
	assert.NoError(t, err)
	assert.Empty(t, clients)
}

func TestOIDC_AddAndList_MultipleRedirectURIs(t *testing.T) {
	o := newTestOIDC(t)

	err := o.AddClient(OIDCClient{
		ID:           "app1",
		Secret:       "secret",
		RedirectURIs: []string{"/auth/openid/callback", "/auth/openid/mobile-redirect"},
		RequirePkce:  true,
	})
	assert.NoError(t, err)

	clients, err := o.Clients()
	assert.NoError(t, err)
	assert.Len(t, clients, 1)
	assert.Equal(t, []string{"/auth/openid/callback", "/auth/openid/mobile-redirect"}, clients[0].RedirectURIs)
}
