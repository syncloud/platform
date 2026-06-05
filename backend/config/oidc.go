package config

type OIDCClient struct {
	ID                      string
	Secret                  string
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

	uriRows, err := db.Query("select client_id, redirect_uri from oidc_redirect_uri order by rowid")
	if err != nil {
		return nil, err
	}
	urisByClient := map[string][]string{}
	for uriRows.Next() {
		var clientID, redirectURI string
		if err := uriRows.Scan(&clientID, &redirectURI); err != nil {
			uriRows.Close()
			return nil, err
		}
		urisByClient[clientID] = append(urisByClient[clientID], redirectURI)
	}
	uriRows.Close()
	if err := uriRows.Err(); err != nil {
		return nil, err
	}

	rows, err := db.Query("select id, secret, require_pkce, token_endpoint_auth_method from oidc_client")
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
			&requirePkce,
			&client.TokenEndpointAuthMethod,
		); err != nil {
			return clients, err
		}
		client.RequirePkce = requirePkce != 0
		client.RedirectURIs = urisByClient[client.ID]
		clients = append(clients, client)
	}
	return clients, rows.Err()
}

func (o *OIDC) AddClient(client OIDCClient) error {
	requirePkce := 0
	if client.RequirePkce {
		requirePkce = 1
	}

	db := o.db.Open()
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(
		"INSERT OR REPLACE INTO oidc_client(id, secret, require_pkce, token_endpoint_auth_method) VALUES (?, ?, ?, ?)",
		client.ID, client.Secret, requirePkce, client.TokenEndpointAuthMethod,
	); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec("DELETE FROM oidc_redirect_uri WHERE client_id = ?", client.ID); err != nil {
		tx.Rollback()
		return err
	}
	for _, redirectURI := range client.RedirectURIs {
		if _, err := tx.Exec(
			"INSERT INTO oidc_redirect_uri(client_id, redirect_uri) VALUES (?, ?)",
			client.ID, redirectURI,
		); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
