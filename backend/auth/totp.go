package auth

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/syncloud/platform/cli"
	"go.uber.org/zap"
)

type TOTP struct {
	executor cli.Executor
	dataDir  string
	keyFile  string
	logger   *zap.Logger
}

func NewTOTP(executor cli.Executor, dataDir string, logger *zap.Logger) *TOTP {
	return &TOTP{
		executor: executor,
		dataDir:  dataDir,
		keyFile:  path.Join(dataDir, KeyFile),
		logger:   logger,
	}
}

func (t *TOTP) sqlitePath() string {
	return path.Join(t.dataDir, "authelia.sqlite3")
}

func (t *TOTP) Generate(username string) (string, error) {
	encryptionKey, err := os.ReadFile(t.keyFile)
	if err != nil {
		return "", fmt.Errorf("failed to read encryption key: %w", err)
	}

	output, err := t.executor.CombinedOutput(
		"snap", "run", "platform.authelia-cli",
		"storage", "user", "totp", "generate", username,
		"--force",
		"--issuer", "Syncloud",
		"--sqlite.path", t.sqlitePath(),
		"--encryption-key", string(encryptionKey),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP: %s", string(output))
	}

	outputStr := string(output)
	uriStart := strings.Index(outputStr, "otpauth://")
	if uriStart == -1 {
		return "", fmt.Errorf("TOTP URI not found in output: %s", outputStr)
	}
	uri := outputStr[uriStart:]
	if idx := strings.Index(uri, "'"); idx != -1 {
		uri = uri[:idx]
	}
	return strings.TrimSpace(uri), nil
}

func (t *TOTP) Has(username string) (bool, error) {
	db, err := sql.Open("sqlite", autheliaDSN(t.sqlitePath()))
	if err != nil {
		return false, fmt.Errorf("failed to open authelia db: %w", err)
	}
	defer db.Close()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM totp_configurations WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, nil
	}
	return count > 0, nil
}

func (t *TOTP) ResetAll() error {
	db, err := sql.Open("sqlite", autheliaDSN(t.sqlitePath()))
	if err != nil {
		return fmt.Errorf("failed to open authelia db: %w", err)
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM totp_configurations")
	return err
}

func (t *TOTP) Reset(username string) error {
	encryptionKey, err := os.ReadFile(t.keyFile)
	if err != nil {
		return fmt.Errorf("failed to read encryption key: %w", err)
	}
	output, err := t.executor.CombinedOutput(
		"snap", "run", "platform.authelia-cli",
		"storage", "user", "totp", "delete", username,
		"--sqlite.path", t.sqlitePath(),
		"--encryption-key", string(encryptionKey),
	)
	if err != nil {
		return fmt.Errorf("failed to delete TOTP: %s", string(output))
	}
	return nil
}
