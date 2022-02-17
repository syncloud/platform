package cert

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/syncloud/platform/date"
	"go.uber.org/zap"
	"io/ioutil"
	"math/big"
	"time"
)

const (
	SubjectCountry      = "UK"
	SubjectProvince     = "Syncloud"
	SubjectLocality     = "Syncloud"
	SubjectOrganization = "Syncloud"
	SubjectCommonName   = "syncloud"
  SubjectCommonNameCa   = "syncloud ca"
	DefaultDuration     = 2 * Month
)

type Fake struct {
	systemConfig        GeneratorSystemConfig
	dateProvider        date.Provider
	subjectOrganization string
	duration            time.Duration
	logger              *zap.Logger
}

type FakeGenerator interface {
	Generate() error
}

func NewFake(systemConfig GeneratorSystemConfig, dateProvider date.Provider, subjectOrganization string, duration time.Duration, logger *zap.Logger) *Fake {
	return &Fake{
		systemConfig:        systemConfig,
		dateProvider:        dateProvider,
		subjectOrganization: subjectOrganization,
		duration:            duration,
		logger:              logger,
	}
}

func (c *Fake) Generate() error {
	c.logger.Info("generating fake certificate")


	now := c.dateProvider.Now()

	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().UnixNano() / int64(time.Millisecond)),
		Subject:               pkix.Name{
		Country:      []string{SubjectCountry},
		Province:     []string{SubjectProvince},
		Locality:     []string{SubjectLocality},
		Organization: []string{c.subjectOrganization},
		CommonName:   SubjectCommonNameCa,
	},
		NotBefore:             now,
		NotAfter:              now.Add(c.duration),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	caPrivKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}
	caPrivateKeyBytes, err := x509.MarshalECPrivateKey(caPrivKey)
	if err != nil {
		return err
	}
	caPrivateKeyPem := &bytes.Buffer{}
	err = pem.Encode(caPrivateKeyPem, &pem.Block{Type: "EC PRIVATE KEY", Bytes: caPrivateKeyBytes})
	if err != nil {
		return err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, caPrivKey.Public(), caPrivKey)
	if err != nil {
		return err
	}
	caPem := &bytes.Buffer{}
	err = pem.Encode(caPem, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes})
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.systemConfig.SslCaKeyFile(), caPrivateKeyPem.Bytes(), 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.systemConfig.SslCaCertificateFile(), caPem.Bytes(), 0644)
	if err != nil {
		return err
	}

	cert := &x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().UnixNano()/int64(time.Millisecond) + 1),
		Subject:               pkix.Name{
		Country:      []string{SubjectCountry},
		Province:     []string{SubjectProvince},
		Locality:     []string{SubjectLocality},
		Organization: []string{c.subjectOrganization},
		CommonName:   SubjectCommonName,
	},
		NotBefore:             now,
		NotAfter:              now.Add(c.duration),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}
	privateKeyPem := &bytes.Buffer{}
	err = pem.Encode(privateKeyPem, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privateKeyBytes})
	if err != nil {
		return err
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, privateKey.Public(), caPrivKey)
	if err != nil {
		return err
	}
	certificatePem := &bytes.Buffer{}
	err = pem.Encode(certificatePem, &pem.Block{Type: "CERTIFICATE", Bytes: certificateBytes})
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.systemConfig.SslKeyFile(), privateKeyPem.Bytes(), 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.systemConfig.SslCertificateFile(), certificatePem.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
