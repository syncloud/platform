package config

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db *Db
}

func NewMigrator(db *Db) *Migrator {
	return &Migrator{db: db}
}

func migrations() []*goose.Migration {
	return []*goose.Migration{
		goose.NewGoMigration(1, &goose.GoFunc{RunTx: createConfigTable}, nil),
		goose.NewGoMigration(2, &goose.GoFunc{RunTx: createOidcClientTable}, nil),
		goose.NewGoMigration(3, &goose.GoFunc{RunTx: createCustomProxyTable}, nil),
		goose.NewGoMigration(4, &goose.GoFunc{RunTx: addCustomProxyHttps}, nil),
		goose.NewGoMigration(5, &goose.GoFunc{RunTx: addCustomProxyAuthelia}, nil),
		goose.NewGoMigration(6, &goose.GoFunc{RunTx: normalizeOidcRedirectUris}, nil),
	}
}

func (m *Migrator) provider() (*goose.Provider, error) {
	return goose.NewProvider(goose.DialectSQLite3, m.db.Open(), nil, goose.WithGoMigrations(migrations()...))
}

func (m *Migrator) Migrate() error {
	provider, err := m.provider()
	if err != nil {
		return err
	}
	defer provider.Close()
	_, err = provider.Up(context.Background())
	return err
}

func (m *Migrator) MigrateTo(version int64) error {
	provider, err := m.provider()
	if err != nil {
		return err
	}
	defer provider.Close()
	_, err = provider.UpTo(context.Background(), version)
	return err
}

func createConfigTable(_ context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("create table if not exists config (key varchar primary key, value varchar)")
	return err
}

func createOidcClientTable(_ context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`create table if not exists oidc_client
		(id varchar primary key, secret varchar, redirect_uri varchar, require_pkce integer, token_endpoint_auth_method varchar)`)
	return err
}

func createCustomProxyTable(_ context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("create table if not exists custom_proxy (name varchar primary key, host varchar, port integer)")
	return err
}

func addCustomProxyHttps(ctx context.Context, tx *sql.Tx) error {
	return addColumnIfMissing(ctx, tx, "custom_proxy", "https", "integer not null default 0")
}

func addCustomProxyAuthelia(ctx context.Context, tx *sql.Tx) error {
	return addColumnIfMissing(ctx, tx, "custom_proxy", "authelia", "integer not null default 0")
}

func normalizeOidcRedirectUris(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`create table if not exists oidc_redirect_uri
		(client_id varchar not null, redirect_uri varchar not null)`)
	if err != nil {
		return err
	}

	hasLegacyColumn, err := columnExists(ctx, tx, "oidc_client", "redirect_uri")
	if err != nil {
		return err
	}
	if !hasLegacyColumn {
		return nil
	}

	_, err = tx.Exec(`insert into oidc_redirect_uri (client_id, redirect_uri)
		select id, redirect_uri from oidc_client
		where redirect_uri != '' and id not in (select client_id from oidc_redirect_uri)`)
	if err != nil {
		return err
	}

	_, err = tx.Exec("ALTER TABLE oidc_client DROP COLUMN redirect_uri")
	return err
}

func addColumnIfMissing(ctx context.Context, tx *sql.Tx, table, column, definition string) error {
	exists, err := columnExists(ctx, tx, table, column)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = tx.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, definition))
	return err
}

func columnExists(ctx context.Context, tx *sql.Tx, table, column string) (bool, error) {
	rows, err := tx.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var cid, notnull, pk int
		var name, columnType string
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &columnType, &notnull, &defaultValue, &pk); err != nil {
			return false, err
		}
		if name == column {
			return true, nil
		}
	}
	return false, rows.Err()
}
