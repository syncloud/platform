package config

import (
	"fmt"
	"os"
)

type Migrator struct {
	db *Db
}

func NewMigrator(db *Db) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Migrate() error {
	_, err := os.Stat(m.db.File())
	if os.IsNotExist(err) {
		if err := m.db.Init(); err != nil {
			return err
		}
	}

	if err := m.addOidcClientTable(); err != nil {
		return err
	}
	if err := m.addCustomProxyTable(); err != nil {
		return err
	}
	if err := m.migrateCustomProxyHttps(); err != nil {
		return err
	}
	if err := m.migrateCustomProxyAuthelia(); err != nil {
		return err
	}
	if err := m.migrateOidcRedirectUris(); err != nil {
		return err
	}
	return nil
}

func (m *Migrator) migrateOidcRedirectUris() error {
	_, _ = m.db.Exec("ALTER TABLE oidc_client ADD COLUMN redirect_uris varchar not null default ''")
	return nil
}

func (m *Migrator) addOidcClientTable() error {
	_, err := m.db.Exec(`create table if not exists oidc_client
		(id varchar primary key, secret varchar, redirect_uri varchar, require_pkce integer, token_endpoint_auth_method varchar)`)
	if err != nil {
		return fmt.Errorf("unable to add oidc_clients: %s", err)
	}
	return nil
}

func (m *Migrator) addCustomProxyTable() error {
	_, err := m.db.Exec(`create table if not exists custom_proxy
		(name varchar primary key, host varchar, port integer)`)
	if err != nil {
		return fmt.Errorf("unable to add custom_proxy table: %s", err)
	}
	return nil
}

func (m *Migrator) migrateCustomProxyHttps() error {
	_, _ = m.db.Exec("ALTER TABLE custom_proxy ADD COLUMN https integer not null default 0")
	return nil
}

func (m *Migrator) migrateCustomProxyAuthelia() error {
	_, _ = m.db.Exec("ALTER TABLE custom_proxy ADD COLUMN authelia integer not null default 0")
	return nil
}
