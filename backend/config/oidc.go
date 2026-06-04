package config

import "strings"

type OIDCClient struct {
	ID                      string
	Secret                  string
	RedirectURI             string
	RedirectURIs            []string
	RequirePkce             bool
	TokenEndpointAuthMethod string
}

type OIDC struct {
	db *Db
}

func NewOIDC(db *Db) *OIDC {
	return &OIDC{db: db}
}

func (o *OIDC) Clients() ([]OIDCClient, error) {
	db := o.db.Open()
	defer db.Close()
	rows, err := db.Query("select id, secret, redirect_uri, require_pkce, token_endpoint_auth_method, redirect_uris from oidc_client")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := make([]OIDCClient, 0)
	for rows.Next() {
		var client OIDCClient
		var requirePkce int
		var redirectURIs string
		if err := rows.Scan(
			&client.ID,
			&client.Secret,
			&client.RedirectURI,
			&requirePkce,
			&client.TokenEndpointAuthMethod,
			&redirectURIs,
		); err != nil {
			return clients, err
		}
		client.RequirePkce = requirePkce != 0
		client.RedirectURIs = splitRedirectURIs(redirectURIs, client.RedirectURI)
		clients = append(clients, client)
	}
	return clients, rows.Err()
}

func (o *OIDC) AddClient(client OIDCClient) error {
	requirePkce := 0
	if client.RequirePkce {
		requirePkce = 1
	}
	uris := client.RedirectURIs
	if len(uris) == 0 && client.RedirectURI != "" {
		uris = []string{client.RedirectURI}
	}
	first := ""
	if len(uris) > 0 {
		first = uris[0]
	}
	_, err := o.db.Exec(
		"INSERT OR REPLACE INTO oidc_client(id, secret, redirect_uri, require_pkce, token_endpoint_auth_method, redirect_uris) VALUES (?, ?, ?, ?, ?, ?)",
		client.ID, client.Secret, first, requirePkce, client.TokenEndpointAuthMethod, strings.Join(uris, " "),
	)
	return err
}

func splitRedirectURIs(redirectURIs, redirectURI string) []string {
	if fields := strings.Fields(redirectURIs); len(fields) > 0 {
		return fields
	}
	if redirectURI != "" {
		return []string{redirectURI}
	}
	return nil
}
