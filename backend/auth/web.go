package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/google/uuid"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/parser"
	"go.uber.org/zap"
	"os"
	"path"
)

type Web interface {
	InitConfig(activated bool) error
}

type Variables struct {
	Domain        string
	AppUrl        string
	EncryptionKey string
	JwtSecret     string
	HmacSecret    string
	DeviceUrl     string
	AuthUrl       string
	IsActivated   bool
	OIDCClients   []config.OIDCClient
}

type Authelia struct {
	inputDir       string
	outDir         string
	keyFile        string
	secretFile     string
	jwksKeyFile    string
	hmacSecretFile string
	userConfig     UserConfig
	systemd        Systemd
	logger         *zap.Logger
}

type UserConfig interface {
	GetDeviceDomainNil() *string
	DeviceUrl() string
	Url(app string) string
	OIDCClients() ([]config.OIDCClient, error)
}

type Systemd interface {
	RestartService(service string) error
}

const (
	KeyFile    = "authelia.storage.encryption.key"
	SecretFile = "authelia.jwt.secret"
	JwksKey    = "authelia.jwks.key"
	HmacSecret = "authelia.hmac_secret.key"
)

func NewWeb(
	inputDir string,
	outDir string,
	outSecretDir string,
	userConfig UserConfig,
	systemd Systemd,
	logger *zap.Logger,
) *Authelia {
	return &Authelia{
		inputDir:       inputDir,
		outDir:         outDir,
		keyFile:        path.Join(outSecretDir, KeyFile),
		secretFile:     path.Join(outSecretDir, SecretFile),
		jwksKeyFile:    path.Join(outSecretDir, JwksKey),
		hmacSecretFile: path.Join(outSecretDir, HmacSecret),
		userConfig:     userConfig,
		systemd:        systemd,
		logger:         logger,
	}
}

func (w *Authelia) InitConfig(activated bool) error {

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
		Domain:        "www.localhost",
		EncryptionKey: encryptionKey,
		JwtSecret:     jwtSecret,
		HmacSecret:    hmacSecret,
		DeviceUrl:     "https://www.localhost",
		AuthUrl:       "https://auth.www.localhost",
		IsActivated:   activated,
		OIDCClients:   clients,
	}

	maybeDomain := w.userConfig.GetDeviceDomainNil()
	if maybeDomain != nil {
		variables.Domain = *maybeDomain
		variables.DeviceUrl = w.userConfig.DeviceUrl()
		variables.AuthUrl = w.userConfig.Url("auth")
	}

	err = parser.Generate(
		w.inputDir,
		w.outDir,
		variables,
	)
	if err != nil {
		return err
	}

	err = w.systemd.RestartService("platform.authelia")
	if err != nil {
		return err
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
