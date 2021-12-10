package cert

import (
	"crypto/x509"
	"github.com/syncloud/platform/date"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

const (
	Log   = "/var/snap/platform/common/log/certbot.log"
	Day   = time.Hour * 24
	Month = Day * 30
)

type Generator interface {
	Generate() error
}

type CertificateGenerator struct {
	systemConfig GeneratorSystemConfig
	userConfig   GeneratorUserConfig
	certbot      CertbotGenerator
	fake         FakeGenerator
	dateProvider date.Provider
	logger       *zap.Logger
}

type GeneratorSystemConfig interface {
	SslCertificateFile() string
	SslKeyFile() string
}

type GeneratorUserConfig interface {
	IsActivated() bool
}

func New(systemConfig GeneratorSystemConfig, userConfig GeneratorUserConfig, dateProvider date.Provider, certbot CertbotGenerator, fake FakeGenerator, logger *zap.Logger) *CertificateGenerator {
	return &CertificateGenerator{
		systemConfig: systemConfig,
		userConfig:   userConfig,
		certbot:      certbot,
		fake:         fake,
		dateProvider: dateProvider,
		logger:       logger,
	}
}

func (g *CertificateGenerator) Generate() error {

	if !g.userConfig.IsActivated() {
		return g.generateFake()
	}

	if !g.isExpired() {
		return nil
	}

	err := g.certbot.Generate()
	if err != nil {
		g.Log("unable to generate fake certificate: %v\n", err)
		return g.generateFake()
	}
	return nil
}

func (g *CertificateGenerator) isExpired() bool {

	certBytes, err := ioutil.ReadFile(g.systemConfig.SslCertificateFile())
	if err != nil {
		return true
	}

	certificate, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return true
	}

	now := g.dateProvider.Now()
	validFor := certificate.NotAfter.Sub(now)
	valid := validFor > Month
	realCert := certificate.Subject.String() != Subject
	if valid && realCert {
		g.Log("not regenerating real certificate, valid for days: %d\n", int(validFor.Hours()/24))
		return false
	}

	return true
}

func (g *CertificateGenerator) generateFake() error {
	err := g.fake.Generate()
	if err != nil {
		g.logger.Info("unable to generate fake certificate: %v\n", err.Error())
		return err
	}
	return nil
}
