package auth

import (
	"database/sql"
	"fmt"
	"path"
)

type Passkeys struct {
	sqlitePath string
}

func NewPasskeys(dataDir string) *Passkeys {
	return &Passkeys{sqlitePath: path.Join(dataDir, "authelia.sqlite3")}
}

func (p *Passkeys) List(username string) ([]Passkey, error) {
	db, err := sql.Open("sqlite", autheliaDSN(p.sqlitePath))
	if err != nil {
		return nil, fmt.Errorf("failed to open authelia db: %w", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, description, created_at, last_used_at FROM webauthn_credentials WHERE username = ? ORDER BY created_at", username)
	if err != nil {
		return nil, fmt.Errorf("failed to list passkeys: %w", err)
	}
	defer rows.Close()
	passkeys := make([]Passkey, 0)
	for rows.Next() {
		var passkey Passkey
		var lastUsed sql.NullTime
		if err := rows.Scan(&passkey.ID, &passkey.Description, &passkey.Created, &lastUsed); err != nil {
			return nil, fmt.Errorf("failed to read passkey: %w", err)
		}
		if lastUsed.Valid {
			passkey.LastUsed = &lastUsed.Time
		}
		passkeys = append(passkeys, passkey)
	}
	return passkeys, rows.Err()
}

func (p *Passkeys) Remove(username string, id int64) error {
	return p.delete("DELETE FROM webauthn_credentials WHERE username = ? AND id = ?", username, id)
}

func (p *Passkeys) Reset(username string) error {
	return p.delete("DELETE FROM webauthn_credentials WHERE username = ?", username)
}

func (p *Passkeys) delete(query string, args ...any) error {
	db, err := sql.Open("sqlite", autheliaDSN(p.sqlitePath))
	if err != nil {
		return fmt.Errorf("failed to open authelia db: %w", err)
	}
	defer db.Close()
	if _, err := db.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to remove passkeys: %w", err)
	}
	return nil
}
