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
		_, err := g.generateFake()
		return err
	}

	err := g.generateReal()
	if err != nil {
		g.logger.Info(fmt.Sprintf("unable to generate certificate: %s", err.Error()))
		generated, err := g.generateFake()
		if err != nil {
			return err
		}
		if generated {
			return g.nginx.ReloadPublic()
		}
	}
	return nil
}

func (g *CertificateGenerator) generateReal() error {
	certInfo := g.ReadCertificateInfo()
	g.logger.Info("certificate info", zap.Int("valid days", certInfo.ValidForDays), zap.Bool("real", certInfo.IsReal))

	if certInfo.IsValid && certInfo.IsReal {
		g.logger.Info("not regenerating real certificate")
		return nil
	}

	err := g.certbot.Generate()
	if err == nil {
		err = g.nginx.ReloadPublic()
	}
	return err
}

func (g *CertificateGenerator) generateFake() (bool, error) {
	certInfo := g.ReadCertificateInfo()
	if certInfo.IsValid {
		return false, nil
	}
	err := g.fake.Generate()
	if err != nil {
		g.logger.Info(fmt.Sprintf("unable to generate fake certificate: %s", err.Error()))
		return false, err
	}
	return true, nil
}

func (g *CertificateGenerator) ReadCertificateInfo() *Info {

	certBytes, err := ioutil.ReadFile(g.systemConfig.SslCertificateFile())
	if err != nil {
		g.logger.Info(fmt.Sprintf("unable to read certificate file: %s", err.Error()))
		return &Info{}
	}

	block, _ := pem.Decode(certBytes)
	certificateData, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		g.logger.Info(fmt.Sprintf("unable to parse certificate: %s", err.Error()))
		return &Info{}
	}

	now := g.dateProvider.Now()
	validFor := certificateData.NotAfter.Sub(now)
	subject := certificateData.Subject.String()
	commonName := certificateData.Subject.CommonName
	return &Info{
		IsValid:      validFor > Month,
		Subject:      subject,
		ValidForDays: int(validFor.Hours() / 24),
		IsReal:       commonName != DefaultSubjectCommonName && commonName != DefaultSubjectCommonNameOld,
	}
}
