package cert

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/syncloud/platform/date"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

const (
	Day   = time.Hour * 24
	Month = Day * 30
)

type Generator interface {
	Generate() error
}

type cert struct {
	validFor time.Duration
	subject  string
}

func (c *cert) IsValid() bool {
	return c.validFor > Month
}

func (c *cert) IsReal() bool {
	return c.subject != fmt.Sprintf("CN=%s,O=%s,L=%s,ST=%s,C=%s", DefaultSubjectCommonName, SubjectOrganization, SubjectLocality, SubjectProvince, SubjectCountry)
}

func (c *cert) ValidForDays() int {
	return int(c.validFor.Hours() / 24)
}

type CertificateGenerator struct {
	systemConfig GeneratorSystemConfig
	userConfig   GeneratorUserConfig
	certbot      CertbotGenerator
	fake         FakeGenerator
	dateProvider date.Provider
	nginx        GeneratorNginx
	logger       *zap.Logger
}

type GeneratorSystemConfig interface {
	SslCertificateFile() string
	SslKeyFile() string
}

type GeneratorUserConfig interface {
	IsActivated() bool
}

type GeneratorNginx interface {
	ReloadPublic() error
}

func New(systemConfig GeneratorSystemConfig, userConfig GeneratorUserConfig, dateProvider date.Provider, certbot CertbotGenerator, fake FakeGenerator, nginx GeneratorNginx, logger *zap.Logger) *CertificateGenerator {
	return &CertificateGenerator{
		systemConfig: systemConfig,
		userConfig:   userConfig,
		certbot:      certbot,
		fake:         fake,
		dateProvider: dateProvider,
		nginx:        nginx,
		logger:       logger,
	}
}

func (g *CertificateGenerator) Generate() error {

	if !g.userConfig.IsActivated() {
		return g.generateFake()
	}

	err := g.generateReal()
	if err != nil {
		g.logger.Info(fmt.Sprintf("unable to generate certificate: %s", err.Error()))
		return g.generateFake()
	}
	return nil
}

func (g *CertificateGenerator) generateReal() error {
	certInfo := g.readCertificateInfo()
	g.logger.Info("certificate info", zap.Int("valid days", certInfo.ValidForDays()), zap.Bool("real", certInfo.IsReal()))

	if certInfo.IsValid() && certInfo.IsReal() {
		g.logger.Info("not regenerating real certificate")
		return nil
	}

	err := g.certbot.Generate()
	if err == nil {
		err = g.nginx.ReloadPublic()
	}
	return err
}

func (g *CertificateGenerator) generateFake() error {
	certInfo := g.readCertificateInfo()
	if certInfo.IsValid() {
		return nil
	}
	err := g.fake.Generate()
	if err != nil {
		g.logger.Info(fmt.Sprintf("unable to generate fake certificate: %s", err.Error()))
		return err
	}
	return g.nginx.ReloadPublic()
}

func (g *CertificateGenerator) readCertificateInfo() *cert {

	certBytes, err := ioutil.ReadFile(g.systemConfig.SslCertificateFile())
	if err != nil {
		g.logger.Info(fmt.Sprintf("unable to read certificate file: %s", err.Error()))
		return &cert{0, ""}
	}

	block, _ := pem.Decode(certBytes)
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		g.logger.Info(fmt.Sprintf("unable to parse certificate: %s", err.Error()))
		return &cert{0, ""}
	}

	now := g.dateProvider.Now()
	validFor := certificate.NotAfter.Sub(now)
	subject := certificate.Subject.String()
	return &cert{validFor, subject}
}
