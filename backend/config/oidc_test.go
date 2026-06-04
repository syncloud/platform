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

func TestOIDC_AddAndList_RoundTripsRawRedirectURI(t *testing.T) {
	o := newTestOIDC(t)

	err := o.AddClient(OIDCClient{
		ID:                      "app1",
		Secret:                  "secret",
		RedirectURI:             "/callback",
		RequirePkce:             true,
		TokenEndpointAuthMethod: "client_secret_post",
	})
	assert.NoError(t, err)

	clients, err := o.Clients()
	assert.NoError(t, err)
	assert.Len(t, clients, 1)
	assert.Equal(t, "app1", clients[0].ID)
	assert.Equal(t, "secret", clients[0].Secret)
	assert.Equal(t, "/callback", clients[0].RedirectURI)
	assert.True(t, clients[0].RequirePkce)
	assert.Equal(t, "client_secret_post", clients[0].TokenEndpointAuthMethod)
}

func TestOIDC_AddClient_OverwritesExistingByID(t *testing.T) {
	o := newTestOIDC(t)

	assert.NoError(t, o.AddClient(OIDCClient{ID: "app1", Secret: "first", RedirectURI: "/a"}))
	assert.NoError(t, o.AddClient(OIDCClient{ID: "app1", Secret: "second", RedirectURI: "/b"}))

	clients, err := o.Clients()
	assert.NoError(t, err)
	assert.Len(t, clients, 1)
	assert.Equal(t, "second", clients[0].Secret)
	assert.Equal(t, "/b", clients[0].RedirectURI)
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
	assert.Equal(t, "/auth/openid/callback", clients[0].RedirectURI)
}

func TestOIDC_RedirectURIContainingCommaIsNotSplit(t *testing.T) {
	o := newTestOIDC(t)

	err := o.AddClient(OIDCClient{
		ID:           "app1",
		Secret:       "secret",
		RedirectURIs: []string{"/cb?ids=1,2,3", "/mobile"},
	})
	assert.NoError(t, err)

	clients, err := o.Clients()
	assert.NoError(t, err)
	assert.Len(t, clients, 1)
	assert.Equal(t, []string{"/cb?ids=1,2,3", "/mobile"}, clients[0].RedirectURIs)
}

func TestOIDC_LegacySingleRedirectRowSurvivesMigration(t *testing.T) {
	db := NewDb(path.Join(t.TempDir(), "db"), log.Default())
	assert.NoError(t, db.Init())
	_, err := db.Exec(`create table oidc_client
		(id varchar primary key, secret varchar, redirect_uri varchar, require_pkce integer, token_endpoint_auth_method varchar)`)
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO oidc_client VALUES ('legacy', 's', '/old/callback', 1, 'client_secret_basic')")
	assert.NoError(t, err)

	assert.NoError(t, NewMigrator(db).Migrate())

	clients, err := NewOIDC(db).Clients()
	assert.NoError(t, err)
	assert.Len(t, clients, 1)
	assert.Equal(t, "/old/callback", clients[0].RedirectURI)
	assert.Equal(t, []string{"/old/callback"}, clients[0].RedirectURIs)
}
