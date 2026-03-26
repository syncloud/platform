package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/parser"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type Web interface {
	InitConfig() error
	WaitForReady() error
}

type Variables struct {
	Domain            string
	AppUrl            string
	EncryptionKey     string
	JwtSecret         string
	HmacSecret        string
	DeviceUrl         string
	AuthUrl           string
	IsActivated       bool
	TwoFactorEnabled  bool
	OIDCClients       []config.OIDCClient
}

type HealthWaiter interface {
	WaitForReady() error
}

type Authelia struct {
	mutex          *sync.Mutex
	inputDir       string
	outDir         string
	dataDir        string
	keyFile        string
	secretFile     string
	jwksKeyFile    string
	hmacSecretFile string
	socketPath     string
	userConfig     UserConfig
	systemd        Systemd
	generator      PasswordGenerator
	executor       cli.Executor
	health         HealthWaiter
	logger         *zap.Logger
}

type UserConfig interface {
	GetDeviceDomain() string
	DeviceUrl() string
	Url(app string) string
	OIDCClients() ([]config.OIDCClient, error)
	AddOIDCClient(client config.OIDCClient) error
	IsActivated() bool
	IsTwoFactorEnabled() bool
}

type Systemd interface {
	RestartService(service string) error
}

type PasswordGenerator interface {
	Generate() (Secret, error)
}

const (
	KeyFile    = "authelia.storage.encryption.key"
	SecretFile = "authelia.jwt.secret"
	JwksKey    = "authelia.jwks.key"
	HmacSecret = "authelia.hmac_secret.key"
)

func NewAuthelia(
	inputDir string,
	outDir string,
	outSecretDir string,
	socketPath string,
	userConfig UserConfig,
	systemd Systemd,
	generator PasswordGenerator,
	executor cli.Executor,
	health HealthWaiter,
	logger *zap.Logger,
) *Authelia {
	return &Authelia{
		mutex:          &sync.Mutex{},
		inputDir:       inputDir,
		outDir:         outDir,
		dataDir:        outSecretDir,
		keyFile:        path.Join(outSecretDir, KeyFile),
		secretFile:     path.Join(outSecretDir, SecretFile),
		jwksKeyFile:    path.Join(outSecretDir, JwksKey),
		hmacSecretFile: path.Join(outSecretDir, HmacSecret),
		socketPath:     socketPath,
		userConfig:     userConfig,
		systemd:        systemd,
		generator:      generator,
		executor:       executor,
		health:         health,
		logger:         logger,
	}
}

func (w *Authelia) RegisterOIDCClient(
	id string,
	redirectURI string,
	requirePkce bool,
	tokenEndpointAuthMethod string,
) (string, error) {
	secret, err := w.generator.Generate()
	if err != nil {
		return "", err
	}

	err = w.userConfig.AddOIDCClient(config.OIDCClient{
		ID:                      id,
		Secret:                  secret.Hash,
		RedirectURI:             redirectURI,
		RequirePkce:             requirePkce,
		TokenEndpointAuthMethod: tokenEndpointAuthMethod,
	})
	if err != nil {
		return "", err
	}
	err = w.InitConfig()
	if err != nil {
		return "", err
	}
	err = w.WaitForReady()
	if err != nil {
		return "", err
	}
	return secret.Password, nil
}

func (w *Authelia) InitConfig() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	activated := w.userConfig.IsActivated()
	encryptionKey, err := getOrCreateUuid(w.keyFile)
	if err != nil {
		return err
	}
	jwtSecret, err := getOrCreateUuid(w.secretFile)
	if err != nil {
		return err
	}
	hmacSecret, err := getOrCreateUuid(w.hmacSecretFile)
	if err != nil {
		return err
	}
	err = createRsaKeyFileIfMissing(w.jwksKeyFile)
	if err != nil {
		return err
	}

	clients, err := w.userConfig.OIDCClients()
	if err != nil {
		return err
	}
	variables := Variables{
		Domain:           w.userConfig.GetDeviceDomain(),
		EncryptionKey:    encryptionKey,
		JwtSecret:        jwtSecret,
		HmacSecret:       hmacSecret,
		DeviceUrl:        w.userConfig.DeviceUrl(),
		AuthUrl:          w.userConfig.Url("auth"),
		IsActivated:      activated,
		TwoFactorEnabled: w.userConfig.IsTwoFactorEnabled(),
		OIDCClients:      clients,
	}

	tmpDir := w.outDir + ".tmp"
	err = os.RemoveAll(tmpDir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	err = copyAssetsDir(
		path.Join(w.inputDir, "assets"),
		path.Join(tmpDir, "assets"),
	)
	if err != nil {
		w.logger.Warn("unable to copy authelia assets", zap.Error(err))
	}
	_ = os.MkdirAll(path.Join(w.outDir, "assets"), 0755)

	err = parser.Generate(
		w.inputDir,
		tmpDir,
		variables,
	)
	if err != nil {
		return err
	}

	output, err := w.executor.CombinedOutput(
		"snap", "run", "platform.authelia-cli",
		"validate-config",
		"--config", path.Join(tmpDir, "config.yml"),
		"--config.experimental.filters", "template",
	)
	if err != nil {
		return fmt.Errorf("authelia config validation failed, contact support: %s", string(output))
	}

	err = copyDir(tmpDir, w.outDir)
	if err != nil {
		return err
	}

	err = w.systemd.RestartService("platform.authelia")
	if err != nil {
		w.logger.Error("unable to restart authelia", zap.Error(err))
		return err
	}

	return nil
}

func (w *Authelia) WaitForReady() error {
	return w.health.WaitForReady()
}

func copyDir(srcDir string, dstDir string) error {
	err := os.MkdirAll(dstDir, 0755)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		src := filepath.Join(srcDir, entry.Name())
		dst := filepath.Join(dstDir, entry.Name())
		if entry.IsDir() {
			err = copyDir(src, dst)
			if err != nil {
				return err
			}
			continue
		}
		srcFile, err := os.Open(src)
		if err != nil {
			return err
		}
		dstFile, err := os.Create(dst)
		if err != nil {
			srcFile.Close()
			return err
		}
		_, err = io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func getOrCreateUuid(file string) (string, error) {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		secret := uuid.New().String()
		err = os.WriteFile(file, []byte(secret), 0644)
		return secret, err
	}
	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func copyAssetsDir(srcDir string, dstDir string) error {
	_, err := os.Stat(srcDir)
	if os.IsNotExist(err) {
		return nil
	}
	err = os.MkdirAll(dstDir, 0755)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := filepath.Join(srcDir, entry.Name())
		dst := filepath.Join(dstDir, entry.Name())
		srcFile, err := os.Open(src)
		if err != nil {
			return err
		}
		dstFile, err := os.Create(dst)
		if err != nil {
			srcFile.Close()
			return err
		}
		_, err = io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Authelia) GenerateTOTP(username string) (string, error) {
	encryptionKey, err := os.ReadFile(w.keyFile)
	if err != nil {
		return "", fmt.Errorf("failed to read encryption key: %w", err)
	}
	sqlitePath := path.Join(w.dataDir, "authelia.sqlite3")

	output, err := w.executor.CombinedOutput(
		"snap", "run", "platform.authelia-cli",
		"storage", "user", "totp", "generate", username,
		"--force",
		"--issuer", "Syncloud",
		"--sqlite.path", sqlitePath,
		"--encryption-key", string(encryptionKey),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP: %s", string(output))
	}

	// Parse the URI from output: "Successfully generated TOTP configuration for user 'X' with URI 'otpauth://...'"
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

func (w *Authelia) HasTOTP(username string) (bool, error) {
	sqlitePath := path.Join(w.dataDir, "authelia.sqlite3")
	db, err := sql.Open("sqlite", sqlitePath)
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

func (w *Authelia) ResetAllTOTP() error {
	sqlitePath := path.Join(w.dataDir, "authelia.sqlite3")
	db, err := sql.Open("sqlite", sqlitePath)
	if err != nil {
		return fmt.Errorf("failed to open authelia db: %w", err)
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM totp_configurations")
	return err
}

func (w *Authelia) ResetTOTP(username string) error {
	encryptionKey, err := os.ReadFile(w.keyFile)
	if err != nil {
		return fmt.Errorf("failed to read encryption key: %w", err)
	}
	sqlitePath := path.Join(w.dataDir, "authelia.sqlite3")
	output, err := w.executor.CombinedOutput(
		"snap", "run", "platform.authelia-cli",
		"storage", "user", "totp", "delete", username,
		"--sqlite.path", sqlitePath,
		"--encryption-key", string(encryptionKey),
	)
	if err != nil {
		return fmt.Errorf("failed to delete TOTP: %s", string(output))
	}
	return nil
}

func createRsaKeyFileIfMissing(file string) error {
	_, err := os.Stat(file)
	if err == nil || !os.IsNotExist(err) {
		return err
	}
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	err = os.WriteFile(file, keyPEM, 0700)
	if err != nil {
		return err
	}
	return nil
}
