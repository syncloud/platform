package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/syncloud/platform/date"
	"go.uber.org/zap"
	"io/ioutil"
	"math/big"
	"time"
)

const (
	SubjectCountry           = "UK"
	SubjectProvince          = "Syncloud"
	SubjectLocality          = "Syncloud"
	SubjectOrganization      = "Syncloud"
	DefaultSubjectCommonName = "syncloud"
	DefaultDuration          = 2 * Month
)

type Fake struct {
	systemConfig      GeneratorSystemConfig
	dateProvider      date.Provider
	subjectCommonName string
	duration          time.Duration
	logger            *zap.Logger
}

type FakeGenerator interface {
	Generate() error
}

func NewFake(systemConfig GeneratorSystemConfig, dateProvider date.Provider, subjectCommonName string, duration time.Duration, logger *zap.Logger) *Fake {
	return &Fake{
		systemConfig:      systemConfig,
		dateProvider:      dateProvider,
		subjectCommonName: subjectCommonName,
		duration:          duration,
		logger:            logger,
	}
}

func (c *Fake) Generate() error {
	c.logger.Info("generating self signed certificate")

	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return err
	}

	subject := pkix.Name{
		Country:      []string{SubjectCountry},
		Province:     []string{SubjectProvince},
		Locality:     []string{SubjectLocality},
		Organization: []string{SubjectOrganization},
		CommonName:   c.subjectCommonName,
	}
	now := c.dateProvider.Now()

	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               subject,
		NotBefore:             now,
		NotAfter:              now.Add(c.duration),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, privateKey.Public(), privateKey)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.systemConfig.SslKeyFile(), privateKeyBytes, 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.systemConfig.SslCertificateFile(), certificateBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
