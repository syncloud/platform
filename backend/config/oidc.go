package config

type OIDCClient struct {
	ID                      string
	Secret                  string
	RedirectURI             string
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
	rows, err := db.Query("select id, secret, redirect_uri, require_pkce, token_endpoint_auth_method from oidc_client")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := make([]OIDCClient, 0)
	for rows.Next() {
		var client OIDCClient
		var requirePkce int
		if err := rows.Scan(
			&client.ID,
			&client.Secret,
			&client.RedirectURI,
			&requirePkce,
			&client.TokenEndpointAuthMethod,
		); err != nil {
			return clients, err
		}
		client.RequirePkce = requirePkce != 0
		clients = append(clients, client)
	}
	return clients, rows.Err()
}

func (o *OIDC) AddClient(client OIDCClient) error {
	requirePkce := 0
	if client.RequirePkce {
		requirePkce = 1
	}
	_, err := o.db.Exec("INSERT OR REPLACE INTO oidc_client VALUES (?, ?, ?, ?, ?)",
		client.ID, client.Secret, client.RedirectURI, requirePkce, client.TokenEndpointAuthMethod,
	)
	return err
}
